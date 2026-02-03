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
			name: "完整Front-matter",
			content: `---
title: 测试文章
slug: test-article
summary: 这是摘要
category: 技术
published_at: 2024-01-01
---

# 文章内容
这是正文。`,
			wantTitle: "测试文章",
			wantSlug:  "test-article",
		},
		{
			name: "无Front-matter",
			content: `# 文章标题
这是内容。`,
			wantTitle: "",
			wantSlug:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm, content, err := ParseMarkdown(tt.content)
			if err != nil {
				t.Fatalf("解析失败: %v", err)
			}

			if fm.Title != tt.wantTitle {
				t.Errorf("期望标题 %q，得到 %q", tt.wantTitle, fm.Title)
			}

			if fm.Slug != tt.wantSlug {
				t.Errorf("期望slug %q，得到 %q", tt.wantSlug, fm.Slug)
			}

			if !strings.Contains(content, "文章") {
				t.Error("内容解析不正确")
			}
		})
	}
}

func TestParseMarkdownFile(t *testing.T) {
	tmpDir := t.TempDir()
	mdPath := filepath.Join(tmpDir, "test.md")

	content := `---
title: 文件测试
slug: file-test
---

# 内容
测试内容。`

	os.WriteFile(mdPath, []byte(content), 0644)

	fm, markdown, err := ParseMarkdownFile(mdPath)
	if err != nil {
		t.Fatalf("解析文件失败: %v", err)
	}

	if fm.Title != "文件测试" {
		t.Errorf("期望标题 '文件测试'，得到 %q", fm.Title)
	}

	if !strings.Contains(markdown, "内容") {
		t.Error("Markdown内容解析不正确")
	}
}

func TestMarkdownToHTML(t *testing.T) {
	markdown := `# 标题

这是一段**粗体**文本。`

	html, err := MarkdownToHTML(markdown)
	if err != nil {
		t.Fatalf("转换失败: %v", err)
	}

	if !strings.Contains(html, "<h1>") {
		t.Error("HTML转换不正确，应该包含h1标签")
	}

	if !strings.Contains(html, "<strong>") {
		t.Error("HTML转换不正确，应该包含strong标签")
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
				t.Errorf("期望 %q，得到 %q", tt.want, got)
			}
		})
	}
}
