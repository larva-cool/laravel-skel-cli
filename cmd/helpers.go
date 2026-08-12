package cmd

import (
	"encoding/json"
	"fmt"

	"laravel-skel-cli/internal/client"
	"laravel-skel-cli/internal/config"
)

// loadClient 读取本地配置并构造 API 客户端。
// baseURL 优先使用命令行传入值，其次使用配置中的值。
func loadClient(flagBaseURL string) (*client.APIClient, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("读取配置失败: %w", err)
	}
	baseURL := flagBaseURL
	if baseURL == "" {
		baseURL = cfg.BaseURL
	}
	if baseURL == "" {
		return nil, fmt.Errorf("未设置 API 地址，请先使用 login 命令或配置 base_url")
	}
	return client.New(baseURL, cfg.Token), nil
}

// printJSON 将任意值以缩进 JSON 形式输出到标准输出。
func printJSON(v any) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println(v)
		return
	}
	fmt.Println(string(data))
}
