package file

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"holiya/evaluator"
	"holiya/lexer"
	"holiya/object"
	"holiya/parser"
)

// ProcessFile 处理指定的文件，逐行读取内容
func ProcessFile(filename string, out io.Writer) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening file %s: %v", filename, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	env := object.NewEnvironment()
	for scanner.Scan() {
		line := scanner.Text()
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

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file %s: %v", filename, err)
	}

	return nil
}
