// Package cmd 提供命令行入口与子命令分发。
// 子命令遵循 <资源> <动作> 的命名方式，例如 `users list`、`users create`。
package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"
)

// version 通过 -ldflags 注入，用于展示当前程序版本。
var version = "dev"

// Execute 解析并执行命令行参数，返回执行过程中产生的错误。
func Execute() error {
	fs := flag.NewFlagSet("laravel-skel-cli", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}

	args := fs.Args()
	if len(args) == 0 {
		printUsage(os.Stdout)
		return nil
	}

	switch args[0] {
	case "version", "-v", "--version":
		fmt.Printf("laravel-skel-cli %s\n", version)
		return nil
	case "login":
		return runLogin(args[1:])
	case "me":
		return runMe(args[1:])
	case "users":
		return runUsers(args[1:])
	case "help", "-h", "--help":
		printUsage(os.Stdout)
		return nil
	default:
		return fmt.Errorf("未知子命令: %s（可用 help 查看帮助）", args[0])
	}
}

// printUsage 输出帮助信息。
func printUsage(w io.Writer) {
	fmt.Fprintf(w, `laravel-skel-cli - 将后端 API 封装为 CLI 命令

用法:
  laravel-skel-cli <子命令> [参数]

认证:
  login       登录并保存 token
  me          获取当前登录用户信息

资源:
  users list         查询用户列表
  users create       创建用户

其他:
  version     显示版本号
  help        显示帮助信息

提示: 可在配置文件 (~/.config/laravel-skel-cli/config.json)
      中设置 base_url，或通过 login 命令获取。
`)
}
