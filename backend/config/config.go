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
	PostsPath string

	// GitHub Webhook
	WebhookSecret string
	GitRepoPath   string // 生产环境的Git仓库路径

	// 定期同步：从远程仓库 git pull 并更新数据库；0 表示不启用
	SyncIntervalMinutes int

	// Environment
	IsDev bool
}

// configFile 对应 config.yaml 结构
type configFile struct {
	Server   struct { Port string `yaml:"port"`; Mode string `yaml:"mode"` }
	Database struct { Path string `yaml:"path"` }
	Posts    struct { Path string `yaml:"path"` }
	Webhook  struct {
		Secret      string `yaml:"secret"`
		GitRepoPath string `yaml:"git_repo_path"`
	}
	Sync struct {
		IntervalMinutes int `yaml:"interval_minutes"`
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

	// 1. 尝试从 config.yaml 加载默认值
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
		if f.Webhook.Secret != "" {
			cfg.WebhookSecret = f.Webhook.Secret
		}
		if f.Webhook.GitRepoPath != "" {
			cfg.GitRepoPath = f.Webhook.GitRepoPath
		}
		if f.Sync.IntervalMinutes > 0 {
			cfg.SyncIntervalMinutes = f.Sync.IntervalMinutes
		}
		break
	}

	// 2. 环境变量覆盖
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

	cfg.IsDev = cfg.Mode == "dev"

	// 3. 文章路径转为绝对路径
	absPostsPath, err := filepath.Abs(cfg.PostsPath)
	if err == nil {
		cfg.PostsPath = absPostsPath
	}

	return cfg
}
