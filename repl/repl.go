package repl

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"holiya/evaluator"
	"holiya/lexer"
	"holiya/object"
	"holiya/parser"
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
	env := object.NewEnvironment()
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
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			io.WriteString(out, "parser errors:\n")
			for _, msg := range p.Errors() {
				io.WriteString(out, "\t"+msg+"\n")
			}
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}
