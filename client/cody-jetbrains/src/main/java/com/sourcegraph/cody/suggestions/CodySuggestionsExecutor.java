package com.sourcegraph.cody.suggestions;

import com.intellij.injected.editor.EditorWindow;
import com.intellij.openapi.application.ApplicationManager;
import com.intellij.openapi.editor.Editor;
import com.intellij.openapi.editor.EditorCustomElementRenderer;
import com.intellij.openapi.editor.Inlay;
import com.intellij.openapi.editor.InlayModel;
import com.intellij.openapi.editor.ex.EditorEx;
import com.intellij.openapi.editor.impl.ImaginaryEditor;
import com.intellij.openapi.project.Project;
import com.intellij.openapi.util.Key;
import com.intellij.util.concurrency.annotations.RequiresEdt;
import com.sourcegraph.cody.CodyCompatibility;
import com.sourcegraph.cody.completions.CompletionsService;
import com.sourcegraph.cody.config.ConfigUtil;
import com.sourcegraph.cody.config.SettingsComponent;
import com.sourcegraph.cody.vscode.*;
import org.jetbrains.annotations.NotNull;

import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

/** Responsible for executing suggestions. */
public class CodySuggestionsExecutor {
  private static final Key<Boolean> KEY_EDITOR_SUPPORTED =
      Key.create("sourcegraph.cody.editorAvailable");
  private final ExecutorService executor = Executors.newSingleThreadExecutor();

  public static @NotNull CodySuggestionsExecutor getInstance() {
    return ApplicationManager.getApplication().getService(CodySuggestionsExecutor.class);
  }

  @RequiresEdt
  public void disposeInlays(Editor editor) {
    System.out.println("TODO: dispose inlay");
  }

  @RequiresEdt
  public boolean isEnabledForEditor(Editor editor) {
    return editor != null && isProjectAvailable(editor.getProject()) && isEditorSupported(editor);
  }

  public void execute(Editor editor, int offset, InlineCompletionTriggerKind requestType) {
    SourcegraphNodeCompletionsClient client =
        new SourcegraphNodeCompletionsClient(completionsService(editor));
    CodyCompletionItemProvider provider =
        new CodyCompletionItemProvider(
            new WebviewErrorMessenger(),
            client,
            new CompletionsDocumentProvider(),
            new History(),
            2048,
            4,
            200,
            0.6,
            0.1);
    TextDocument textDocument = new EditorTextDocument(editor);
    CancellationToken token = new CancellationToken();
    System.out.println("<<<<<<<<------------- COMPLETION -------------->>>>>>>>>>");
    this.executor.submit(
        () ->
            provider
                .provideInlineCompletions(
                    textDocument,
                    textDocument.positionAt(offset),
                    new InlineCompletionContext(
                        com.sourcegraph.cody.vscode.InlineCompletionTriggerKind.Automatic, null),
                    token)
                .thenAccept(
                    result -> {
                      System.out.println("=============> COMPLETION: " + result.items.size());
                      if (result.items.isEmpty()) {
                        return;
                      }
                      InlayModel inlayModel = editor.getInlayModel();
                      InlineCompletionItem item = result.items.get(0);
                      if (item.insertText.isEmpty()) {
                        return;
                      }
                      System.out.println("line: " + item.insertText);
                      try {
                        EditorCustomElementRenderer renderer =
                            new SimpleEditorCustomElementRenderer(item.insertText, editor);
                        ApplicationManager.getApplication()
                            .invokeLater(
                                () -> {
                                  Inlay<EditorCustomElementRenderer> inlay =
                                      inlayModel.addInlineElement(offset, true, renderer);
                                  System.out.println("inlay: " + inlay);
                                });
                      } catch (Exception e) {
                        e.printStackTrace();
                      }
                    }));
  }

  private boolean isProjectAvailable(Project project) {
    return project != null && !project.isDisposed();
  }

  private boolean isEditorSupported(Editor editor) {
    if (editor.isDisposed()) {
      return false;
    }

    Boolean fromCache = KEY_EDITOR_SUPPORTED.get(editor);
    if (fromCache != null) {
      return fromCache;
    }

    boolean isSupported =
        isEditorInstanceSupported(editor)
            && CodyCompatibility.isSupportedProject(editor.getProject());
    KEY_EDITOR_SUPPORTED.set(editor, isSupported);
    return isSupported;
  }

  private boolean isEditorInstanceSupported(Editor editor) {
    return !(editor instanceof EditorWindow)
        && !(editor instanceof ImaginaryEditor)
        && (!(editor instanceof EditorEx) || !((EditorEx) editor).isEmbeddedIntoDialogWrapper())
        && !editor.isViewer()
        && !editor.isOneLineMode();
  }

  private CompletionsService completionsService(Editor editor) {
    Project project = editor.getProject();
    boolean isEnterprise =
        ConfigUtil.getInstanceType(project).equals(SettingsComponent.InstanceType.ENTERPRISE);
    String srcEndpoint = System.getenv("SRC_ENDPOINT");
    String instanceUrl =
        srcEndpoint != null
            ? srcEndpoint
            : isEnterprise ? ConfigUtil.getEnterpriseUrl(project) : "https://sourcegraph.com/";
    String accessToken =
        isEnterprise
            ? ConfigUtil.getEnterpriseAccessToken(project)
            : ConfigUtil.getDotcomAccessToken(project);
    if (!instanceUrl.endsWith("/")) {
      instanceUrl = instanceUrl + "/";
    }
    if (accessToken == null) {
      accessToken = System.getenv("SRC_ACCESS_TOKEN");
    }
    if (accessToken == null) {
      throw new IllegalArgumentException("ACCESS_TOKEN is null");
    }
    return new CompletionsService(instanceUrl, accessToken);
  }

  public boolean hasCompletionInlays(Editor editor) {
    System.out.println("TODO: hasCompletionInlays()");
    return true;
  }
}
