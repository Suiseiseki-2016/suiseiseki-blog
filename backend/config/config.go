package config

import (
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Config struct {
	// Server
	Port string
	Mode string // "dev" or "prod"

	// Database
	DBPath string

	// Posts
	PostsPath    string
	PostsRemoteURL string // Remote repo URL to clone when posts dir is empty (e.g. https://github.com/xxx/blog-posts.git)

	// GitHub Webhook
	WebhookSecret string
	GitRepoPath   string // Production path to the posts repo on the server

	// Sync: git pull interval (minutes); 0 = disabled
	SyncIntervalMinutes int

	// Frontend dev server port (for scripts / docs)
	FrontendPort string

	// Environment
	IsDev bool
}

// configFile mirrors config.yaml structure.
type configFile struct {
	Server   struct { Port string `yaml:"port"`; Mode string `yaml:"mode"` }
	Database struct { Path string `yaml:"path"` }
	Posts    struct {
		Path      string `yaml:"path"`
		RemoteURL string `yaml:"remote_url"`
	}
	Webhook  struct {
		Secret      string `yaml:"secret"`
		GitRepoPath string `yaml:"git_repo_path"`
	}
	Sync    struct { IntervalMinutes int `yaml:"interval_minutes"` }
	Frontend struct {
		Port       string `yaml:"port"`
		APIBaseURL string `yaml:"api_base_url"`
	}
}

func Load() *Config {
	cfg := &Config{
		Port:                "8080",
		Mode:                "dev",
		DBPath:              "./blog.db",
		PostsPath:           "../posts",
		WebhookSecret:       "",
		GitRepoPath:         "",
		SyncIntervalMinutes: 0,
	}

	// 1. Load defaults from config.yaml
	configPaths := []string{
		os.Getenv("CONFIG_PATH"),
		"./config.yaml",
		"../config.yaml",
	}
	for _, p := range configPaths {
		if p == "" {
			continue
		}
		abs, _ := filepath.Abs(p)
		data, err := os.ReadFile(abs)
		if err != nil {
			continue
		}
		var f configFile
		if err := yaml.Unmarshal(data, &f); err != nil {
			continue
		}
		if f.Server.Port != "" {
			cfg.Port = f.Server.Port
		}
		if f.Server.Mode != "" {
			cfg.Mode = f.Server.Mode
		}
		if f.Database.Path != "" {
			cfg.DBPath = f.Database.Path
		}
		if f.Posts.Path != "" {
			cfg.PostsPath = f.Posts.Path
		}
		if f.Posts.RemoteURL != "" {
			cfg.PostsRemoteURL = f.Posts.RemoteURL
		}
		if f.Webhook.Secret != "" {
			cfg.WebhookSecret = f.Webhook.Secret
		}
		if f.Webhook.GitRepoPath != "" {
			cfg.GitRepoPath = f.Webhook.GitRepoPath
		}
		if f.Sync.IntervalMinutes > 0 {
			cfg.SyncIntervalMinutes = f.Sync.IntervalMinutes
		}
		if f.Frontend.Port != "" {
			cfg.FrontendPort = f.Frontend.Port
		}
		break
	}

	// 2. Environment variables override
	if v := os.Getenv("PORT"); v != "" {
		cfg.Port = v
	}
	if v := os.Getenv("MODE"); v != "" {
		cfg.Mode = v
	}
	if v := os.Getenv("DB_PATH"); v != "" {
		cfg.DBPath = v
	}
	if v := os.Getenv("POSTS_PATH"); v != "" {
		cfg.PostsPath = v
	}
	if v := os.Getenv("POSTS_REMOTE_URL"); v != "" {
		cfg.PostsRemoteURL = v
	}
	if cfg.Mode != "dev" && (cfg.PostsPath == "" || cfg.PostsPath == "../posts") {
		if v := os.Getenv("GIT_REPO_PATH"); v != "" {
			cfg.PostsPath = v
		} else if cfg.GitRepoPath != "" {
			cfg.PostsPath = cfg.GitRepoPath
		} else {
			cfg.PostsPath = "/var/lib/blog/posts"
		}
	}
	if v := os.Getenv("WEBHOOK_SECRET"); v != "" {
		cfg.WebhookSecret = v
	}
	if v := os.Getenv("GIT_REPO_PATH"); v != "" {
		cfg.GitRepoPath = v
	}
	if v := os.Getenv("SYNC_INTERVAL_MINUTES"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			cfg.SyncIntervalMinutes = n
		}
	}
	if v := os.Getenv("FRONTEND_PORT"); v != "" {
		cfg.FrontendPort = v
	}

	cfg.IsDev = cfg.Mode == "dev"
	if cfg.FrontendPort == "" {
		cfg.FrontendPort = "3000"
	}

	// 3. Resolve posts path to absolute
	absPostsPath, err := filepath.Abs(cfg.PostsPath)
	if err == nil {
		cfg.PostsPath = absPostsPath
	}

	return cfg
}
