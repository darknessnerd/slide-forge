package config

import (
	"errors"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

var placeholderRe = regexp.MustCompile(`\$\{env:([A-Za-z_][A-Za-z0-9_]*)(?:=([^}]*))?\}`)

// Load reads configuration from CONFIG_PATH env var or ./config/config.yaml.
// If the file does not exist, hardcoded defaults are returned so the server
// works with zero configuration in stdio-only deployments.
func Load() (*AppConfig, error) {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "./config/config.yaml"
	}
	cfg, err := loadFrom(path)
	if err != nil {
		return nil, err
	}
	applyDefaults(cfg)
	return cfg, nil
}

func loadFrom(path string) (*AppConfig, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &AppConfig{}, nil
		}
		return nil, err
	}
	substituted := placeholderRe.ReplaceAllFunc(raw, func(match []byte) []byte {
		groups := placeholderRe.FindSubmatch(match)
		varName := string(groups[1])
		def := string(groups[2])
		val := os.Getenv(varName)
		if val == "" {
			val = def
		}
		return []byte(val)
	})
	var cfg AppConfig
	if err := yaml.Unmarshal(substituted, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func applyDefaults(cfg *AppConfig) {
	if cfg.Transport == "" {
		cfg.Transport = "stdio"
	}
	if cfg.Addr == "" {
		cfg.Addr = ":8080"
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}
}
