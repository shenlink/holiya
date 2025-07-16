package repl

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"holiya/lexer"
	"holiya/token"
)

// 命令行提示符
const PROMPT = ">> "

// 启动	repl ，循环读取用户输入，并执行，之后输出执行结果
func Start(in io.Reader, out io.Writer) {
	fmt.Println("Hello! This is the holiya programming language!")
	fmt.Printf("Feel free to type in commands\n")
	// 监听 Ctrl+C 信号
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	scanner := bufio.NewScanner(in)
	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()

		// 如果输入 exit，退出
		if line == "exit" {
			return
		}

		l := lexer.New(line)
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}
