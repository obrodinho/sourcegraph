# This Dockerfile was generated from github.com/sourcegraph/godockerize. It
# was not written by a human, and as such looks janky. As you change this
# file, please don't be scared to make it more pleasant / remove hadolint
# ignores.

FROM sourcegraph/alpine-3.14:213466_2023-04-17_5.0-bdda34a71619@sha256:6354a4ff578b685e36c8fbde81f62125ae0011b047fb2cc22d1b0de616b3c59a

ARG COMMIT_SHA="unknown"
ARG DATE="unknown"
ARG VERSION="unknown"

LABEL org.opencontainers.image.revision=${COMMIT_SHA}
LABEL org.opencontainers.image.created=${DATE}
LABEL org.opencontainers.image.version=${VERSION}
LABEL com.sourcegraph.github.url=https://github.com/sourcegraph/sourcegraph/commit/${COMMIT_SHA}

# Firecracker doesn't work in docker, so disable it by default.
ENV EXECUTOR_USE_FIRECRACKER=false

# Preconfigure bundled-executor to take 1 parallel job and restart afterwards, this is to keep the
# environment clean-ish.
ENV EXECUTOR_MAXIMUM_NUM_JOBS=1
ENV EXECUTOR_NUM_TOTAL_JOBS=1

# Install git and docker. We use the same version here as we use in gitserver.
RUN apk add --no-cache \
    # Minimal version requirement to address vulnerabilities
    # https://github.blog/2023-02-14-git-security-vulnerabilities-announced-3/
    # Don't use alpine/edge, the git release on this segfaults
    'git>=2.38.0' --repository=http://dl-cdn.alpinelinux.org/alpine/v3.17/main \
    ca-certificates

# Install src-cli.
ARG SRC_CLI_VERSION
RUN set -ex && \
    curl -f -L -o src-cli.tar.gz "https://github.com/sourcegraph/src-cli/releases/download/${SRC_CLI_VERSION}/src-cli_${SRC_CLI_VERSION}_linux_amd64.tar.gz" && \
    tar -xvzf src-cli.tar.gz src && \
    mv src /usr/local/bin/src && \
    chmod +x /usr/local/bin/src && \
    rm -rf src-cli.tar.gz

# Install additional common tools for running batch changes.
RUN apk add --no-cache \
    xmlstarlet \
    python3 py3-pip

# Install batcheshelper.
COPY batcheshelper /usr/local/bin/

USER sourcegraph
ENTRYPOINT ["/sbin/tini", "--", "/usr/local/bin/executor"]
COPY executor /usr/local/bin/
