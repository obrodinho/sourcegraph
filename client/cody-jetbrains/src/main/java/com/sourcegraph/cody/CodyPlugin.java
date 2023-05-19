package com.sourcegraph.cody;

import com.intellij.ide.lightEdit.LightEdit;
import com.intellij.openapi.project.Project;
import org.jetbrains.annotations.Nullable;

public class CodyPlugin {

    public static boolean isSupportedIDE(@Nullable Project project) {
        if (isRemoteIDE()) {
            return true;
        } else {
            return !LightEdit.owns(project);
        }
    }

    public static boolean isRemoteIDE() {
        return "true".equals(System.getProperty("org.jetbrains.projector.server.enable"));
    }
}
