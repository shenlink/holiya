package token

import (
	"testing"
)

// 测试 Token 类型是否正确设置
func TestTokenType(t *testing.T) {
	tests := []struct {
		input    Token
		expected TokenType
	}{
		// 标识符
		{Token{Type: IDENTIFIER, Literal: "x"}, IDENTIFIER},

		// 数据类型
		{Token{Type: INT, Literal: "42"}, INT},
		{Token{Type: STRING, Literal: "hello"}, STRING},
		{Token{Type: FLOAT, Literal: "43.0"}, FLOAT},

		// 数学运算符
		{Token{Type: ASSIGN, Literal: "="}, ASSIGN},
		{Token{Type: PLUS, Literal: "+"}, PLUS},
		{Token{Type: MINUS, Literal: "-"}, MINUS},
		{Token{Type: MUL, Literal: "*"}, MUL},
		{Token{Type: DIV, Literal: "/"}, DIV},
		{Token{Type: MOD, Literal: "%"}, MOD},

		// 逻辑运算符
		{Token{Type: BAND, Literal: "!"}, BAND},
		{Token{Type: AND, Literal: "&&"}, AND},
		{Token{Type: OR, Literal: "||"}, OR},

		// 比较运算符
		{Token{Type: LT, Literal: "<"}, LT},
		{Token{Type: GT, Literal: ">"}, GT},
		{Token{Type: LTE, Literal: "<="}, LTE},
		{Token{Type: GTE, Literal: ">="}, GTE},
		{Token{Type: EQ, Literal: "=="}, EQ},
		{Token{Type: NEQ, Literal: "!="}, NEQ},

		// 分隔符
		{Token{Type: LPAREN, Literal: "("}, LPAREN},
		{Token{Type: RPAREN, Literal: ")"}, RPAREN},
		{Token{Type: LBRACE, Literal: "{"}, LBRACE},
		{Token{Type: RBRACE, Literal: "}"}, RBRACE},
		{Token{Type: LBRACKET, Literal: "["}, LBRACKET},
		{Token{Type: RBRACKET, Literal: "]"}, RBRACKET},
		{Token{Type: COMMA, Literal: ","}, COMMA},
		{Token{Type: SEMICOLON, Literal: ";"}, SEMICOLON},
		{Token{Type: COLON, Literal: ":"}, COLON},

		// 关键字
		{Token{Type: FUNCTION, Literal: "fn"}, FUNCTION},
		{Token{Type: LET, Literal: "let"}, LET},
		{Token{Type: TRUE, Literal: "true"}, TRUE},
		{Token{Type: FALSE, Literal: "false"}, FALSE},
		{Token{Type: IF, Literal: "if"}, IF},
		{Token{Type: ELSE, Literal: "else"}, ELSE},
		{Token{Type: RETURN, Literal: "return"}, RETURN},
	}

	for _, tt := range tests {
		if tt.input.Type != tt.expected {
			t.Errorf("Token.Type = %v, expected %v", tt.input.Type, tt.expected)
		}
	}
}

// 测试 Token 字面量是否正确设置
func TestTokenLiteral(t *testing.T) {
	tests := []struct {
		input    Token
		expected string
	}{
		// 标识符
		{Token{Type: IDENTIFIER, Literal: "x"}, "x"},

		// 数据类型
		{Token{Type: INT, Literal: "42"}, "42"},
		{Token{Type: FLOAT, Literal: "43.0"}, "43.0"},
		{Token{Type: STRING, Literal: "hello"}, "hello"},

		// 数学运算符
		{Token{Type: ASSIGN, Literal: "="}, "="},
		{Token{Type: PLUS, Literal: "+"}, "+"},
		{Token{Type: MINUS, Literal: "-"}, "-"},
		{Token{Type: MUL, Literal: "*"}, "*"},
		{Token{Type: DIV, Literal: "/"}, "/"},
		{Token{Type: MOD, Literal: "%"}, "%"},

		// 逻辑运算符
		{Token{Type: BAND, Literal: "!"}, "!"},
		{Token{Type: AND, Literal: "&&"}, "&&"},
		{Token{Type: OR, Literal: "||"}, "||"},

		// 比较运算符
		{Token{Type: LT, Literal: "<"}, "<"},
		{Token{Type: GT, Literal: ">"}, ">"},
		{Token{Type: LTE, Literal: "<="}, "<="},
		{Token{Type: GTE, Literal: ">="}, ">="},
		{Token{Type: EQ, Literal: "=="}, "=="},
		{Token{Type: NEQ, Literal: "!="}, "!="},

		// 分隔符
		{Token{Type: LPAREN, Literal: "("}, "("},
		{Token{Type: RPAREN, Literal: ")"}, ")"},
		{Token{Type: LBRACE, Literal: "{"}, "{"},
		{Token{Type: RBRACE, Literal: "}"}, "}"},
		{Token{Type: LBRACKET, Literal: "["}, "["},
		{Token{Type: RBRACKET, Literal: "]"}, "]"},
		{Token{Type: COMMA, Literal: ","}, ","},
		{Token{Type: SEMICOLON, Literal: ";"}, ";"},
		{Token{Type: COLON, Literal: ":"}, ":"},

		// 关键字
		{Token{Type: FUNCTION, Literal: "fn"}, "fn"},
		{Token{Type: LET, Literal: "let"}, "let"},
		{Token{Type: TRUE, Literal: "true"}, "true"},
		{Token{Type: FALSE, Literal: "false"}, "false"},
		{Token{Type: IF, Literal: "if"}, "if"},
		{Token{Type: ELSE, Literal: "else"}, "else"},
		{Token{Type: RETURN, Literal: "return"}, "return"},
	}

	for _, tt := range tests {
		if tt.input.Literal != tt.expected {
			t.Errorf("Token.Literal = %q, expected %q", tt.input.Literal, tt.expected)
		}
	}
}

// 测试 LookupIdentifier 是否能正确识别关键字
func TestLookupIdentifier(t *testing.T) {
	tests := []struct {
		input    string
		expected TokenType
	}{
		{"fn", FUNCTION},
		{"let", LET},
		{"true", TRUE},
		{"false", FALSE},
		{"if", IF},
		{"else", ELSE},
		{"return", RETURN},
		{"unknown", IDENTIFIER}, // 非关键字应返回 IDENTIFIER
	}

	for _, tt := range tests {
		result := LookupIdentifier(tt.input)
		if result != tt.expected {
			t.Errorf("LookupIdentifier(%q) = %v, expected %v", tt.input, result, tt.expected)
		}
	}
}
