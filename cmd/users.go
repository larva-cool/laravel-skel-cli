package cmd

import (
	"flag"
	"fmt"
)

// runUsers 实现 users 资源命令，分发到 list / create 等动作。
func runUsers(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("users 需要指定动作，如: users list / users create")
	}

	switch args[0] {
	case "list":
		return runUsersList(args[1:])
	case "create":
		return runUsersCreate(args[1:])
	default:
		return fmt.Errorf("未知动作: users %s（可用: list / create）", args[0])
	}
}

// runUsersList 实现 users list：查询用户列表。
func runUsersList(args []string) error {
	fs := flag.NewFlagSet("users list", flag.ContinueOnError)
	baseURL := fs.String("base-url", "", "后端 API 地址（覆盖配置值）")
	page := fs.Int("page", 1, "页码")
	perPage := fs.Int("per-page", 20, "每页数量")
	if err := fs.Parse(args); err != nil {
		return err
	}

	api, err := loadClient(*baseURL)
	if err != nil {
		return err
	}

	var result map[string]any
	if err := api.Get(fmt.Sprintf("/api/v1/users?page=%d&per_page=%d", *page, *perPage), &result); err != nil {
		return err
	}
	printJSON(result)
	return nil
}

// runUsersCreate 实现 users create：创建用户。
func runUsersCreate(args []string) error {
	fs := flag.NewFlagSet("users create", flag.ContinueOnError)
	baseURL := fs.String("base-url", "", "后端 API 地址（覆盖配置值）")
	name := fs.String("name", "", "用户名")
	email := fs.String("email", "", "邮箱")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *name == "" || *email == "" {
		return fmt.Errorf("users create 需要提供 --name 和 --email")
	}

	api, err := loadClient(*baseURL)
	if err != nil {
		return err
	}

	body := map[string]any{"name": *name, "email": *email}
	var result map[string]any
	if err := api.Post("/api/v1/users", body, &result); err != nil {
		return err
	}
	printJSON(result)
	return nil
}
