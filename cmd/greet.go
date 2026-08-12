package cmd

import (
	"flag"
	"fmt"
)

// runGreet 实现 greet 子命令：向指定用户打招呼。
func runGreet(args []string) error {
	fs := flag.NewFlagSet("greet", flag.ContinueOnError)

	name := fs.String("name", "世界", "要打招呼的用户名")

	if err := fs.Parse(args); err != nil {
		return err
	}

	fmt.Printf("你好，%s！欢迎使用 laravel-skel-cli。\n", *name)
	return nil
}
