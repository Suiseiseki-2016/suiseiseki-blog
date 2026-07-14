package utils

import (
	"bytes"
	"html"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"gopkg.in/yaml.v3"
)

type FrontMatter struct {
	Title       string   `yaml:"title"`
	Summary     string   `yaml:"summary"`
	Category    string   `yaml:"category"`
	Tags        []string `yaml:"tags"`
	PublishedAt string   `yaml:"published_at"`
	Slug        string   `yaml:"slug"`
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
	// Normalize CRLF to LF so front-matter detection works on Windows-authored files
	content = strings.ReplaceAll(content, "\r\n", "\n")

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

var reMermaid = regexp.MustCompile(`(?s)<pre><code class="language-mermaid">(.*?)</code></pre>`)

// MarkdownToHTML converts Markdown to HTML, with mermaid fenced blocks rendered
// as <div class="mermaid"> elements for client-side rendering.
func MarkdownToHTML(markdown string) (string, error) {
	var buf bytes.Buffer
	md := goldmark.New()
	if err := md.Convert([]byte(markdown), &buf); err != nil {
		return "", err
	}
	return convertMermaidBlocks(buf.String()), nil
}

// convertMermaidBlocks replaces goldmark-rendered mermaid code blocks with
// <div class="mermaid"> elements, HTML-unescaping the diagram source.
func convertMermaidBlocks(h string) string {
	return reMermaid.ReplaceAllStringFunc(h, func(match string) string {
		subs := reMermaid.FindStringSubmatch(match)
		if len(subs) != 2 {
			return match
		}
		return `<div class="mermaid">` + html.UnescapeString(subs[1]) + `</div>`
	})
}

// GenerateSlug generates a slug from the file path when not set in front-matter.
func GenerateSlug(filePath string) string {
	base := filepath.Base(filePath)
	name := strings.TrimSuffix(base, ".md")
	return strings.ToLower(strings.ReplaceAll(name, " ", "-"))
}
