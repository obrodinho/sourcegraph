package com.sourcegraph.cody.completions;

public enum Speaker {
  HUMAN,
  ASSISTANT;

  public String prompt() {
    if (this == HUMAN) return "\n\nHuman:";
    return "\n\nAssistant:";
  }
}
