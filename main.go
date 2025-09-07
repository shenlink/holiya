package main

import (
	"fmt"
	"holiya/file"
	"os"
)

func main() {
	// 执行 go run main.go filename.holiya 或 ./holiya filename.holiya
	err := file.ProcessFile(os.Args[1], os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
