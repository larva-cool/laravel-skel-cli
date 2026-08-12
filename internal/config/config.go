// Package config 负责持久化 CLI 配置（API 根地址、认证 token）。
// 配置文件存放于用户主目录下，默认权限 0600 以保护 token。
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config 表示 CLI 的本地配置。
type Config struct {
	BaseURL string `json:"base_url"`
	Token   string `json:"token"`
}

// configPath 返回配置文件路径（~/.config/laravel-skel-cli/config.json）。
func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".config", "laravel-skel-cli")
	return filepath.Join(dir, "config.json"), nil
}

// Load 读取配置；若文件不存在则返回空配置。
func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Save 将配置写入磁盘。
func (c *Config) Save() error {
	path, err := configPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}
