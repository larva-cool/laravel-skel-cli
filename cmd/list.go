package cmd

import (
	"fmt"

	"laravel-skel-cli/internal/apidefs"
)

// runList 实现 list 子命令：列出所有可调用的接口。
func runList(_ []string) error {
	fmt.Printf("可调用接口（共 %d 个）:\n", len(apidefs.All()))
	for _, ep := range apidefs.All() {
		fmt.Printf("  %-32s %-6s %s\n", ep.Slug, ep.Method, ep.Path)
	}
	fmt.Println("\n用法: laravel-skel-cli call <slug> [--参数 值]")
	return nil
}
