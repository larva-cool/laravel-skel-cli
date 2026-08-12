package cmd

import (
	"flag"
	"fmt"

	"laravel-skel-cli/internal/client"
	"laravel-skel-cli/internal/config"
)

// loginRequest 对应登录接口的请求体。
type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// loginResponse 对应登录接口的响应体，token 字段需按实际后端调整。
type loginResponse struct {
	Token string `json:"token"`
}

// runLogin 实现 login 子命令：调用登录接口获取 token 并保存到配置。
func runLogin(args []string) error {
	fs := flag.NewFlagSet("login", flag.ContinueOnError)
	baseURL := fs.String("base-url", "", "后端 API 地址（可选，首次登录建议指定）")
	email := fs.String("email", "", "登录邮箱")
	password := fs.String("password", "", "登录密码")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if *email == "" || *password == "" {
		return fmt.Errorf("login 需要提供 --email 和 --password")
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
	if err := api.Post("/api/v1/login", &loginRequest{Email: *email, Password: *password}, &resp); err != nil {
		return err
	}
	if resp.Token == "" {
		return fmt.Errorf("登录成功但响应中未包含 token 字段")
	}

	cfg.BaseURL = *baseURL
	cfg.Token = resp.Token
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("保存配置失败: %w", err)
	}

	fmt.Printf("登录成功，token 已保存。base_url: %s\n", *baseURL)
	return nil
}
