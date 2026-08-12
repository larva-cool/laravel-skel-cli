package cmd

import "flag"

// runMe 实现 me 子命令：获取当前登录用户信息。
func runMe(args []string) error {
	fs := flag.NewFlagSet("me", flag.ContinueOnError)
	baseURL := fs.String("base-url", "", "后端 API 地址（覆盖配置值）")
	if err := fs.Parse(args); err != nil {
		return err
	}

	api, err := loadClient(*baseURL)
	if err != nil {
		return err
	}

	var result map[string]any
	if err := api.Get("/api/v1/me", &result); err != nil {
		return err
	}
	printJSON(result)
	return nil
}
