package com.sourcegraph.cody.suggestions;

import com.intellij.openapi.Disposable;
import com.intellij.openapi.application.ApplicationManager;
import com.intellij.openapi.command.CommandProcessor;
import com.intellij.openapi.editor.Editor;
import com.intellij.openapi.editor.event.*;
import com.intellij.openapi.editor.ex.util.EditorUtil;
import com.intellij.openapi.fileEditor.FileEditor;
import com.intellij.openapi.fileEditor.FileEditorManager;
import com.intellij.openapi.fileEditor.TextEditor;
import com.intellij.openapi.fileEditor.impl.FileEditorManagerImpl;
import com.intellij.openapi.project.Project;
import com.intellij.openapi.util.Disposer;
import com.sourcegraph.cody.vscode.InlineCompletionTriggerKind;
import org.jetbrains.annotations.NotNull;

public class CodyEditorFactoryListener implements EditorFactoryListener {
  CodySelectionListener selectionListener = new CodySelectionListener();

  @Override
  public void editorCreated(@NotNull EditorFactoryEvent event) {
    Editor editor = event.getEditor();
    System.out.println("CODY: editorCreated() " + editor.getDocument());
    Project project = editor.getProject();
    if (project == null || project.isDisposed()) {
      return;
    }
    Disposable disposable = Disposer.newDisposable("codyEditorListener");
    EditorUtil.disposeWithEditor(editor, disposable);
    editor.getCaretModel().addCaretListener(new CodyCaretListener(editor), disposable);
    editor.getSelectionModel().addSelectionListener(this.selectionListener, disposable);
    editor.getDocument().addDocumentListener(new CodyDocumentListener(editor), disposable);
  }

  private static class CodyCaretListener implements CaretListener {
    private final Editor editor;

    public CodyCaretListener(Editor editor) {
      this.editor = editor;
    }

    @Override
    public void caretPositionChanged(@NotNull CaretEvent e) {
      System.out.println("CODY: caretPositionChanged()");
      CodySuggestionsExecutor suggestions = CodySuggestionsExecutor.getInstance();
      if (suggestions.isEnabledForEditor(e.getEditor())
          && CodyEditorFactoryListener.isSelectedEditor(e.getEditor())) {
        if (suggestions.hasCompletionInlays(this.editor)) {
          suggestions.execute(
              this.editor,
              editor.getCaretModel().getOffset(),
              InlineCompletionTriggerKind.Automatic);
        } else {
          suggestions.disposeInlays(this.editor);
        }
      }
    }
  }

  private static class CodySelectionListener implements SelectionListener {
    @Override
    public void selectionChanged(@NotNull SelectionEvent e) {
      if (CodySuggestionsExecutor.getInstance().isEnabledForEditor(e.getEditor())
          && CodyEditorFactoryListener.isSelectedEditor(e.getEditor())) {
        ApplicationManager.getApplication()
            .getService(CodySuggestionsExecutor.class)
            .disposeInlays(e.getEditor());
      }
    }
  }

  private static class CodyDocumentListener implements BulkAwareDocumentListener {
    private final Editor editor;

    public CodyDocumentListener(Editor editor) {
      this.editor = editor;
    }

    public void documentChangedNonBulk(@NotNull DocumentEvent event) {
      if (!CodyEditorFactoryListener.isSelectedEditor(this.editor)) {
        return;
      }
      CodySuggestionsExecutor editorManager = CodySuggestionsExecutor.getInstance();
      if (editorManager.isEnabledForEditor(this.editor)) {
        CommandProcessor commandProcessor = CommandProcessor.getInstance();
        if (!commandProcessor.isUndoTransparentActionInProgress()) {
          int changeOffset = event.getOffset() + event.getNewLength();
          if (this.editor.getCaretModel().getOffset() == changeOffset) {
            InlineCompletionTriggerKind requestType =
                event.getOldLength() != event.getNewLength()
                    ? InlineCompletionTriggerKind.Invoke
                    : InlineCompletionTriggerKind.Automatic;
            editorManager.execute(this.editor, changeOffset, requestType);
          }
        }
      }
    }
  }

  private static boolean isSelectedEditor(Editor editor) {
    if (editor == null) {
      return false;
    }
    Project project = editor.getProject();
    if (project == null || project.isDisposed()) {
      return false;
    }
    FileEditorManager editorManager = FileEditorManager.getInstance(project);
    if (editorManager == null) {
      return false;
    }
    if (editorManager instanceof FileEditorManagerImpl) {
      Editor current = ((FileEditorManagerImpl) editorManager).getSelectedTextEditor(true);
      return current != null && current.equals(editor);
    }
    FileEditor current = editorManager.getSelectedEditor();
    return current instanceof TextEditor && editor.equals(((TextEditor) current).getEditor());
  }
}
