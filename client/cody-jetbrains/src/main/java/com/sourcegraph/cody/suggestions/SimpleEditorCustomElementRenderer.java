package com.sourcegraph.cody.suggestions;

import com.intellij.openapi.editor.DefaultLanguageHighlighterColors;
import com.intellij.openapi.editor.Editor;
import com.intellij.openapi.editor.EditorCustomElementRenderer;
import com.intellij.openapi.editor.Inlay;
import com.intellij.openapi.editor.colors.EditorFontType;
import com.intellij.openapi.editor.impl.EditorImpl;
import com.intellij.openapi.editor.markup.TextAttributes;
import com.intellij.ui.JBColor;
import org.jetbrains.annotations.NotNull;

import java.awt.*;

public class SimpleEditorCustomElementRenderer implements EditorCustomElementRenderer {
  private final String text;
  private final TextAttributes themeAttributes;
  private final Editor editor;

  public SimpleEditorCustomElementRenderer(String text, Editor editor) {
    this.text = text;
    this.themeAttributes =
        editor
            .getColorsScheme()
            .getAttributes(DefaultLanguageHighlighterColors.INLAY_TEXT_WITHOUT_BACKGROUND);
    this.editor = editor;
  }

  @Override
  public int calcWidthInPixels(@NotNull Inlay inlay) {
    EditorImpl editor = (EditorImpl) inlay.getEditor();
    return editor.getFontMetrics(Font.PLAIN).stringWidth(text);
  }

  @Override
  public void paint(
      @NotNull Inlay inlay,
      @NotNull Graphics g,
      @NotNull Rectangle targetRegion,
      @NotNull TextAttributes textAttributes) {
    Font font = this.editor.getColorsScheme().getFont(EditorFontType.PLAIN).deriveFont(Font.ITALIC);
    g.setFont(font);
    g.setColor(this.themeAttributes.getForegroundColor());
    g.drawString(this.text, targetRegion.x, targetRegion.y + g.getFontMetrics().getAscent());
  }

  @Override
  public String toString() {
    return "SimpleEditorCustomElementRenderer{" + "text='" + text + '\'' + '}';
  }
}
