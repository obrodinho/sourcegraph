export interface ActiveTextEditor {
    content: string
    filePath: string
    repoName?: string
    revision?: string
}

export interface ActiveTextEditorSelection {
    fileName: string
    repoName?: string
    revision?: string
    precedingText: string
    selectedText: string
    followingText: string
}

export interface ActiveTextEditorVisibleContent {
    content: string
    fileName: string
    repoName?: string
    revision?: string
}

interface VsCodeInlineController {
    selection: ActiveTextEditorSelection | null
}

interface VsCodeTaskContoller {
    add(input: string, selection: ActiveTextEditorSelection): string | null
    stop(taskID: string): void
}

export interface ActiveTextEditorViewControllers {
    inline: VsCodeInlineController
    task: VsCodeTaskContoller
}

export interface Editor {
    controllers?: ActiveTextEditorViewControllers
    getWorkspaceRootPath(): string | null
    getActiveTextEditor(): ActiveTextEditor | null
    getActiveTextEditorSelection(): ActiveTextEditorSelection | null

    /**
     * Gets the active text editor's selection, or the entire file if the selected range is empty.
     */
    getActiveTextEditorSelectionOrEntireFile(): ActiveTextEditorSelection | null

    getActiveTextEditorVisibleContent(): ActiveTextEditorVisibleContent | null
    replaceSelection(fileName: string, selectedText: string, replacement: string): Promise<void>
    showQuickPick(labels: string[]): Promise<string | undefined>
    showWarningMessage(message: string): Promise<void>
    showInputBox(prompt?: string): Promise<string | undefined>
}
