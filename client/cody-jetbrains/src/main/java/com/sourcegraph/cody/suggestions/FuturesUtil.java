package com.sourcegraph.cody.suggestions;

import java.util.ArrayList;
import java.util.List;
import java.util.Queue;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentLinkedQueue;

public class FuturesUtil {
  public static <T> CompletableFuture<List<T>> all(List<CompletableFuture<T>> promises) {
    Queue<T> responses = new ConcurrentLinkedQueue<>();
    for (CompletableFuture<T> promise : promises) {
      promise.thenAccept(responses::add);
    }
    return CompletableFuture.allOf(promises.toArray(new CompletableFuture[0]))
        .thenApply((ignore) -> new ArrayList<>(responses));
  }
}
