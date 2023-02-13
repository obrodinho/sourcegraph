import { readFileSync, rmdirSync, writeFileSync } from 'fs'
import * as path from 'path'

import commandExists from 'command-exists'
import { addMinutes } from 'date-fns'
import execa from 'execa'

import * as batchChanges from './batchChanges'
import * as changelog from './changelog'
import {
    addRelease,
    Config,
    loadReleaseConfig,
    newRelease,
    newReleaseFromInput,
    releaseVersions, removeRelease,
    saveReleaseConfig
} from './config'
import {
    getAuthenticatedGitHubClient,
    listIssues,
    getTrackingIssue,
    createChangesets,
    CreatedChangeset,
    createTag,
    ensureTrackingIssues,
    closeTrackingIssue,
    releaseName,
    commentOnIssue,
    queryIssues,
    IssueLabel,
    createLatestRelease,
    cloneRepo,
    getCandidateTags,
    localSourcegraphRepo,
} from './github'
import { ensureEvent, getClient, EventOptions, calendarTime } from './google-calendar'
import { postMessage, slackURL } from './slack'
import {
    cacheFolder,
    formatDate,
    timezoneLink,
    ensureDocker,
    changelogURL,
    ensureReleaseBranchUpToDate,
    ensureSrcCliEndpoint,
    ensureSrcCliUpToDate,
    getLatestTag,
    getAllUpgradeGuides,
    updateUpgradeGuides,
    verifyWithInput, readLine,
} from './util'

const sed = process.platform === 'linux' ? 'sed' : 'gsed'

export type StepID =
    | 'help'
    // release tracking
    | 'tracking:timeline'
    | 'tracking:issues'
    // branch cut
    | 'changelog:cut'
    | 'release:branch-cut'
    // release
    | 'release:status'
    | 'release:create-candidate'
    | 'release:promote-candidate'
    | 'release:check-candidate'
    | 'release:stage'
    | 'release:add-to-batch-change'
    | 'release:finalize'
    | 'release:announce'
    | 'release:close'
    | 'release:prepare'
    | 'release:remove'
    // util
    | 'util:clear-cache'
    // testing
    | '_test:google-calendar'
    | '_test:slack'
    | '_test:batchchange-create-from-changes'
    | '_test:config'
    | '_test:dockerensure'
    | '_test:srccliensure'
    | '_test:release-guide-content'
    | '_test:release-guide-update'

/**
 * Runs given release step with the provided configuration and arguments.
 */
export async function runStep(config: Config, step: StepID, ...args: string[]): Promise<void> {
    if (!steps.map(({ id }) => id as string).includes(step)) {
        throw new Error(`Unrecognized step ${JSON.stringify(step)}`)
    }
    await Promise.all(
        steps
            .filter(({ id }) => id === step)
            .map(async step => {
                if (step.run) {
                    await step.run(config, ...args)
                }
            })
    )
}

interface Step {
    id: StepID
    description: string
    run?: ((config: Config, ...args: string[]) => Promise<void>) | ((config: Config, ...args: string[]) => void)
    argNames?: string[]
}

const steps: Step[] = [
    {
        id: 'help',
        description: 'Output help text about this tool',
        argNames: ['all'],
        run: (_config, all) => {
            console.error('Sourcegraph release tool - https://handbook.sourcegraph.com/engineering/releases')
            console.error('\nUSAGE\n')
            console.error('\tpnpm run release <step>')
            console.error('\nAVAILABLE STEPS\n')
            console.error(
                steps
                    .filter(({ id }) => all || !id.startsWith('_'))
                    .map(
                        ({ id, argNames, description }) =>
                            '\t' +
                            id +
                            (argNames && argNames.length > 0
                                ? ' ' + argNames.map(argumentName => `<${argumentName}>`).join(' ')
                                : '') +
                            '\n\t\t' +
                            description
                    )
                    .join('\n') + '\n'
            )
        },
    },
    {
        id: 'tracking:timeline',
        description: 'Generate a set of Google Calendar events for a MAJOR.MINOR release',
        run: async config => {
            const { upcoming: release } = await releaseVersions(config)
            const name = releaseName(release)
            const events: EventOptions[] = [
                {
                    title: `Security Team to Review Release Container Image Scans ${name}`,
                    description: '(This is not an actual event to attend, just a calendar marker.)',
                    anyoneCanAddSelf: true,
                    attendees: [config.teamEmail],
                    transparency: 'transparent',
                    ...calendarTime(config.oneWorkingWeekBeforeRelease),
                },
                {
                    title: `Cut Sourcegraph ${name}`,
                    description: '(This is not an actual event to attend, just a calendar marker.)',
                    anyoneCanAddSelf: true,
                    attendees: [config.teamEmail],
                    transparency: 'transparent',
                    ...calendarTime(config.threeWorkingDaysBeforeRelease),
                },
                {
                    title: `Release Sourcegraph ${name}`,
                    description: '(This is not an actual event to attend, just a calendar marker.)',
                    anyoneCanAddSelf: true,
                    attendees: [config.teamEmail],
                    transparency: 'transparent',
                    ...calendarTime(config.releaseDate),
                },
                {
                    title: `Start deploying Sourcegraph ${name} to Cloud instances`,
                    description: '(This is not an actual event to attend, just a calendar marker.)',
                    anyoneCanAddSelf: true,
                    attendees: [config.teamEmail],
                    transparency: 'transparent',
                    ...calendarTime(config.oneWorkingDayAfterRelease),
                },
                {
                    title: `All Cloud instances upgraded to Sourcegraph ${name}`,
                    description: '(This is not an actual event to attend, just a calendar marker.)',
                    anyoneCanAddSelf: true,
                    attendees: [config.teamEmail],
                    transparency: 'transparent',
                    ...calendarTime(config.oneWorkingWeekAfterRelease),
                },
            ]

            if (!config.dryRun.calendar) {
                const googleCalendar = await getClient()
                for (const event of events) {
                    console.log(`Create calendar event: ${event.title}: ${event.startDateTime || 'undefined'}`)
                    await ensureEvent(event, googleCalendar)
                }
            } else {
                console.log('dryRun.calendar=true, skipping calendar event creation', events)
            }
        },
    },
    {
        id: 'tracking:issues',
        description: 'Generate GitHub tracking issue for the configured release',
        run: async (config: Config) => {
            const {
                releaseDate,
                captainGitHubUsername,
                oneWorkingWeekBeforeRelease,
                threeWorkingDaysBeforeRelease,
                oneWorkingDayAfterRelease,
                captainSlackUsername,
                slackAnnounceChannel,
                dryRun,
            } = config
            const { upcoming: release } = await releaseVersions(config)
            const date = new Date(releaseDate)

            // Create issue
            const trackingIssues = await ensureTrackingIssues({
                version: release,
                assignees: [captainGitHubUsername],
                releaseDate: date,
                oneWorkingWeekBeforeRelease: new Date(oneWorkingWeekBeforeRelease),
                threeWorkingDaysBeforeRelease: new Date(threeWorkingDaysBeforeRelease),
                oneWorkingDayAfterRelease: new Date(oneWorkingDayAfterRelease),
                dryRun: dryRun.trackingIssues || false,
            })
            console.log('Rendered tracking issues', trackingIssues)

            // If at least one issue was created, post to Slack
            if (trackingIssues.find(({ created }) => created)) {
                const name = releaseName(release)
                const releaseDateString = slackURL(formatDate(date), timezoneLink(date, `${name} release`))
                let annoncement = `:mega: *${name} release*

:captain: Release captain: @${captainSlackUsername}
:spiral_calendar_pad: Scheduled for: ${releaseDateString}
:pencil: Tracking issues:
${trackingIssues.map(index => `- ${slackURL(index.title, index.url)}`).join('\n')}`
                if (release.patch !== 0) {
                    const patchRequestTemplate = `https://github.com/sourcegraph/sourcegraph/issues/new?assignees=&labels=team%2Fdistribution&template=request_patch_release.md&title=${release.version}%3A+`
                    annoncement += `\n\nIf you have changes that should go into this patch release, ${slackURL(
                        'please *file a patch request issue*',
                        patchRequestTemplate
                    )}, or it will not be included.`
                }
                if (!dryRun.slack) {
                    await postMessage(annoncement, slackAnnounceChannel)
                    console.log(`Posted to Slack channel ${slackAnnounceChannel}`)
                }
            } else {
                console.log('No tracking issues were created, skipping Slack announcement')
            }
        },
    },
    {
        id: 'changelog:cut',
        description: 'Generate pull requests to perform a changelog cut for branch cut',
        argNames: ['changelogFile'],
        run: async (config, changelogFile = 'CHANGELOG.md') => {
            const { upcoming: release } = await releaseVersions(config)
            const prMessage = `changelog: cut sourcegraph@${release.version}`
            const pullRequest = await createChangesets({
                requiredCommands: [],
                changes: [
                    {
                        owner: 'sourcegraph',
                        repo: 'sourcegraph',
                        base: 'main',
                        head: `changelog-${release.version}`,
                        title: prMessage,
                        commitMessage: prMessage + '\n\n ## Test plan\n\nn/a',
                        edits: [
                            (directory: string) => {
                                console.log(`Updating '${changelogFile} for ${release.format()}'`)
                                const changelogPath = path.join(directory, changelogFile)
                                let changelogContents = readFileSync(changelogPath).toString()

                                // Convert 'unreleased' to a release
                                const releaseHeader = `## ${release.format()}`
                                const unreleasedHeader = '## Unreleased'
                                changelogContents = changelogContents.replace(unreleasedHeader, releaseHeader)

                                // Add a blank changelog template for the next release
                                changelogContents = changelogContents.replace(
                                    changelog.divider,
                                    changelog.releaseTemplate
                                )

                                // Update changelog
                                writeFileSync(changelogPath, changelogContents)
                            },
                        ],
                    },
                    {
                        owner: 'sourcegraph',
                        repo: 'deploy-sourcegraph-helm',
                        base: 'main',
                        head: `changelog-${release.version}`,
                        title: prMessage,
                        commitMessage: prMessage,
                        body: prMessage + '\n\n ## Test plan\n\nn/a',
                        edits: [
                            (directory: string) => {
                                console.log(`Updating '${changelogFile} for ${release.format()}'`)
                                const changelogPath = path.join(directory, 'charts', 'sourcegraph', changelogFile)
                                let changelogContents = readFileSync(changelogPath).toString()

                                // Convert 'unreleased' to a release
                                const releaseHeader = `## ${release.format()}`
                                const releaseUpdate =
                                    releaseHeader + `\n\n- Sourcegraph ${release.format()} is now available\n`
                                const unreleasedHeader = '## Unreleased\n'
                                changelogContents = changelogContents.replace(unreleasedHeader, releaseUpdate)

                                // Add a blank changelog template for the next release
                                changelogContents = changelogContents.replace(
                                    changelog.divider,
                                    changelog.simpleReleaseTemplate
                                )

                                // Update changelog
                                writeFileSync(changelogPath, changelogContents)
                            },
                        ],
                    },
                ],
                dryRun: config.dryRun.changesets,
            })
            const changeLogPrUrl = pullRequest[0].pullRequestURL
            console.log(
                `\nPlease review the changelog PR at ${changeLogPrUrl}, and merge manually when checks have passed.`
            )
        },
    },
    {
        id: 'release:branch-cut',
        description: 'Create release branch',
        run: async config => {
            const { upcoming: release } = await releaseVersions(config)
            const branch = `${release.major}.${release.minor}`
            let message: string
            // notify cs team on patch release cut
            if (release.patch !== 0) {
                message = `:mega: *${release.version}* branch has been cut cc: @cs`
            } else {
                message = `:mega: *${release.version}* branch has been cut.`
            }
            try {
                // Create and push new release branch from changelog commit
                await execa('git', ['branch', branch])
                await execa('git', ['push', 'origin', branch])
                await postMessage(message, config.slackAnnounceChannel)
                console.log(`To check the status of the branch, run:\nsg ci status -branch ${release.version} --wait\n`)
            } catch (error) {
                console.error('Failed to create release branch', error)
            }
        },
    },
    {
        id: 'release:status',
        description: 'Post a message in Slack summarizing the progress of a release',
        run: async config => {
            const githubClient = await getAuthenticatedGitHubClient()
            const { upcoming: release } = await releaseVersions(config)

            const trackingIssue = await getTrackingIssue(githubClient, release)
            if (!trackingIssue) {
                throw new Error(`Tracking issue for version ${release.version} not found - has it been created yet?`)
            }
            const latestTag = (await getLatestTag('sourcegraph', 'sourcegraph')).toString()
            const latestBuildURL = `https://buildkite.com/sourcegraph/sourcegraph/builds?branch=${latestTag}`
            const latestBuildMessage = `Latest release build: ${latestTag}. See the build status on <${latestBuildURL}|Buildkite>`
            const blockingQuery = 'is:open org:sourcegraph label:release-blocker'
            const blockingIssues = await listIssues(githubClient, blockingQuery)
            const blockingIssuesURL = `https://github.com/issues?q=${encodeURIComponent(blockingQuery)}`
            const blockingMessage =
                blockingIssues.length === 0
                    ? 'There are no release-blocking issues'
                    : `There ${
                          blockingIssues.length === 1
                              ? 'is 1 release-blocking issue'
                              : `are ${blockingIssues.length} release-blocking issues`
                      }`

            const message = `:mega: *${release.version} Release Status Update*

* Tracking issue: ${trackingIssue.url}
* ${blockingMessage}: ${blockingIssuesURL}
* ${latestBuildMessage}`
            if (!config.dryRun.slack) {
                await postMessage(message, config.slackAnnounceChannel)
            }
        },
    },
    {
        id: 'release:create-candidate',
        description: 'Generate the Nth release candidate. Set <candidate> to "final" to generate a final release',
        run: async config => {
            const { upcoming: release } = await releaseVersions(config)
            const branch = `${release.major}.${release.minor}`
            ensureReleaseBranchUpToDate(branch)

            const owner = 'sourcegraph'
            const repo = 'sourcegraph'

            try {
                const client = await getAuthenticatedGitHubClient()
                const { workdir } = await cloneRepo(client, owner, repo, { revision: branch, revisionMustExist: true })

                const tags = getCandidateTags(workdir, release.version)
                let nextCandidate = 1
                for (const tag of tags) {
                    const num = parseInt(tag.slice(-1), 10)
                    if (num >= nextCandidate) {
                        nextCandidate = num + 1
                    }
                }
                const tag = `v${release.version}-rc.${nextCandidate}`

                console.log(`Detected next candidate: ${nextCandidate}, attempting to create tag: ${tag}`)
                await createTag(
                    client,
                    workdir,
                    {
                        owner,
                        repo,
                        branch,
                        tag,
                    },
                    config.dryRun.tags || false
                )
                console.log(`To check the status of the build, run:\nsg ci status -branch ${tag} --wait\n`)
            } catch (error) {
                console.error('Failed to create tag', error)
            }
        },
    },
    {
        id: 'release:promote-candidate',
        description:
            'Promote a release candidate to release build. Specify the full candidate tag to promote the tagged commit to release.',
        argNames: ['candidate'],
        run: async (config, candidate) => {
            const { upcoming: release } = await releaseVersions(config)
            const releaseBranch = `${release.major}.${release.minor}`
            ensureReleaseBranchUpToDate(releaseBranch)

            const warnMsg =
                'Verify the provided tag is correct to promote to release. Note: it is very unusual to require a non-standard tag to promote to release, proceed with caution.'
            const exampleTag = `v${release.version}-rc.1`
            if (!candidate) {
                throw new Error(
                    `Candidate tag is a required argument. This should be the git tag of the commit to promote to release (ex.${exampleTag}`
                )
            } else if (!candidate.match('v\\d\\.\\d(?:\\.\\d)?-rc\\.\\d')) {
                await verifyWithInput(
                    `Warning!\nCandidate tag: ${candidate} does not match the standard convention (ex. ${exampleTag}). ${warnMsg}`
                )
            } else if (!candidate.match(`${release.version}-rc\\.\\d`)) {
                await verifyWithInput(
                    `Warning!\nCandidate tag: ${candidate} does not match the expected version ${release.version} (ex. ${exampleTag}). ${warnMsg}`
                )
            }

            const owner = 'sourcegraph'
            const repo = 'sourcegraph'

            const releaseTag = `v${release.version}`

            try {
                const client = await getAuthenticatedGitHubClient()
                // passing the tag as branch so that only the specified tag is shallow cloned
                const { workdir } = await cloneRepo(client, owner, repo, {
                    revision: candidate,
                    revisionMustExist: true,
                })

                execa.sync('git', ['fetch', '--tags'], { stdio: 'inherit', cwd: workdir })
                await createTag(
                    client,
                    workdir,
                    {
                        owner,
                        repo,
                        branch: candidate,
                        tag: releaseTag,
                    },
                    config.dryRun.tags || false
                )
                console.log(`To check the status of the build, run:\nsg ci status -branch ${releaseTag} --wait\n`)
            } catch (error) {
                console.error('Failed to create tag', error)
            }
        },
    },
    {
        id: 'release:check-candidate',
        description: 'Check release candidates.',
        argNames: ['version'],
        run: async (config, version) => {
            if (!version) {
                const { upcoming: release } = await releaseVersions(config)
                version = release.version
            }
            const tags = getCandidateTags(localSourcegraphRepo, version)
            console.log(`Release candidate tags for version: ${version}\n${tags}`)
            console.log('To check the status of the build, run:\nsg ci status -branch tag\n')
        },
    },
    {
        id: 'release:stage',
        description: 'Open pull requests and a batch change staging a release',
        run: async config => {
            const { slackAnnounceChannel, dryRun } = config
            const { upcoming: release, previous } = await releaseVersions(config)
            // ensure docker is running for 'batch changes'
            try {
                await ensureDocker()
            } catch (error) {
                console.log(error)
                console.log('Docker required for batch changes')
                process.exit(1)
            }
            // ensure $SRC_ENDPOINT is set
            ensureSrcCliEndpoint()
            // ensure src-cli is up to date
            await ensureSrcCliUpToDate()
            // set up batch change config
            const batchChange = batchChanges.releaseTrackingBatchChange(
                release.version,
                await batchChanges.sourcegraphCLIConfig()
            )

            // default values
            const notPatchRelease = release.patch === 0
            const versionRegex = '[0-9]+\\.[0-9]+\\.[0-9]+'
            const batchChangeURL = batchChanges.batchChangeURL(batchChange)
            const trackingIssue = await getTrackingIssue(await getAuthenticatedGitHubClient(), release)
            if (!trackingIssue) {
                throw new Error(`Tracking issue for version ${release.version} not found - has it been created yet?`)
            }

            // default PR content
            const defaultPRMessage = `release: sourcegraph@${release.version}`
            const prBodyAndDraftState = (
                actionItems: string[],
                customMessage?: string
            ): { draft: boolean; body: string } => {
                const defaultBody = `This pull request is part of the Sourcegraph ${release.version} release.
${customMessage || ''}

* [Release batch change](${batchChangeURL})
* ${trackingIssue ? `[Tracking issue](${trackingIssue.url})` : 'No tracking issue exists for this release'}

### Test plan

CI checks in this repository should pass, and a manual review should confirm if the generated changes are correct.`

                if (!actionItems || actionItems.length === 0) {
                    return { draft: false, body: defaultBody }
                }
                return {
                    draft: true, // further actions required before merge
                    body: `${defaultBody}

### :warning: Additional changes required

These steps must be completed before this PR can be merged, unless otherwise stated. Push any required changes directly to this PR branch.

${actionItems.map(item => `- [ ] ${item}`).join('\n')}

cc @${config.captainGitHubUsername}

`,
                }
            }

            // Render changes
            const createdChanges = await createChangesets({
                requiredCommands: ['comby', sed, 'find', 'go', 'src', 'sg'],
                changes: [
                    {
                        owner: 'sourcegraph',
                        repo: 'sourcegraph',
                        base: 'main',
                        head: `publish-${release.version}`,
                        commitMessage: notPatchRelease
                            ? `draft sourcegraph@${release.version} release`
                            : defaultPRMessage,
                        title: defaultPRMessage,
                        edits: [
                            // Update references to Sourcegraph versions in docs
                            `${sed} -i -E 's/version \`${versionRegex}\`/version \`${release.version}\`/g' doc/index.md`,
                            // Update sourcegraph/server:VERSION everywhere except changelog
                            `find . -type f -name '*.md' ! -name 'CHANGELOG.md' -exec ${sed} -i -E 's/sourcegraph\\/server:${versionRegex}/sourcegraph\\/server:${release.version}/g' {} +`,
                            // Update Sourcegraph versions in installation guides
                            `find ./doc/admin/deploy/ -type f -name '*.md' -exec ${sed} -i -E 's/SOURCEGRAPH_VERSION="v${versionRegex}"/SOURCEGRAPH_VERSION="v${release.version}"/g' {} +`,
                            `find ./doc/admin/deploy/ -type f -name '*.md' -exec ${sed} -i -E 's/--version ${versionRegex}/--version ${release.version}/g' {} +`,
                            // Update fork variables in installation guides
                            `find ./doc/admin/deploy/ -type f -name '*.md' -exec ${sed} -i -E "s/DEPLOY_SOURCEGRAPH_DOCKER_FORK_REVISION='v${versionRegex}'/DEPLOY_SOURCEGRAPH_DOCKER_FORK_REVISION='v${release.version}'/g" {} +`,

                            notPatchRelease
                                ? `comby -in-place '{{$previousReleaseRevspec := ":[1]"}} {{$previousReleaseVersion := ":[2]"}} {{$currentReleaseRevspec := ":[3]"}} {{$currentReleaseVersion := ":[4]"}}' '{{$previousReleaseRevspec := ":[3]"}} {{$previousReleaseVersion := ":[4]"}} {{$currentReleaseRevspec := "v${release.version}"}} {{$currentReleaseVersion := "${release.major}.${release.minor}"}}' doc/_resources/templates/document.html`
                                : `comby -in-place 'currentReleaseRevspec := ":[1]"' 'currentReleaseRevspec := "v${release.version}"' doc/_resources/templates/document.html`,

                            // Update references to Sourcegraph deployment versions
                            `comby -in-place 'latestReleaseKubernetesBuild = newBuild(":[1]")' "latestReleaseKubernetesBuild = newBuild(\\"${release.version}\\")" cmd/frontend/internal/app/updatecheck/handler.go`,
                            `comby -in-place 'latestReleaseDockerServerImageBuild = newBuild(":[1]")' "latestReleaseDockerServerImageBuild = newBuild(\\"${release.version}\\")" cmd/frontend/internal/app/updatecheck/handler.go`,
                            `comby -in-place 'latestReleaseDockerComposeOrPureDocker = newBuild(":[1]")' "latestReleaseDockerComposeOrPureDocker = newBuild(\\"${release.version}\\")" cmd/frontend/internal/app/updatecheck/handler.go`,

                            // Support current release as the "previous release" going forward
                            notPatchRelease
                                ? `comby -in-place 'const minimumUpgradeableVersion = ":[1]"' 'const minimumUpgradeableVersion = "${release.version}"' enterprise/dev/ci/internal/ci/*.go`
                                : 'echo "Skipping minimumUpgradeableVersion bump on patch release"',
                            updateUpgradeGuides(previous.version, release.version),
                        ],
                        ...prBodyAndDraftState(
                            ((): string[] => {
                                const items: string[] = []
                                items.push(
                                    'Ensure all other pull requests in the batch change have been merged',
                                    'Run `pnpm run release release:finalize` to generate the tags required. CI will not pass until this command is run.',
                                    'Re-run the build on this branch (using either the command `sg ci build` or the Buildkite UI) and merge when the build passes.'
                                )
                                return items
                            })()
                        ),
                    },
                    {
                        owner: 'sourcegraph',
                        repo: 'about',
                        base: 'main',
                        head: `publish-${release.version}`,
                        commitMessage: defaultPRMessage,
                        title: defaultPRMessage,
                        edits: [
                            // Update sourcegraph/server:VERSION in all tsx files
                            `find . -type f -name '*.tsx' -exec ${sed} -i -E 's/sourcegraph\\/server:${versionRegex}/sourcegraph\\/server:${release.version}/g' {} +`,
                        ],
                        ...prBodyAndDraftState(
                            [],
                            notPatchRelease ? 'Note that this PR does *not* include the release blog post.' : undefined
                        ),
                    },
                    {
                        owner: 'sourcegraph',
                        repo: 'deploy-sourcegraph',
                        base: `${release.major}.${release.minor}`,
                        head: `publish-${release.version}`,
                        commitMessage: defaultPRMessage,
                        title: defaultPRMessage,
                        edits: [`tools/update-docker-tags.sh ${release.version}`],
                        ...prBodyAndDraftState([]),
                    },
                    {
                        owner: 'sourcegraph',
                        repo: 'deploy-sourcegraph-k8s',
                        base: `${release.major}.${release.minor}`,
                        head: `publish-${release.version}`,
                        commitMessage: defaultPRMessage,
                        title: defaultPRMessage,
                        edits: [`sg ops update-images -pin-tag ${release.version} base/`],
                        ...prBodyAndDraftState([]),
                    },
                    {
                        owner: 'sourcegraph',
                        repo: 'deploy-sourcegraph-docker',
                        base: `${release.major}.${release.minor}`,
                        head: `publish-${release.version}`,
                        commitMessage: defaultPRMessage,
                        title: defaultPRMessage,
                        edits: [`tools/update-docker-tags.sh ${release.version}`],
                        ...prBodyAndDraftState([]),
                    },
                    {
                        owner: 'sourcegraph',
                        repo: 'deploy-sourcegraph-docker-customer-replica-1',
                        base: `${release.major}.${release.minor}`,
                        head: `publish-${release.version}`,
                        commitMessage: defaultPRMessage,
                        title: defaultPRMessage,
                        edits: [`tools/update-docker-tags.sh ${release.version}`],
                        ...prBodyAndDraftState([]),
                    },
                    {
                        owner: 'sourcegraph',
                        repo: 'deploy-sourcegraph-aws',
                        base: 'master',
                        head: `publish-${release.version}`,
                        commitMessage: defaultPRMessage,
                        title: defaultPRMessage,
                        edits: [
                            `${sed} -i -E 's/export SOURCEGRAPH_VERSION=${versionRegex}/export SOURCEGRAPH_VERSION=${release.version}/g' resources/amazon-linux2.sh`,
                        ],
                        ...prBodyAndDraftState([]),
                    },
                    {
                        owner: 'sourcegraph',
                        repo: 'deploy-sourcegraph-digitalocean',
                        base: 'master',
                        head: `publish-${release.version}`,
                        commitMessage: defaultPRMessage,
                        title: defaultPRMessage,
                        edits: [
                            `${sed} -i -E 's/export SOURCEGRAPH_VERSION=${versionRegex}/export SOURCEGRAPH_VERSION=${release.version}/g' resources/user-data.sh`,
                        ],
                        ...prBodyAndDraftState([]),
                    },
                    {
                        owner: 'sourcegraph',
                        repo: 'deploy-sourcegraph-helm',
                        base: `release/${release.major}.${release.minor}`,
                        head: `publish-${release.version}`,
                        commitMessage: defaultPRMessage,
                        title: defaultPRMessage,
                        edits: [
                            `for i in charts/*; do sg ops update-images -kind helm -pin-tag ${release.version} $i/.; done`,
                            `${sed} -i 's/appVersion:.*/appVersion: "${release.version}"/g' charts/*/Chart.yaml`,
                            `${sed} -i 's/version:.*/version: "${release.version}"/g' charts/*/Chart.yaml`,
                            './scripts/helm-docs.sh',
                        ],
                        ...prBodyAndDraftState([]),
                    },
                ],
                dryRun: dryRun.changesets,
            })

            // if changesets were actually published, set up a batch change and post in Slack
            if (!dryRun.changesets) {
                // Create batch change to track changes
                try {
                    console.log(`Creating batch change in ${batchChange.cliConfig.SRC_ENDPOINT}`)
                    await batchChanges.createBatchChange(
                        createdChanges,
                        batchChange,
                        `Track publishing of sourcegraph v${release.version}: ${trackingIssue?.url}`
                    )
                } catch (error) {
                    console.error(error)
                    console.error('Failed to create batch change for this release, continuing with announcement')
                }

                // Announce release update in Slack
                if (!dryRun.slack) {
                    await postMessage(
                        `:captain: *Sourcegraph ${release.version} has been staged.*

Batch change: ${batchChangeURL}`,
                        slackAnnounceChannel
                    )
                }
            }
        },
    },
    {
        id: 'release:add-to-batch-change',
        description: 'Manually add a change to a release batch change',
        argNames: ['changeRepo', 'changeID'],
        // Example: pnpm run release release:add-to-batch-change sourcegraph/about 1797
        run: async (config, changeRepo, changeID) => {
            const { upcoming: release } = await releaseVersions(config)
            if (!changeRepo || !changeID) {
                throw new Error('Missing parameters (required: version, repo, change ID)')
            }

            const batchChange = batchChanges.releaseTrackingBatchChange(
                release.version,
                await batchChanges.sourcegraphCLIConfig()
            )
            await batchChanges.addToBatchChange(
                [
                    {
                        repository: changeRepo,
                        pullRequestNumber: parseInt(changeID, 10),
                    },
                ],
                batchChange
            )
            console.log(`Added ${changeRepo}#${changeID} to batch change ${batchChanges.batchChangeURL(batchChange)}`)
        },
    },
    {
        id: 'release:finalize',
        description: 'Run final tasks for sourcegraph/sourcegraph release pull requests',
        run: async config => {
            const { upcoming: release } = await releaseVersions(config)
            let failed = false

            const owner = 'sourcegraph'
            // Push final tags
            const branch = `${release.major}.${release.minor}`
            const tag = `v${release.version}`
            for (const repo of [
                'deploy-sourcegraph',
                'deploy-sourcegraph-docker',
                'deploy-sourcegraph-docker-customer-replica-1',
                'deploy-sourcegraph-k8s',
            ]) {
                try {
                    const client = await getAuthenticatedGitHubClient()
                    const { workdir } = await cloneRepo(client, owner, repo, {
                        revision: branch,
                        revisionMustExist: true,
                    })
                    await createTag(
                        client,
                        workdir,
                        {
                            owner,
                            repo,
                            branch,
                            tag,
                        },
                        config.dryRun.tags || false
                    )
                } catch (error) {
                    console.error(error)
                    console.error(`Failed to create tag ${tag} on ${repo}@${branch}`)
                    failed = true
                }
            }

            if (failed) {
                throw new Error('Error occured applying some changes - please check log output')
            }
        },
    },
    {
        id: 'release:announce',
        description: 'Announce a release as live',
        run: async config => {
            const { slackAnnounceChannel, dryRun } = config
            const { upcoming: release } = await releaseVersions(config)
            const githubClient = await getAuthenticatedGitHubClient()

            // Create final GitHub release
            let githubRelease = ''
            try {
                githubRelease = await createLatestRelease(
                    githubClient,
                    {
                        owner: 'sourcegraph',
                        repo: 'sourcegraph',
                        release,
                    },
                    dryRun.tags
                )
            } catch (error) {
                console.error('Failed to generate GitHub release:', error)
                // Do not block process
            }

            // Set up announcement message
            const batchChangeURL = batchChanges.batchChangeURL(
                batchChanges.releaseTrackingBatchChange(release.version, await batchChanges.sourcegraphCLIConfig())
            )
            const releaseMessage = `*Sourcegraph ${release.version} has been published*

* Changelog: ${changelogURL(release.format())}
* GitHub release: ${githubRelease || 'No release generated'}
* Release batch change: ${batchChangeURL}`

            // Slack
            const slackMessage = `:captain: ${releaseMessage}`
            if (!dryRun.slack) {
                await postMessage(slackMessage, slackAnnounceChannel)
                console.log(`Posted to Slack channel ${slackAnnounceChannel}`)
            } else {
                console.log(`dryRun enabled, skipping Slack post to ${slackAnnounceChannel}: ${slackMessage}`)
            }

            // GitHub tracking issues
            const trackingIssue = await getTrackingIssue(githubClient, release)
            if (!trackingIssue) {
                console.warn(`Could not find tracking issue for release ${release.version} - skipping`)
            } else {
                // Note patch release requests if there are any outstanding
                let comment = `${releaseMessage}

@${config.captainGitHubUsername}: Please complete the post-release steps before closing this issue.`
                const patchRequestIssues = await queryIssues(githubClient, '', [IssueLabel.PATCH_REQUEST])
                if (patchRequestIssues && patchRequestIssues.length > 0) {
                    comment += `
Please also update outstanding patch requests, if relevant:

${patchRequestIssues.map(issue => `* #${issue.number}`).join('\n')}`
                }
                if (!dryRun.trackingIssues) {
                    const commentURL = await commentOnIssue(githubClient, trackingIssue, comment)
                    console.log(`Please make sure to follow up on the release issue: ${commentURL}`)
                } else {
                    console.log(`dryRun enabled, skipping GitHub comment to ${trackingIssue.url}: ${comment}`)
                }
            }
        },
    },
    {
        id: 'release:close',
        description: 'Close tracking issues for current release',
        run: async config => {
            const { previous: release } = await releaseVersions(config)
            // close tracking issue
            await closeTrackingIssue(release)
        },
    },
    {
        id: 'release:prepare',
        description: 'Prepare a feature release',
        run: async config => {
            const rconfig = loadReleaseConfig()
            const rel = await newReleaseFromInput()
            addRelease(rconfig, rel)
            saveReleaseConfig(rconfig)
        },
    },
    {
        id: 'release:remove',
        description: 'Remove a release from the config',
        argNames: ['version'],
        run: async (config, version) => {
            await verifyWithInput(`Confirm you want to remove release: ${version} from the release config?`)
            const rconfig = loadReleaseConfig()
            removeRelease(rconfig, version)
            saveReleaseConfig(rconfig)
        },
    },
    {
        id: 'util:clear-cache',
        description: 'Clear release tool cache',
        run: () => {
            rmdirSync(cacheFolder, { recursive: true })
        },
    },
    {
        id: '_test:release-guide-content',
        description: 'Generate upgrade guides',
        argNames: ['previous', 'next'],
        run: (config, previous, next) => {
            for (const content of getAllUpgradeGuides(previous, next)) {
                console.log(content)
            }
        },
    },
    {
        id: '_test:release-guide-update',
        description: 'Test update the upgrade guides',
        argNames: ['previous', 'next', 'dir'],
        run: (config, previous, next, dir) => {
            updateUpgradeGuides(previous, next)(dir)
        },
    },
    {
        id: '_test:google-calendar',
        description: 'Test Google Calendar integration',
        run: async config => {
            const googleCalendar = await getClient()
            await ensureEvent(
                {
                    title: 'TEST EVENT',
                    startDateTime: new Date(config.releaseDate).toISOString(),
                    endDateTime: addMinutes(new Date(config.releaseDate), 1).toISOString(),
                    transparency: 'transparent',
                },
                googleCalendar
            )
        },
    },
    {
        id: '_test:slack',
        description: 'Test Slack integration',
        argNames: ['channel', 'message'],
        run: async ({ dryRun }, channel, message) => {
            if (!dryRun.slack) {
                await postMessage(message, channel)
            }
        },
    },
    {
        id: '_test:batchchange-create-from-changes',
        description: 'Test batch changes integration',
        argNames: ['batchchangeConfigJSON'],
        // Example: pnpm run release _test:batchchange-create-from-changes "$(cat ./.secrets/test-batch-change-import.json)"
        run: async (_config, batchchangeConfigJSON) => {
            const batchChangeConfig = JSON.parse(batchchangeConfigJSON) as {
                changes: CreatedChangeset[]
                name: string
                description: string
            }

            // set up src-cli
            await commandExists('src')
            const batchChange = {
                name: batchChangeConfig.name,
                description: batchChangeConfig.description,
                namespace: 'sourcegraph',
                cliConfig: await batchChanges.sourcegraphCLIConfig(),
            }

            await batchChanges.createBatchChange(
                batchChangeConfig.changes,
                batchChange,
                'release tool testing batch change'
            )
            console.log(`Created batch change ${batchChanges.batchChangeURL(batchChange)}`)
        },
    },
    {
        id: '_test:config',
        description: 'Test release configuration loading',
        run: async config => {
           const rconfig = loadReleaseConfig()
            const rel = await newReleaseFromInput()
            rconfig.releases[rel.current] = rel
            saveReleaseConfig(rconfig)
        },
    },
    {
        id: '_test:dockerensure',
        description: 'test docker ensure function',
        run: async () => {
            try {
                await ensureDocker()
            } catch (error) {
                console.log(error)
                console.log('Docker required for batch changes')
                process.exit(1)
            }
        },
    },
    {
        id: '_test:srccliensure',
        description: 'test srccli version',
        run: async () => {
            ensureSrcCliEndpoint()
            await ensureSrcCliUpToDate()
        },
    },
]
