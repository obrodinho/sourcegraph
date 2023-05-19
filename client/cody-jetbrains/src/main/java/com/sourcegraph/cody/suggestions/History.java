package com.sourcegraph.cody.suggestions;

import com.intellij.openapi.Disposable;

import java.util.ArrayList;
import java.util.List;

public class History implements Disposable {
  List<Disposable> subscriptions = new ArrayList<>();

  @Override
  public void dispose() {}
}
