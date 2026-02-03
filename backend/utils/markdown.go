package utils

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	"gopkg.in/yaml.v3"
)

type FrontMatter struct {
	Title       string `yaml:"title"`
	Summary     string `yaml:"summary"`
	Category    string `yaml:"category"`
	PublishedAt string `yaml:"published_at"`
	Slug        string `yaml:"slug"`
}

// ParseMarkdownFile parses a Markdown file and extracts front-matter and body.
func ParseMarkdownFile(filePath string) (*FrontMatter, string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, "", err
	}

	return ParseMarkdown(string(content))
}

// ParseMarkdown parses a Markdown string and extracts front-matter and body.
func ParseMarkdown(content string) (*FrontMatter, string, error) {
	if !strings.HasPrefix(content, "---\n") {
		return &FrontMatter{}, content, nil
	}

	parts := strings.SplitN(content, "---\n", 3)
	if len(parts) < 3 {
		return &FrontMatter{}, content, nil
	}

	frontMatterStr := parts[1]
	markdownContent := parts[2]

	var fm FrontMatter
	if err := yaml.Unmarshal([]byte(frontMatterStr), &fm); err != nil {
		return nil, "", err
	}

	return &fm, markdownContent, nil
}

// MarkdownToHTML converts Markdown to HTML.
func MarkdownToHTML(markdown string) (string, error) {
	var buf bytes.Buffer
	md := goldmark.New()
	if err := md.Convert([]byte(markdown), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GenerateSlug generates a slug from the file path when not set in front-matter.
func GenerateSlug(filePath string) string {
	base := filepath.Base(filePath)
	name := strings.TrimSuffix(base, ".md")
	return strings.ToLower(strings.ReplaceAll(name, " ", "-"))
}
