package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantTitle string
		wantSlug  string
	}{
		{
			name: "with front-matter",
			content: `---
title: Test Article
slug: test-article
summary: Summary
category: tech
published_at: 2024-01-01
---

# Article content
Body.`,
			wantTitle: "Test Article",
			wantSlug:  "test-article",
		},
		{
			name: "no front-matter",
			content: `# Title
Content.`,
			wantTitle: "",
			wantSlug:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm, content, err := ParseMarkdown(tt.content)
			if err != nil {
				t.Fatalf("parse: %v", err)
			}

			if fm.Title != tt.wantTitle {
				t.Errorf("want title %q, got %q", tt.wantTitle, fm.Title)
			}

			if fm.Slug != tt.wantSlug {
				t.Errorf("want slug %q, got %q", tt.wantSlug, fm.Slug)
			}

			if tt.wantTitle != "" && !strings.Contains(content, "Article") && !strings.Contains(content, "Content") && !strings.Contains(content, "Body") {
				t.Error("content parse incorrect")
			}
		})
	}
}

func TestParseMarkdownFile(t *testing.T) {
	tmpDir := t.TempDir()
	mdPath := filepath.Join(tmpDir, "test.md")

	content := `---
title: File Test
slug: file-test
---

# Content
Test content.`

	os.WriteFile(mdPath, []byte(content), 0644)

	fm, markdown, err := ParseMarkdownFile(mdPath)
	if err != nil {
		t.Fatalf("parse file: %v", err)
	}

	if fm.Title != "File Test" {
		t.Errorf("want title File Test, got %q", fm.Title)
	}

	if !strings.Contains(markdown, "Content") && !strings.Contains(markdown, "content") {
		t.Error("markdown content parse incorrect")
	}
}

func TestMarkdownToHTML(t *testing.T) {
	markdown := `# Title

This is **bold** text.`

	html, err := MarkdownToHTML(markdown)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}

	if !strings.Contains(html, "<h1>") {
		t.Error("HTML should contain h1 tag")
	}

	if !strings.Contains(html, "<strong>") {
		t.Error("HTML should contain strong tag")
	}
}

func TestGenerateSlug(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"/path/to/My Article.md", "my-article"},
		{"/path/to/test-post.md", "test-post"},
		{"/path/to/Hello World.md", "hello-world"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := GenerateSlug(tt.path)
			if got != tt.want {
				t.Errorf("want %q, got %q", tt.want, got)
			}
		})
	}
}
