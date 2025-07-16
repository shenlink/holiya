package main

import (
	"fmt"
	"holiya/file"
	"holiya/repl"
	"os"
)

func main() {
	// 直接运行，启动 repl 模式
	if len(os.Args) < 2 {
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	// 执行 go run main.go filename.holiya 或 ./holira filename.holiya
	err := file.ProcessFile(os.Args[1], os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
