package models

import "time"

type Post struct {
	ID          int       `json:"id"`
	Slug        string    `json:"slug"`
	Title       string    `json:"title"`
	Summary     string    `json:"summary"`
	Tags        []string  `json:"tags"`
	PublishedAt time.Time `json:"published_at"`
	ContentPath string    `json:"-"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PostWithContent struct {
	Post
	Content string `json:"content"`
}
