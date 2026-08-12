package cmd

import (
	"fmt"
	"strings"

	"laravel-skel-cli/internal/config"
)

// runConfig 实现 config 子命令：查看或修改本地配置。
// 子命令: set-base-url <url>  /  get
func runConfig(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("config 需要指定动作，如: config set-base-url <url> / config get")
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("读取配置失败: %w", err)
	}

	switch args[0] {
	case "set-base-url":
		if len(args) < 2 || strings.TrimSpace(args[1]) == "" {
			return fmt.Errorf("用法: laravel-skel-cli config set-base-url <url>")
		}
		cfg.BaseURL = strings.TrimRight(args[1], "/")
		if err := cfg.Save(); err != nil {
			return fmt.Errorf("保存配置失败: %w", err)
		}
		fmt.Printf("已设置 base_url: %s\n", cfg.BaseURL)
		return nil

	case "set-token":
		if len(args) < 2 || strings.TrimSpace(args[1]) == "" {
			return fmt.Errorf("用法: laravel-skel-cli config set-token <token>")
		}
		cfg.Token = strings.TrimSpace(args[1])
		if err := cfg.Save(); err != nil {
			return fmt.Errorf("保存配置失败: %w", err)
		}
		fmt.Printf("已设置 token: %s\n", maskToken(cfg.Token))
		return nil

	case "get":
		fmt.Printf("base_url: %s\n", cfg.BaseURL)
		fmt.Printf("token:    %s\n", maskToken(cfg.Token))
		return nil

	default:
		return fmt.Errorf("未知动作: config %s（可用: set-base-url / set-token / get）", args[0])
	}
}

// maskToken 对 token 打码，避免明文展示。
func maskToken(t string) string {
	if t == "" {
		return "(未登录)"
	}
	if len(t) <= 8 {
		return "****"
	}
	return t[:4] + "****" + t[len(t)-4:]
}
