package handlers

import (
	"database/sql"
	"net/http"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"blog-suiseiseki/models"
	"blog-suiseiseki/utils"
)

// Match <img ... src="path" ...>
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

// GetPosts returns the list of posts.
func (h *PostsHandler) GetPosts(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

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
		var publishedAt, updatedAt, category string
		err := rows.Scan(
			&p.ID,
			&p.Slug,
			&p.Title,
			&p.Summary,
			&category,
			&publishedAt,
			&p.ContentPath,
			&updatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		p.Tags = splitTags(category)
		p.PublishedAt, _ = parseTime(publishedAt)
		p.UpdatedAt, _ = parseTime(updatedAt)

		posts = append(posts, p)
	}

	// Get actual total count of all posts (not just this page)
	var total int
	h.db.QueryRow("SELECT COUNT(*) FROM posts").Scan(&total)

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"total": total,
	})
}

// GetPost returns a single post with HTML content.
func (h *PostsHandler) GetPost(c *gin.Context) {
	slug := c.Param("slug")

	var p models.Post
	var publishedAt, updatedAt, category string
	err := h.db.QueryRow(`
		SELECT id, slug, title, summary, category, published_at, content_path, updated_at
		FROM posts
		WHERE slug = ?
	`, slug).Scan(
		&p.ID,
		&p.Slug,
		&p.Title,
		&p.Summary,
		&category,
		&publishedAt,
		&p.ContentPath,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	p.Tags = splitTags(category)
	p.PublishedAt, _ = parseTime(publishedAt)
	p.UpdatedAt, _ = parseTime(updatedAt)

	_, markdownContent, err := utils.ParseMarkdownFile(p.ContentPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read post content"})
		return
	}

	htmlContent, err := utils.MarkdownToHTML(markdownContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "markdown to HTML failed"})
		return
	}

	// Rewrite relative img src to /api/posts-assets/... so repo images display correctly
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

func splitTags(s string) []string {
	tags := []string{}
	for _, t := range strings.Split(s, ",") {
		if trimmed := strings.TrimSpace(t); trimmed != "" {
			tags = append(tags, trimmed)
		}
	}
	return tags
}

func parseTime(s string) (time.Time, error) {
	for _, format := range []string{
		"2006-01-02 15:04:05-07:00",
		"2006-01-02 15:04:05.999999999-07:00",
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02",
	} {
		if t, err := time.Parse(format, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, nil
}

// rewriteRelativeImgSrc rewrites relative img src in HTML to /api/posts-assets/{postDirRel}/{src}.
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
		// Only rewrite relative paths; skip http(s) and empty
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

// ServePostAsset serves static assets (e.g. images) from the posts repo; GET /api/posts-assets/*path.
func (h *PostsHandler) ServePostAsset(c *gin.Context) {
	rawPath := strings.TrimPrefix(c.Param("path"), "/")
	if rawPath == "" {
		c.Status(http.StatusNotFound)
		return
	}
	// Prevent path traversal
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
