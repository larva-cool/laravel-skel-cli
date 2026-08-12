// Package cmd 提供命令行入口与子命令分发。
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
	case "logout":
		return runLogout(args[1:])
	case "config":
		return runConfig(args[1:])
	case "whoami":
		return runCall(append([]string{"auth.info"}, args[1:]...))
	case "call":
		return runCall(args[1:])
	case "list":
		return runList(args[1:])
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
  login --account <账号> --password <密码> [--base-url <地址>]   登录并保存 token
  logout                                                       退出并清除 token
  whoami                                                       获取当前登录管理员信息
  config set-base-url <地址>   设置后端 API 地址到配置
  config set-token <token>     设置 token 到配置
  config get                  查看当前配置

接口调用（数据驱动，见 list）:
  call <slug> [--参数 值]   调用指定接口并输出 JSON
  list                      列出全部可调用接口

其他:
  version     显示版本号
  help        显示帮助信息
`)
}
