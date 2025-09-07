package main

import (
	"fmt"
	"holiya/file"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		// 没有参数时，输出帮助信息
		printHelp()
		return
	}

	// 执行 go run main.go filename.holiya 或 ./holiya filename.holiya
	err := file.ProcessFile(os.Args[1], os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

// printHelp 输出帮助信息
func printHelp() {
	fmt.Println("Usage:")
	fmt.Println("  ./holiya filename.holiya     Process the specified file")
	fmt.Println("  go run main.go filename.holiya     Process the specified file")
}
