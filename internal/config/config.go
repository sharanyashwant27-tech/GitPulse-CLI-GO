package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds application-wide settings.
type Config struct {
	RepoPath    string        `mapstructure:"repo"`
	Theme       string        `mapstructure:"theme"`
	LogLevel    string        `mapstructure:"log_level"`
	WatchInterval time.Duration `mapstructure:"watch_interval"`
	ExportDir   string        `mapstructure:"export_dir"`
	Since       string        `mapstructure:"since"`
	Until       string        `mapstructure:"until"`
	Limit       int           `mapstructure:"limit"`
	NoColor     bool          `mapstructure:"no_color"`
	Verbose     bool          `mapstructure:"verbose"`
}

// Default returns a Config with sensible defaults.
func Default() *Config {
	cwd, _ := os.Getwd()
	return &Config{
		RepoPath:      cwd,
		Theme:         "default",
		LogLevel:      "info",
		WatchInterval: 5 * time.Second,
		ExportDir:     "reports",
		Limit:         20,
		NoColor:       false,
		Verbose:       false,
	}
}

// Load reads configuration from flags, env, and optional config file.
func Load(cfgFile string) (*Config, error) {
	cfg := Default()

	v := viper.New()
	v.SetEnvPrefix("GITPULSE")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	v.SetDefault("repo", cfg.RepoPath)
	v.SetDefault("theme", cfg.Theme)
	v.SetDefault("log_level", cfg.LogLevel)
	v.SetDefault("watch_interval", cfg.WatchInterval)
	v.SetDefault("export_dir", cfg.ExportDir)
	v.SetDefault("limit", cfg.Limit)

	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err == nil {
			v.AddConfigPath(filepath.Join(home, ".config", "gitpulse"))
		}
		v.AddConfigPath(".")
		v.SetConfigName("gitpulse")
		v.SetConfigType("yaml")
	}

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok && cfgFile != "" {
			return nil, fmt.Errorf("read config: %w", err)
		}
	}

	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	if cfg.RepoPath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		cfg.RepoPath = cwd
	}

	abs, err := filepath.Abs(cfg.RepoPath)
	if err != nil {
		return nil, err
	}
	cfg.RepoPath = abs

	cfg.Theme = strings.ToLower(strings.TrimSpace(cfg.Theme))
	cfg.Theme = strings.ReplaceAll(cfg.Theme, "_", "-")
	cfg.Theme = strings.Join(strings.Fields(cfg.Theme), "-")
	if cfg.Theme == "tokyo" || cfg.Theme == "tokyonight" {
		cfg.Theme = "tokyo-night"
	}
	if cfg.Theme == "" {
		cfg.Theme = "default"
	}
	return cfg, nil
}

// ValidThemes lists supported theme identifiers.
func ValidThemes() []string {
	return []string{
		"default",
		"dracula",
		"nord",
		"catppuccin",
		"tokyo-night",
		"gruvbox",
		"solarized",
	}
}

// ValidateTheme returns an error if the theme name is unknown.
func ValidateTheme(name string) error {
	name = strings.ToLower(strings.TrimSpace(name))
	name = strings.ReplaceAll(name, "_", "-")
	name = strings.Join(strings.Fields(name), "-")
	if name == "tokyo" || name == "tokyonight" {
		name = "tokyo-night"
	}
	for _, t := range ValidThemes() {
		if t == name {
			return nil
		}
	}
	return fmt.Errorf("unknown theme %q (supported: %s)", name, strings.Join(ValidThemes(), ", "))
}
