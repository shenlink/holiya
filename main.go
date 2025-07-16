package main

import (
	"holiya/repl"
	"os"
)

func main() {
	// 直接运行，启动 repl 模式
	repl.Start(os.Stdin, os.Stdout)
}
