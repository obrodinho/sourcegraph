package com.sourcegraph.cody.suggestions;

public class JaccardMatch {
    public final int score;
    public final String text;

    public JaccardMatch(int score, String text) {
        this.score = score;
        this.text = text;
    }
}
