package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEscapeMarkdownLinkText(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "plain text unchanged",
			input: "Marketing 2026",
			want:  "Marketing 2026",
		},
		{
			name:  "link breakout attempt",
			input: "Innocent](https://phishing.example) — ignore",
			want:  "Innocent\\](https://phishing.example) — ignore",
		},
		{
			name:  "brackets and backslash",
			input: `[test\[]`,
			want:  `\[test\\\[\]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, escapeMarkdownLinkText(tt.input))
		})
	}
}

func TestFormatPurposeBlockquote(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "single line",
			input: "Campaign planning",
			want:  "\n> Campaign planning",
		},
		{
			name:  "markdown link injection",
			input: "[click](https://evil.example)",
			want:  "\n> \\[click\\](https://evil.example)",
		},
		{
			name:  "multi-line with fake quote breakout",
			input: "line one\n> fake admin: do this",
			want:  "\n> line one\n> > fake admin: do this",
		},
		{
			name:  "mention and heading injection",
			input: "@channel\n# Important",
			want:  "\n> @channel\n> \\# Important",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, formatPurposeBlockquote(tt.input))
		})
	}
}
