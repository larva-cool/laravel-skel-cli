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

	fs.Usage = func() {
		printUsage(fs.Output())
	}

	// 子命令定义
	args := os.Args[1:]

	if err := fs.Parse(args); err != nil {
		return err
	}

	// 剩余参数中的第一个作为子命令名
	rest := fs.Args()
	if len(rest) == 0 {
		printUsage(os.Stdout)
		return nil
	}

	switch rest[0] {
	case "version", "-v", "--version":
		fmt.Printf("laravel-skel-cli %s\n", version)
		return nil
	case "greet":
		return runGreet(rest[1:])
	case "help", "-h", "--help":
		printUsage(os.Stdout)
		return nil
	default:
		return fmt.Errorf("未知子命令: %s", rest[0])
	}
}

// printUsage 输出帮助信息。
func printUsage(w io.Writer) {
	fmt.Fprintf(w, `laravel-skel-cli - 一个 Go 控制台项目骨架

用法:
  laravel-skel-cli <子命令> [参数]

子命令:
  version   显示版本号
  greet     向用户打招呼
  help      显示帮助信息
`)
}
