// Copyright (c) 2026-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package main

import "strings"

func escapeMarkdownLinkText(s string) string {
	replacer := strings.NewReplacer(
		"\\", "\\\\",
		"[", "\\[",
		"]", "\\]",
		"@", "\\@",
	)
	return replacer.Replace(s)
}

func escapeMarkdownPlainText(s string) string {
	replacer := strings.NewReplacer(
		"\\", "\\\\",
		"[", "\\[",
		"]", "\\]",
		"*", "\\*",
		"_", "\\_",
		"`", "\\`",
		"~", "\\~",
		"#", "\\#",
		"|", "\\|",
		"@", "\\@",
	)
	return replacer.Replace(s)
}

func formatPurposeBlockquote(purpose string) string {
	lines := strings.Split(purpose, "\n")
	for i, line := range lines {
		lines[i] = "> " + escapeMarkdownPlainText(line)
	}
	return "\n" + strings.Join(lines, "\n")
}
