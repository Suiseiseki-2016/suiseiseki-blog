package handlers

import (
	"database/sql"
	"net/http"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"blog-suiseiseki/models"
	"blog-suiseiseki/utils"
)

// 匹配 <img ... src="任意路径" ...>
var reImgSrc = regexp.MustCompile(`(?i)<img([^>]*)\s+src="([^"]+)"([^>]*)>`)

type PostsHandler struct {
	db        *sql.DB
	postsPath string
}

func NewPostsHandler(db *sql.DB, postsPath string) *PostsHandler {
	return &PostsHandler{
		db:        db,
		postsPath: postsPath,
	}
}

// GetPosts 获取文章列表
func (h *PostsHandler) GetPosts(c *gin.Context) {
	limit := c.DefaultQuery("limit", "20")
	offset := c.DefaultQuery("offset", "0")

	query := `
		SELECT id, slug, title, summary, category, published_at, content_path, updated_at
		FROM posts
		ORDER BY published_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := h.db.Query(query, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		var publishedAt, updatedAt string
		err := rows.Scan(
			&p.ID,
			&p.Slug,
			&p.Title,
			&p.Summary,
			&p.Category,
			&publishedAt,
			&p.ContentPath,
			&updatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		p.PublishedAt, _ = time.Parse("2006-01-02 15:04:05", publishedAt)
		p.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

		posts = append(posts, p)
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"total": len(posts),
	})
}

// GetPost 获取单篇文章详情（包含 HTML 内容）
func (h *PostsHandler) GetPost(c *gin.Context) {
	slug := c.Param("slug")

	var p models.Post
	var publishedAt, updatedAt string
	err := h.db.QueryRow(`
		SELECT id, slug, title, summary, category, published_at, content_path, updated_at
		FROM posts
		WHERE slug = ?
	`, slug).Scan(
		&p.ID,
		&p.Slug,
		&p.Title,
		&p.Summary,
		&p.Category,
		&publishedAt,
		&p.ContentPath,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	p.PublishedAt, _ = time.Parse("2006-01-02 15:04:05", publishedAt)
	p.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

	_, markdownContent, err := utils.ParseMarkdownFile(p.ContentPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取文章内容失败"})
		return
	}

	htmlContent, err := utils.MarkdownToHTML(markdownContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "转换Markdown失败"})
		return
	}

	// 将相对图片路径重写为 /api/posts-assets/...，使仓库内图片能正确展示
	postDirRel := "."
	if absPosts, err := filepath.Abs(h.postsPath); err == nil {
		if postDirAbs, err := filepath.Abs(filepath.Dir(p.ContentPath)); err == nil {
			if rel, err := filepath.Rel(absPosts, postDirAbs); err == nil {
				postDirRel = filepath.ToSlash(rel)
			}
		}
	}
	if strings.Contains(postDirRel, "..") {
		postDirRel = "."
	}
	htmlContent = rewriteRelativeImgSrc(htmlContent, postDirRel)

	postWithContent := models.PostWithContent{
		Post:    p,
		Content: htmlContent,
	}

	c.JSON(http.StatusOK, postWithContent)
}

// rewriteRelativeImgSrc 将 HTML 中相对路径的 img src 重写为 /api/posts-assets/{postDirRel}/{src}
func rewriteRelativeImgSrc(html, postDirRel string) string {
	postDirRel = path.Clean(postDirRel)
	if strings.Contains(postDirRel, "..") {
		postDirRel = "."
	}
	return reImgSrc.ReplaceAllStringFunc(html, func(match string) string {
		subs := reImgSrc.FindStringSubmatch(match)
		if len(subs) != 4 {
			return match
		}
		prefix, src, suffix := subs[1], subs[2], subs[3]
		src = strings.TrimSpace(src)
		// 只重写相对路径，跳过 http(s) 和空
		if src == "" || strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") || strings.Contains(src, "..") {
			return match
		}
		assetPath := path.Join(postDirRel, src)
		assetPath = path.Clean(assetPath)
		if strings.HasPrefix(assetPath, "..") {
			return match
		}
		return `<img` + prefix + ` src="/api/posts-assets/` + assetPath + `"` + suffix + `>`
	})
}

// ServePostAsset 托管文章仓库内的静态资源（图片等），GET /api/posts-assets/*path
func (h *PostsHandler) ServePostAsset(c *gin.Context) {
	rawPath := strings.TrimPrefix(c.Param("path"), "/")
	if rawPath == "" {
		c.Status(http.StatusNotFound)
		return
	}
	// 禁止路径穿越
	rawPath = filepath.Clean(filepath.FromSlash(rawPath))
	if strings.Contains(rawPath, "..") || filepath.IsAbs(rawPath) {
		c.Status(http.StatusNotFound)
		return
	}
	absPosts, err := filepath.Abs(h.postsPath)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	fullPath := filepath.Join(absPosts, rawPath)
	rel, err := filepath.Rel(absPosts, fullPath)
	if err != nil || strings.Contains(rel, "..") {
		c.Status(http.StatusNotFound)
		return
	}
	c.File(fullPath)
}
