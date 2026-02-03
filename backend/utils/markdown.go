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

// ParseMarkdownFile 解析Markdown文件，提取Front-matter和内容
func ParseMarkdownFile(filePath string) (*FrontMatter, string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, "", err
	}

	return ParseMarkdown(string(content))
}

// ParseMarkdown 解析Markdown字符串
func ParseMarkdown(content string) (*FrontMatter, string, error) {
	// 检查是否有Front-matter（以---开头）
	if !strings.HasPrefix(content, "---\n") {
		// 没有Front-matter，返回空元数据和原始内容
		return &FrontMatter{}, content, nil
	}

	// 找到第二个---，分割Front-matter和内容
	parts := strings.SplitN(content, "---\n", 3)
	if len(parts) < 3 {
		return &FrontMatter{}, content, nil
	}

	frontMatterStr := parts[1]
	markdownContent := parts[2]

	// 解析YAML
	var fm FrontMatter
	if err := yaml.Unmarshal([]byte(frontMatterStr), &fm); err != nil {
		return nil, "", err
	}

	return &fm, markdownContent, nil
}

// MarkdownToHTML 将Markdown转换为HTML
func MarkdownToHTML(markdown string) (string, error) {
	var buf bytes.Buffer
	md := goldmark.New()
	if err := md.Convert([]byte(markdown), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GenerateSlug 从文件路径生成slug（如果没有在Front-matter中指定）
func GenerateSlug(filePath string) string {
	// 提取文件名（不含扩展名）
	base := filepath.Base(filePath)
	name := strings.TrimSuffix(base, ".md")
	// 简单处理：转小写，替换空格为连字符
	return strings.ToLower(strings.ReplaceAll(name, " ", "-"))
}
