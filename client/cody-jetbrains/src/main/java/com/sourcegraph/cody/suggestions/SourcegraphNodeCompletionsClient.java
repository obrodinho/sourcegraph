package com.sourcegraph.cody.suggestions;

import com.google.gson.Gson;
import com.google.gson.GsonBuilder;
import com.sourcegraph.cody.completions.CompletionsCallbacks;
import com.sourcegraph.cody.completions.CompletionsInput;
import com.sourcegraph.cody.completions.CompletionsService;
import com.sourcegraph.cody.completions.Message;
import com.sourcegraph.cody.vscode.CancellationToken;

import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.stream.Collectors;

public class SourcegraphNodeCompletionsClient {
  public final CompletionsService completionsService;

  public SourcegraphNodeCompletionsClient(CompletionsService completionsService) {
    this.completionsService = completionsService;
  }

  public CompletableFuture<CompletionResponse> complete(
      CompletionParameters params, CancellationToken abortSignal) {
    SuggestionsCallbacks callbacks = new SuggestionsCallbacks();
    //    System.out.println(
    //        "QUERY: " +
    // params.messages.stream().map(Message::prompt).collect(Collectors.joining("")));
    completionsService.streamCompletion(
        new CompletionsInput(
            params.messages,
            params.temperature,
            params.stopSequences,
            params.maxTokensToSample,
            params.topK,
            params.topP),
        callbacks,
        CompletionsService.Endpoint.Code);
    return callbacks.promise;
  }

  private static class SuggestionsCallbacks implements CompletionsCallbacks {
    CompletableFuture<CompletionResponse> promise = new CompletableFuture<>();
    List<String> chunks = new ArrayList<>();

    @Override
    public void onSubscribed() {
      // Do nothing
    }

    @Override
    public void onData(String data) {
      chunks.add(data);
    }

    @Override
    public void onError(Throwable error) {
      error.printStackTrace();
      promise.completeExceptionally(error);
    }

    @Override
    public void onComplete() {
      String json = String.join("", chunks);
      CompletionResponse completionResponse = new Gson().fromJson(json, CompletionResponse.class);
      promise.complete(completionResponse);
    }
  }
}
