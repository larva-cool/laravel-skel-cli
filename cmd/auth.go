package cmd

import (
	"flag"
	"fmt"

	"laravel-skel-cli/internal/client"
	"laravel-skel-cli/internal/config"
)

// loginResponse 对应 /admin/auth/login 的响应，token 字段为 Sanctum Bearer 令牌。
type loginResponse struct {
	AccessToken string `json:"access_token"`
}

// runLogin 实现 login 子命令：调用登录接口获取 token 并保存到配置。
func runLogin(args []string) error {
	fs := flag.NewFlagSet("login", flag.ContinueOnError)
	account := fs.String("account", "", "登录账号")
	password := fs.String("password", "", "登录密码")
	baseURL := fs.String("base-url", "", "后端 API 地址（覆盖配置值）")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *account == "" || *password == "" {
		return fmt.Errorf("login 需要提供 --account 和 --password")
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("读取配置失败: %w", err)
	}
	if *baseURL == "" {
		*baseURL = cfg.BaseURL
	}
	if *baseURL == "" {
		return fmt.Errorf("请通过 --base-url 指定后端 API 地址")
	}

	api := client.New(*baseURL, "")
	var resp loginResponse
	if err := api.Post("/admin/auth/login",
		map[string]any{"account": *account, "password": *password}, &resp); err != nil {
		return err
	}
	if resp.AccessToken == "" {
		return fmt.Errorf("登录成功但响应中未包含 access_token 字段")
	}

	cfg.BaseURL = *baseURL
	cfg.Token = resp.AccessToken
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("保存配置失败: %w", err)
	}
	fmt.Printf("登录成功，token 已保存。base_url: %s\n", *baseURL)
	return nil
}

// runLogout 实现 logout 子命令：调用退出接口并清除本地 token。
func runLogout(args []string) error {
	fs := flag.NewFlagSet("logout", flag.ContinueOnError)
	baseURL := fs.String("base-url", "", "后端 API 地址（覆盖配置值）")
	if err := fs.Parse(args); err != nil {
		return err
	}

	api, err := loadClient(*baseURL)
	if err != nil {
		return err
	}
	// 尽力通知服务端退出，忽略响应内容
	_ = api.Post("/admin/auth/logout", map[string]any{}, nil)

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("读取配置失败: %w", err)
	}
	cfg.Token = ""
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("保存配置失败: %w", err)
	}
	fmt.Println("已退出登录，本地 token 已清除。")
	return nil
}
