package com.sourcegraph.cody.vscode;

import com.intellij.openapi.editor.Editor;
import com.sourcegraph.cody.suggestions.InlineCompletionList;

import java.util.concurrent.CompletableFuture;

public abstract class InlineCompletionItemProvider {
  public abstract CompletableFuture<InlineCompletionList> provideInlineCompletions(
      TextDocument document,
      Position position,
      InlineCompletionContext context,
      CancellationToken token);
}
