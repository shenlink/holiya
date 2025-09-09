package ast

import (
	"holiya/token"
	"testing"
)

// 表达式测试结构体
type expressions struct {
	expression      Expression
	expectedLiteral string
	expectedString  string
}

// 语句测试结构体
type statements struct {
	statement       Statement
	expectedLiteral string
	expectedString  string
}

// 程序测试结构体
type programs struct {
	program        Program
	expectedString string
}

// 测试 Identifier
func TestIdentifier(t *testing.T) {
	identifiers := []expressions{
		{
			expression: &Identifier{
				Token: token.Token{Type: token.IDENTIFIER, Literal: "test"},
				Value: "test",
			},
			expectedLiteral: "test",
			expectedString:  "test",
		},
	}
	if !testExpression(t, identifiers) {
		return
	}
}

// 测试 IntegerLiteral
func TestIntegerLiteral(t *testing.T) {
	integers := []expressions{
		{
			expression: &IntegerLiteral{
				Token: token.Token{Type: token.INT, Literal: "5"},
				Value: 5,
			},
			expectedLiteral: "5",
			expectedString:  "5",
		},
	}

	if !testExpression(t, integers) {
		return
	}
}

// 测试 FloatLiteral
func TestFloatLiteral(t *testing.T) {
	floats := []expressions{
		{
			expression: &FloatLiteral{
				Token: token.Token{Type: token.INT, Literal: "5.0"},
				Value: 5.0,
			},
			expectedLiteral: "5.0",
			expectedString:  "5.0",
		},
	}

	if !testExpression(t, floats) {
		return
	}
}

// 测试 StringLiteral
func TestStringLiteral(t *testing.T) {
	stringLiterals := []expressions{
		{
			expression: &StringLiteral{
				Token: token.Token{Type: token.STRING, Literal: "hello world"},
				Value: "hello world",
			},
			expectedLiteral: "hello world",
			expectedString:  "hello world",
		},
	}

	if !testExpression(t, stringLiterals) {
		return
	}
}

// 测试 InfixExpression
func TestInfixExpression(t *testing.T) {
	infixes := []expressions{
		{
			expression: &InfixExpression{
				Token: token.Token{Type: token.INT, Literal: "+"},
				Left: &IntegerLiteral{
					Token: token.Token{Type: token.INT, Literal: "5"},
					Value: 5,
				},
				Operator: "+",
				Right: &IntegerLiteral{
					Token: token.Token{Type: token.INT, Literal: "5"},
					Value: 5,
				},
			},
			expectedLiteral: "+",
			expectedString:  "(5 + 5)",
		},
		{
			expression: &InfixExpression{
				Token: token.Token{Type: token.INT, Literal: "+"},
				Left: &IntegerLiteral{
					Token: token.Token{Type: token.INT, Literal: "5"},
					Value: 5,
				},
				Operator: "+",
				Right: &InfixExpression{
					Token: token.Token{Type: token.INT, Literal: "*"},
					Left: &IntegerLiteral{
						Token: token.Token{Type: token.INT, Literal: "5"},
						Value: 5,
					},
					Operator: "*",
					Right: &IntegerLiteral{
						Token: token.Token{Type: token.INT, Literal: "5"},
						Value: 5,
					},
				},
			},
			expectedLiteral: "+",
			expectedString:  "(5 + (5 * 5))",
		},
	}

	if !testExpression(t, infixes) {
		return
	}
}

// 测试 PrefixExpression
func TestPrefixExpression(t *testing.T) {
	prefixes := []expressions{
		{
			expression: &PrefixExpression{
				Token:    token.Token{Type: token.MINUS, Literal: "-"},
				Operator: "-",
				Right: &IntegerLiteral{
					Token: token.Token{Type: token.INT, Literal: "5"},
					Value: 5,
				},
			},
			expectedLiteral: "-",
			expectedString:  "(-5)",
		},
	}

	if !testExpression(t, prefixes) {
		return
	}
}

// 测试 FunctionStatement
func TestFunctionStatement(t *testing.T) {
	funcs := []statements{
		{
			statement: &FunctionStatement{
				Token:      token.Token{Type: token.FUNCTION, Literal: "fn"},
				Parameters: []*Identifier{},
				Body:       getBlockStatement(),
			},
			expectedLiteral: "fn",
			expectedString:  "fn()let myVar = myVar;let myVar = 5;return 5;",
		},
		{
			statement: &FunctionStatement{
				Token: token.Token{Type: token.FUNCTION, Literal: "fn"},
				Parameters: []*Identifier{
					{
						Token: token.Token{Type: token.IDENTIFIER, Literal: "myVar"},
						Value: "myVar",
					},
				},
				Body: getBlockStatement(),
			},
			expectedLiteral: "fn",
			expectedString:  "fn(myVar)let myVar = myVar;let myVar = 5;return 5;",
		},
	}

	if !testStatement(t, funcs) {
		return
	}
}

// 测试 Boolean
func TestBoolean(t *testing.T) {
	booleans := []expressions{
		{
			expression: &Boolean{
				Token: token.Token{Type: token.TRUE, Literal: "true"},
				Value: true,
			},
			expectedLiteral: "true",
			expectedString:  "true",
		},
		{
			expression: &Boolean{
				Token: token.Token{Type: token.FALSE, Literal: "false"},
				Value: false,
			},
			expectedLiteral: "false",
			expectedString:  "false",
		},
	}

	if !testExpression(t, booleans) {
		return
	}
}

// 测试 IfExpression
func TestIfExpression(t *testing.T) {
	ifs := []expressions{
		{
			expression: &IfExpression{
				Token: token.Token{Type: token.IF, Literal: "if"},
				Condition: &InfixExpression{
					Token: token.Token{Type: token.IF, Literal: "if"},
					Left: &InfixExpression{
						Token: token.Token{Type: token.INT, Literal: "5"},
						Left: &IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
						Operator: "+",
						Right: &IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
					},
					Operator: ">",
					Right: &IntegerLiteral{
						Token: token.Token{Type: token.INT, Literal: "10"},
						Value: 10,
					},
				},
				Consequence: getBlockStatement(),
				Alternative: getBlockStatement(),
			},
			expectedLiteral: "if",
			expectedString:  "if((5 + 5) > 10) let myVar = myVar;let myVar = 5;return 5;else let myVar = myVar;let myVar = 5;return 5;",
		},
	}

	if !testExpression(t, ifs) {
		return
	}
}

// 测试 CallExpression
func TestCallExpression(t *testing.T) {
	calls := []expressions{
		{
			expression: &CallExpression{
				Token: token.Token{Type: token.LPAREN, Literal: "("},
				Function: &Identifier{
					Token: token.Token{Type: token.IDENTIFIER, Literal: "add"},
					Value: "add",
				},
				Arguments: []Expression{
					&Identifier{
						Token: token.Token{Type: token.IDENTIFIER, Literal: "myVar"},
						Value: "myVar",
					},
					&Identifier{
						Token: token.Token{Type: token.IDENTIFIER, Literal: "myVar2"},
						Value: "myVar2",
					},
				},
			},
			expectedLiteral: "(",
			expectedString:  "add(myVar, myVar2)",
		},
	}

	if !testExpression(t, calls) {
		return
	}
}

// 测试 HashLiteral
func TestHashLiteral(t *testing.T) {
	maps := []expressions{
		{
			expression: &HashLiteral{
				Token: token.Token{Type: token.STRING, Literal: ""},
				Pairs: map[Expression]Expression{
					&StringLiteral{
						Token: token.Token{
							Type:    token.STRING,
							Literal: "a",
						},
						Value: "a",
					}: &StringLiteral{
						Token: token.Token{
							Type:    token.STRING,
							Literal: "a",
						},
						Value: "a",
					},
					&StringLiteral{
						Token: token.Token{
							Type:    token.STRING,
							Literal: "b",
						},
						Value: "b",
					}: &StringLiteral{
						Token: token.Token{
							Type:    token.STRING,
							Literal: "b",
						},
						Value: "b",
					},
					&StringLiteral{
						Token: token.Token{
							Type:    token.STRING,
							Literal: "c",
						},
						Value: "c",
					}: &IntegerLiteral{
						Token: token.Token{
							Type:    token.INT,
							Literal: "1",
						},
						Value: 1,
					},
				},
			},
			expectedString: "{a:a, b:b, c:1}",
		},
	}

	if !testExpression(t, maps) {
		return
	}
}

// 测试 ArrayLiteral
func TestArrayLiteral(t *testing.T) {
	maps := []expressions{
		{
			expression: &ArrayLiteral{
				Token: token.Token{
					Type:    token.IDENTIFIER,
					Literal: "lists",
				},
				Elements: []Expression{
					&IntegerLiteral{
						Token: token.Token{
							Type:    token.INT,
							Literal: "1",
						},
						Value: 1,
					},
					&IntegerLiteral{
						Token: token.Token{
							Type:    token.INT,
							Literal: "2",
						},
						Value: 2,
					},
					&IntegerLiteral{
						Token: token.Token{
							Type:    token.INT,
							Literal: "3",
						},
						Value: 3,
					},
				},
			},
			expectedLiteral: "lists",
			expectedString:  "[1, 2, 3]",
		},
	}

	if !testExpression(t, maps) {
		return
	}
}

// 测试 IndexExpression
func TestIndexExpression(t *testing.T) {
	maps := []expressions{
		{
			expression: &IndexExpression{
				Token: token.Token{
					Type:    token.IDENTIFIER,
					Literal: "[",
				},
				Left: &ArrayLiteral{
					Token: token.Token{
						Type:    token.IDENTIFIER,
						Literal: "[",
					},
					Elements: []Expression{
						&IntegerLiteral{
							Token: token.Token{
								Type:    token.INT,
								Literal: "1",
							},
							Value: 1,
						},
						&IntegerLiteral{
							Token: token.Token{
								Type:    token.INT,
								Literal: "2",
							},
							Value: 2,
						},
						&IntegerLiteral{
							Token: token.Token{
								Type:    token.INT,
								Literal: "3",
							},
							Value: 3,
						},
					},
				},
				Index: &IntegerLiteral{
					Token: token.Token{
						Type:    token.INT,
						Literal: "0",
					},
					Value: 0,
				},
			},
			expectedLiteral: "[",
			expectedString:  "([1, 2, 3][0])",
		},
		{
			expression: &IndexExpression{
				Token: token.Token{
					Type:    token.STRING,
					Literal: "[",
				},
				Left: &HashLiteral{
					Token: token.Token{Type: token.STRING, Literal: "{"},
					Pairs: map[Expression]Expression{
						&StringLiteral{
							Token: token.Token{
								Type:    token.STRING,
								Literal: "a",
							},
							Value: "a",
						}: &StringLiteral{
							Token: token.Token{
								Type:    token.STRING,
								Literal: "a",
							},
							Value: "a",
						},
						&StringLiteral{
							Token: token.Token{
								Type:    token.STRING,
								Literal: "b",
							},
							Value: "b",
						}: &StringLiteral{
							Token: token.Token{
								Type:    token.STRING,
								Literal: "b",
							},
							Value: "b",
						},
						&StringLiteral{
							Token: token.Token{
								Type:    token.STRING,
								Literal: "c",
							},
							Value: "c",
						}: &IntegerLiteral{
							Token: token.Token{
								Type:    token.INT,
								Literal: "1",
							},
							Value: 1,
						},
					},
				},
				Index: &StringLiteral{
					Token: token.Token{
						Type:    token.STRING,
						Literal: "a",
					},
					Value: "a",
				},
			},
			expectedLiteral: "[",
			expectedString:  "({a:a, b:b, c:1}[a])",
		},
		{
			expression: &IndexExpression{
				Token: token.Token{
					Type:    token.STRING,
					Literal: "[",
				},
				Left: &StringLiteral{
					Token: token.Token{
						Type:    token.STRING,
						Literal: "hello world",
					},
					Value: "hello world",
				},
				Index: &IntegerLiteral{
					Token: token.Token{
						Type:    token.INT,
						Literal: "1",
					},
					Value: 1,
				},
			},
			expectedLiteral: "[",
			expectedString:  "(hello world[1])",
		},
	}

	if !testExpression(t, maps) {
		return
	}
}

// 测试 LetStatement
func TestLetStatement(t *testing.T) {
	lets := []statements{
		{
			statement:       getLetStatement("myVar", "myVar"),
			expectedLiteral: "let",
			expectedString:  "let myVar = myVar;",
		},
		{
			statement:       getLetStatement("myVar", "5"),
			expectedLiteral: "let",
			expectedString:  "let myVar = 5;",
		},
	}

	if !testStatement(t, lets) {
		return
	}
}

// 测试 ReturnStatement
func TestReturnStatement(t *testing.T) {
	returnPrograms := []statements{
		{
			statement: &ReturnStatement{
				Token: token.Token{Type: token.RETURN, Literal: "return"},
				ReturnValue: &IntegerLiteral{
					Token: token.Token{Type: token.INT, Literal: "5"},
					Value: 5,
				},
			},
			expectedLiteral: "return",
			expectedString:  "return 5;",
		},
		{
			statement: &ReturnStatement{
				Token: token.Token{Type: token.RETURN, Literal: "return"},
				ReturnValue: &Identifier{
					Token: token.Token{Type: token.INT, Literal: "identifier"},
					Value: "identifier",
				},
			},
			expectedLiteral: "return",
			expectedString:  "return identifier;",
		},
		{
			statement: &ReturnStatement{
				Token: token.Token{Type: token.RETURN, Literal: "return"},
				ReturnValue: &PrefixExpression{
					Token:    token.Token{Type: token.INT, Literal: "5"},
					Operator: "-",
					Right: &IntegerLiteral{
						Token: token.Token{Type: token.INT, Literal: "5"},
						Value: 5,
					},
				},
			},
			expectedLiteral: "return",
			expectedString:  "return (-5);",
		},
		{
			statement: &ReturnStatement{
				Token: token.Token{Type: token.RETURN, Literal: "return"},
				ReturnValue: &InfixExpression{
					Token: token.Token{Type: token.INT, Literal: "5 + 5"},
					Left: &IntegerLiteral{
						Token: token.Token{Type: token.INT, Literal: "5"},
						Value: 5,
					},
					Operator: "+",
					Right: &IntegerLiteral{
						Token: token.Token{Type: token.INT, Literal: "5"},
						Value: 5,
					},
				},
			},
			expectedLiteral: "return",
			expectedString:  "return (5 + 5);",
		},
		{
			statement: &ReturnStatement{
				Token: token.Token{Type: token.RETURN, Literal: "return"},
				ReturnValue: &InfixExpression{
					Token: token.Token{Type: token.INT, Literal: "5 + 5"},
					Left: &IntegerLiteral{
						Token: token.Token{Type: token.INT, Literal: "5"},
						Value: 5,
					},
					Operator: "+",
					Right: &InfixExpression{
						Token: token.Token{Type: token.INT, Literal: "5"},
						Left: &IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
						Operator: "*",
						Right: &IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
					},
				},
			},
			expectedLiteral: "return",
			expectedString:  "return (5 + (5 * 5));",
		},
	}

	if !testStatement(t, returnPrograms) {
		return
	}
}

// 测试 ExpressionStatement
func TestExpressionStatement(t *testing.T) {
	expressionPrograms := []statements{
		{
			statement: &ExpressionStatement{
				Token: token.Token{Type: token.INT, Literal: "5"},
				Expression: &IntegerLiteral{
					Token: token.Token{Type: token.INT, Literal: "5"},
					Value: 5,
				},
			},
			expectedLiteral: "5",
			expectedString:  "5",
		},
		{
			statement: &ExpressionStatement{
				Token: token.Token{Type: token.INT, Literal: "-"},
				Expression: &PrefixExpression{
					Token:    token.Token{Type: token.MINUS, Literal: "-"},
					Operator: "-",
					Right: &IntegerLiteral{
						Token: token.Token{Type: token.INT, Literal: "5"},
						Value: 5,
					},
				},
			},
			expectedLiteral: "-",
			expectedString:  "(-5)",
		},
		{
			statement: &ExpressionStatement{
				Token: token.Token{Type: token.INT, Literal: "-"},
				Expression: &InfixExpression{
					Token: token.Token{Type: token.MINUS, Literal: "5"},
					Left: &IntegerLiteral{
						Token: token.Token{Type: token.INT, Literal: "5"},
						Value: 5,
					},
					Operator: "-",
					Right: &IntegerLiteral{
						Token: token.Token{Type: token.INT, Literal: "5"},
						Value: 5,
					},
				},
			},
			expectedLiteral: "-",
			expectedString:  "(5 - 5)",
		},
		{
			statement: &ExpressionStatement{
				Token: token.Token{Type: token.INT, Literal: "-"},
				Expression: &InfixExpression{
					Token: token.Token{Type: token.MINUS, Literal: "5"},
					Left: &IntegerLiteral{
						Token: token.Token{Type: token.INT, Literal: "5"},
						Value: 5,
					},
					Operator: "-",
					Right: &InfixExpression{
						Token: token.Token{Type: token.INT, Literal: "5"},
						Left: &IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
						Operator: "*",
						Right: &IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
					},
				},
			},
			expectedLiteral: "-",
			expectedString:  "(5 - (5 * 5))",
		},
	}

	if !testStatement(t, expressionPrograms) {
		return
	}
}

// 测试 BlockStatement
func TestBlockStatement(t *testing.T) {
	blocks := []statements{
		{
			statement:       getBlockStatement(),
			expectedLiteral: "{",
			expectedString:  "let myVar = myVar;let myVar = 5;return 5;",
		},
	}

	if !testStatement(t, blocks) {
		return
	}
}

// 测试 let 语句程序
func TestLetStatementProgram(t *testing.T) {
	lets := []programs{
		{
			program: Program{
				Statements: []Statement{
					getLetStatement("myVar", "myVar"),
				},
			},
			expectedString: "let myVar = myVar;",
		},
		{
			program: Program{
				Statements: []Statement{
					getLetStatement("myVar", "5"),
				},
			},
			expectedString: "let myVar = 5;",
		},
	}

	if !testProgram(t, lets) {
		return
	}
}

// 测试 return 语句程序
func TestReturnStatementProgram(t *testing.T) {
	returnPrograms := []programs{
		{
			program: Program{
				Statements: []Statement{
					&ReturnStatement{
						Token: token.Token{Type: token.RETURN, Literal: "return"},
						ReturnValue: &IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
					},
				},
			},
			expectedString: "return 5;",
		},
		{
			program: Program{
				Statements: []Statement{
					&ReturnStatement{
						Token: token.Token{Type: token.RETURN, Literal: "return"},
						ReturnValue: &Identifier{
							Token: token.Token{Type: token.INT, Literal: "ident"},
							Value: "ident",
						},
					},
				},
			},
			expectedString: "return ident;",
		},
		{
			program: Program{
				Statements: []Statement{
					&ReturnStatement{
						Token: token.Token{Type: token.RETURN, Literal: "return"},
						ReturnValue: &PrefixExpression{
							Token:    token.Token{Type: token.INT, Literal: "5"},
							Operator: "-",
							Right: &IntegerLiteral{
								Token: token.Token{Type: token.INT, Literal: "5"},
								Value: 5,
							},
						},
					},
				},
			},
			expectedString: "return (-5);",
		},
		{
			program: Program{
				Statements: []Statement{
					&ReturnStatement{
						Token: token.Token{Type: token.RETURN, Literal: "return"},
						ReturnValue: &InfixExpression{
							Token: token.Token{Type: token.INT, Literal: "5 + 5"},
							Left: &IntegerLiteral{
								Token: token.Token{Type: token.INT, Literal: "5"},
								Value: 5,
							},
							Operator: "+",
							Right: &IntegerLiteral{
								Token: token.Token{Type: token.INT, Literal: "5"},
								Value: 5,
							},
						},
					},
				},
			},
			expectedString: "return (5 + 5);",
		},
		{
			program: Program{
				Statements: []Statement{
					&ReturnStatement{
						Token: token.Token{Type: token.RETURN, Literal: "return"},
						ReturnValue: &InfixExpression{
							Token: token.Token{Type: token.INT, Literal: "5 + 5"},
							Left: &IntegerLiteral{
								Token: token.Token{Type: token.INT, Literal: "5"},
								Value: 5,
							},
							Operator: "+",
							Right: &InfixExpression{
								Token: token.Token{Type: token.INT, Literal: "5"},
								Left: &IntegerLiteral{
									Token: token.Token{Type: token.INT, Literal: "5"},
									Value: 5,
								},
								Operator: "*",
								Right: &IntegerLiteral{
									Token: token.Token{Type: token.INT, Literal: "5"},
									Value: 5,
								},
							},
						},
					},
				},
			},
			expectedString: "return (5 + (5 * 5));",
		},
	}

	if !testProgram(t, returnPrograms) {
		return
	}
}

// 测试表达式语句程序
func TestExpressionStatementProgram(t *testing.T) {
	expressionPrograms := []programs{
		{
			program: Program{
				Statements: []Statement{
					&ExpressionStatement{
						Token: token.Token{Type: token.INT, Literal: "5"},
						Expression: &IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
					},
				},
			},
			expectedString: "5",
		},
		{
			program: Program{
				Statements: []Statement{
					&ExpressionStatement{
						Token: token.Token{Type: token.INT, Literal: "-"},
						Expression: &PrefixExpression{
							Token:    token.Token{Type: token.MINUS, Literal: "-"},
							Operator: "-",
							Right: &IntegerLiteral{
								Token: token.Token{Type: token.INT, Literal: "5"},
								Value: 5,
							},
						},
					},
				},
			},
			expectedString: "(-5)",
		},
		{
			program: Program{
				Statements: []Statement{
					&ExpressionStatement{
						Token: token.Token{Type: token.INT, Literal: "5"},
						Expression: &InfixExpression{
							Token: token.Token{Type: token.MINUS, Literal: "5"},
							Left: &IntegerLiteral{
								Token: token.Token{Type: token.INT, Literal: "5"},
								Value: 5,
							},
							Operator: "-",
							Right: &IntegerLiteral{
								Token: token.Token{Type: token.INT, Literal: "5"},
								Value: 5,
							},
						},
					},
				},
			},
			expectedString: "(5 - 5)",
		},
		{
			program: Program{
				Statements: []Statement{
					&ExpressionStatement{
						Token: token.Token{Type: token.INT, Literal: "5"},
						Expression: &InfixExpression{
							Token: token.Token{Type: token.MINUS, Literal: "5"},
							Left: &IntegerLiteral{
								Token: token.Token{Type: token.INT, Literal: "5"},
								Value: 5,
							},
							Operator: "-",
							Right: &InfixExpression{
								Token: token.Token{Type: token.INT, Literal: "5"},
								Left: &IntegerLiteral{
									Token: token.Token{Type: token.INT, Literal: "5"},
									Value: 5,
								},
								Operator: "*",
								Right: &IntegerLiteral{
									Token: token.Token{Type: token.INT, Literal: "5"},
									Value: 5,
								},
							},
						},
					},
				},
			},
			expectedString: "(5 - (5 * 5))",
		},
	}

	if !testProgram(t, expressionPrograms) {
		return
	}
}

// 测试块语句程序
func TestBlockStatementProgram(t *testing.T) {
	blocks := []programs{
		{
			program: Program{
				Statements: []Statement{
					getBlockStatement(),
				},
			},
			expectedString: "let myVar = myVar;let myVar = 5;return 5;",
		},
	}

	if !testProgram(t, blocks) {
		return
	}
}

// 获取一个 BlockStatement 实例
func getBlockStatement() *BlockStatement {
	return &BlockStatement{
		Token: token.Token{Type: token.LBRACE, Literal: "{"},
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENTIFIER, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.LET, Literal: "let"},
					Value: "myVar",
				},
			},
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENTIFIER, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &IntegerLiteral{
					Token: token.Token{Type: token.INT, Literal: "5"},
					Value: 5,
				},
			},
			&ReturnStatement{
				Token: token.Token{Type: token.RETURN, Literal: "return"},
				ReturnValue: &IntegerLiteral{
					Token: token.Token{Type: token.INT, Literal: "5"},
					Value: 5,
				},
			},
		},
	}
}

// 获取一个 LetStatement 实例
func getLetStatement(name string, value string) *LetStatement {
	return &LetStatement{
		Token: token.Token{Type: token.LET, Literal: "let"},
		Name: &Identifier{
			Token: token.Token{Type: token.IDENTIFIER, Literal: name},
			Value: name,
		},
		Value: &Identifier{
			Token: token.Token{Type: token.IDENTIFIER, Literal: value},
			Value: value,
		},
	}
}

// 测试 Expression
func testExpression(t *testing.T, exps []expressions) bool {
	for _, expression := range exps {
		if expression.expression.TokenLiteral() != expression.expectedLiteral {
			t.Errorf("Expression wrong, expectedLiteral=%q, got=%q",
				expression.expectedLiteral, expression.expression.TokenLiteral())
		}
		if expression.expression.String() != expression.expectedString {
			t.Errorf("Expression wrong, expectedString=%q, got=%q",
				expression.expectedString, expression.expression.String())
		}
	}
	return false
}

// 测试 Statement
func testStatement(t *testing.T, sts []statements) bool {
	for _, st := range sts {
		if st.statement.TokenLiteral() != st.expectedLiteral {
			t.Errorf("Statement wrong, expectedLiteral=%q, got=%q",
				st.expectedLiteral, st.statement.TokenLiteral())
		}
		if st.statement.String() != st.expectedString {
			t.Errorf("Statement wrong, expectedString=%q, got=%q",
				st.expectedString, st.statement.String())
		}
	}
	return false
}

// 测试 Statement
func testProgram(t *testing.T, ps []programs) bool {
	for _, p := range ps {
		if p.program.String() != p.expectedString {
			t.Errorf("Statement wrong, expected=%q, got=%q",
				p.expectedString, p.program.String())
		}
	}
	return false
}
