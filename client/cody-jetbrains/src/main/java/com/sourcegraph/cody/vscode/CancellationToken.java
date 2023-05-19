package com.sourcegraph.cody.vscode;

import java.util.concurrent.CompletableFuture;
import java.util.function.Consumer;

public class CancellationToken {
  public final CompletableFuture<Boolean> cancelled = new CompletableFuture<>();

  public void onCancellationRequested(Runnable callback) {
    this.cancelled.thenAccept(
        (cancelled) -> {
          if (cancelled) {
            callback.run();
          }
        });
  }

  public boolean isCancelled() {
    return this.cancelled.isDone();
  }

  public void abort() {
    this.cancelled.complete(true);
  }
}
