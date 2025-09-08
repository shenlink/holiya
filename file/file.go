package file

import (
	"io"
	"os"

	"holiya/evaluator"
	"holiya/lexer"
	"holiya/object"
	"holiya/parser"
)

// ProcessFile 处理指定的文件，逐行读取内容
func ProcessFile(filename string, out io.Writer) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	content := string(data)
	env := object.NewEnvironment()
	l := lexer.New(content)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		io.WriteString(out, "parser errors:\n")
		for _, msg := range p.Errors() {
			io.WriteString(out, "\t"+msg+"\n")
		}
	}

	for _, statement := range program.Statements {
		evaluated := evaluator.Eval(statement, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}

	return nil
}
