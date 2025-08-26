package parser

import (
	"fmt"
	"holiya/ast"
	"holiya/lexer"
	"holiya/token"
	"testing"
)

// 测试 parseIdentifier 函数
func TestParseIdentifier(t *testing.T) {
	tests := []struct {
		tokenType     token.TokenType
		tokenLiteral  string
		expectedValue string
	}{
		{
			tokenType:     token.IDENTIFIER,
			tokenLiteral:  "foobar",
			expectedValue: "foobar",
		},
		{
			tokenType:     token.IDENTIFIER,
			tokenLiteral:  "foo_bar",
			expectedValue: "foo_bar",
		},
		{
			tokenType:     token.IDENTIFIER,
			tokenLiteral:  "x123",
			expectedValue: "x123",
		},
		{
			tokenType:     token.IDENTIFIER,
			tokenLiteral:  "",
			expectedValue: "",
		},
	}

	// 执行测试用例
	for _, tt := range tests {
		tok := token.Token{
			Type:    tt.tokenType,
			Literal: tt.tokenLiteral,
		}
		parser := &Parser{
			currToken: tok,
		}
		result := parser.parseIdentifier()
		identifier, ok := result.(*ast.Identifier)
		if !ok {
			t.Fatalf("parseIdentifier() returned wrong type. Expected *ast.Identifier, got %T", result)
		}
		if identifier.TokenLiteral() != tt.tokenLiteral {
			t.Errorf("identifier.String() = %v, want %v", identifier.TokenLiteral(), tt.tokenLiteral)
		}
		if identifier.String() != tt.expectedValue {
			t.Errorf("identifier.String() = %v, want %v", identifier.String(), tt.expectedValue)
		}
	}
}

// 测试 parseIntegerLiteral 函数
func TestParseIntegerLiteral(t *testing.T) {
	tests := []struct {
		tokenLiteral   string
		expectedValue  int64
		expectedString string
		expectError    bool
	}{
		{
			tokenLiteral:   "123",
			expectedValue:  123,
			expectedString: "123",
			expectError:    false,
		},
		{
			tokenLiteral:   "0",
			expectedValue:  0,
			expectedString: "0",
			expectError:    false,
		},
		{
			tokenLiteral:   "-456",
			expectedValue:  -456,
			expectedString: "-456",
			expectError:    false,
		},
		{
			// int64最大值
			tokenLiteral:   "9223372036854775807",
			expectedValue:  9223372036854775807,
			expectedString: "9223372036854775807",
			expectError:    false,
		},
		{
			tokenLiteral:   "abc",
			expectedValue:  0,
			expectedString: "",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		tok := token.Token{
			Type:    token.INT,
			Literal: tt.tokenLiteral,
		}
		parser := &Parser{
			currToken: tok,
			errors:    []string{},
		}
		result := parser.parseIntegerLiteral()
		if tt.expectError {
			if len(parser.errors) == 0 {
				t.Errorf("expected error for input %s, but got none", tt.tokenLiteral)
			}
			if result != nil {
				t.Errorf("expected nil result for invalid input %s, but got %T", tt.tokenLiteral, result)
			}
			continue
		}

		// 验证没有错误
		if len(parser.errors) > 0 {
			t.Errorf("unexpected error for input %s: %v", tt.tokenLiteral, parser.errors)
			continue
		}

		integerLiteral, ok := result.(*ast.IntegerLiteral)
		if !ok {
			t.Fatalf("parseIntegerLiteral() returned wrong type. Expected *ast.IntegerLiteral, got %T", result)
		}
		if integerLiteral.TokenLiteral() != tt.tokenLiteral {
			t.Errorf("integerLiteral.TokenLiteral() = %v, want %v", integerLiteral.TokenLiteral(), tt.tokenLiteral)
		}
		if integerLiteral.String() != tt.expectedString {
			t.Errorf("integerLiteral.String() = %v, want %v", integerLiteral.String(), tt.expectedString)
		}
		if integerLiteral.Value != tt.expectedValue {
			t.Errorf("integerLiteral.Value = %v, want %v", integerLiteral.Value, tt.expectedValue)
		}
	}
}

// 测试 parseIntegerLiteral
func TestParseFloatLiteral(t *testing.T) {
	tests := []struct {
		tokenLiteral   string
		expectedValue  float64
		expectedString string
		expectError    bool
	}{
		{
			tokenLiteral:   "123.456",
			expectedValue:  123.456,
			expectedString: "123.456",
			expectError:    false,
		},
		{
			tokenLiteral:   "0.0",
			expectedValue:  0.0,
			expectedString: "0.0",
			expectError:    false,
		},
		{
			tokenLiteral:   "-456.789",
			expectedValue:  -456.789,
			expectedString: "-456.789",
			expectError:    false,
		},
		{
			tokenLiteral:   "3.141592653589793",
			expectedValue:  3.141592653589793,
			expectedString: "3.141592653589793",
			expectError:    false,
		},
		{
			tokenLiteral:   "10000000000",
			expectedValue:  10000000000,
			expectedString: "10000000000",
			expectError:    false,
		},
		{
			tokenLiteral:   "0.000123",
			expectedValue:  0.000123,
			expectedString: "0.000123",
			expectError:    false,
		},
		{
			tokenLiteral:   "abc",
			expectedValue:  0,
			expectedString: "",
			expectError:    true,
		},
		{
			tokenLiteral:   "",
			expectedValue:  0,
			expectedString: "",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		tok := token.Token{
			Type:    token.FLOAT,
			Literal: tt.tokenLiteral,
		}
		parser := &Parser{
			currToken: tok,
			errors:    []string{},
		}
		result := parser.parseFloatLiteral()
		if tt.expectError {
			if len(parser.errors) == 0 {
				t.Errorf("expected error for input %s, but got none", tt.tokenLiteral)
			}
			if result != nil {
				t.Errorf("expected nil result for invalid input %s, but got %T", tt.tokenLiteral, result)
			}
			continue
		}

		// 验证没有错误
		if len(parser.errors) > 0 {
			t.Errorf("unexpected error for input %s: %v", tt.tokenLiteral, parser.errors)
			continue
		}

		floatLiteral, ok := result.(*ast.FloatLiteral)
		if !ok {
			t.Fatalf("parseFloatLiteral() returned wrong type. Expected *ast.FloatLiteral, got %T", result)
		}
		if floatLiteral.TokenLiteral() != tt.tokenLiteral {
			t.Errorf("floatLiteral.TokenLiteral() = %v, want %v", floatLiteral.TokenLiteral(), tt.tokenLiteral)
		}
		if floatLiteral.String() != tt.expectedString {
			t.Errorf("floatLiteral.TokenLiteral() = %v, want %v", floatLiteral.String(), tt.expectedString)
		}
		if floatLiteral.Value != tt.expectedValue {
			t.Errorf("floatLiteral.Value = %v, want %v", floatLiteral.Value, tt.expectedValue)
		}
	}
}

// 测试 parseStringLiteral 函数
func TestParseStringLiteral(t *testing.T) {
	tests := []struct {
		tokenLiteral  string
		expectedValue string
	}{
		{
			tokenLiteral:  "hello world",
			expectedValue: "hello world",
		},
		{
			tokenLiteral:  "",
			expectedValue: "",
		},
		{
			tokenLiteral:  "123",
			expectedValue: "123",
		},
		{
			tokenLiteral:  "foo_bar",
			expectedValue: "foo_bar",
		},
		{
			tokenLiteral:  "hello\nworld",
			expectedValue: "hello\nworld",
		},
	}

	// 执行测试用例
	for _, tt := range tests {
		tok := token.Token{
			Type:    token.STRING,
			Literal: tt.tokenLiteral,
		}
		parser := &Parser{
			currToken: tok,
		}
		result := parser.parseStringLiteral()
		stringLiteral, ok := result.(*ast.StringLiteral)
		if !ok {
			t.Fatalf("parseStringLiteral() returned wrong type. Expected *ast.StringLiteral, got %T", result)
		}
		if stringLiteral.TokenLiteral() != tt.tokenLiteral {
			t.Errorf("stringLiteral.TokenLiteral() = %v, want %v", stringLiteral.TokenLiteral(), tt.tokenLiteral)
		}
		if stringLiteral.String() != tt.expectedValue {
			t.Errorf("stringLiteral.String() = %v, want %v", stringLiteral.String(), tt.expectedValue)
		}
	}
}

// 测试 parsePrefixExpression 函数
func TestParsePrefixExpression(t *testing.T) {
	tests := []struct {
		input          string
		token          token.Token
		operator       string
		operatorType   token.TokenType
		rightLiteral   string
		rightType      token.TokenType
		expectedString string
	}{
		{
			input:          "!true",
			token:          token.Token{Type: token.BAND, Literal: "!"},
			operator:       "!",
			operatorType:   token.BAND,
			rightLiteral:   "true",
			rightType:      token.FALSE,
			expectedString: "(!true)",
		},
		{
			input:          "!false",
			token:          token.Token{Type: token.BAND, Literal: "!"},
			operator:       "!",
			operatorType:   token.BAND,
			rightLiteral:   "false",
			rightType:      token.TRUE,
			expectedString: "(!false)",
		},
		{
			input:          "-10",
			token:          token.Token{Type: token.MINUS, Literal: "-"},
			operator:       "-",
			operatorType:   token.MINUS,
			rightLiteral:   "10",
			rightType:      token.INT,
			expectedString: "(-10)",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		parser := New(l)
		result := parser.parsePrefixExpression()
		prefixExpression, ok := result.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("parsePrefixExpression() returned wrong type. Expected *ast.PrefixExpression, got %T", result)
		}

		if prefixExpression.Token != tt.token {
			t.Errorf("prefixExpression.Token = %v, want %v", prefixExpression.Token, tt.token)
		}

		if prefixExpression.Operator != tt.operator {
			t.Errorf("prefixExpression.Operator = %v, want %v", prefixExpression.Operator, tt.operator)
		}

		if prefixExpression.Right.TokenLiteral() != tt.rightLiteral {
			t.Errorf("prefixExpression..Right.TokenLiteral() = %v, want %v", prefixExpression.Operator, tt.rightLiteral)
		}

		if prefixExpression.String() != tt.expectedString {
			t.Errorf("prefixExpression.String() = %v, want %v", prefixExpression.String(), tt.expectedString)
		}
	}
}

// 测试 parseExpression 函数
func TestParseExpression(t *testing.T) {
	tests := []struct {
		input          string
		token          token.Token
		operator       string
		operatorType   token.TokenType
		rightLiteral   string
		rightType      token.TokenType
		expectedString string
	}{
		{
			input:          "!true",
			token:          token.Token{Type: token.BAND, Literal: "!"},
			operator:       "!",
			operatorType:   token.BAND,
			rightLiteral:   "true",
			rightType:      token.FALSE,
			expectedString: "(!true)",
		},
		{
			input:          "!false",
			token:          token.Token{Type: token.BAND, Literal: "!"},
			operator:       "!",
			operatorType:   token.BAND,
			rightLiteral:   "false",
			rightType:      token.TRUE,
			expectedString: "(!false)",
		},
		{
			input:          "-10",
			token:          token.Token{Type: token.MINUS, Literal: "-"},
			operator:       "-",
			operatorType:   token.MINUS,
			rightLiteral:   "10",
			rightType:      token.INT,
			expectedString: "(-10)",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		parser := New(l)
		result := parser.parseExpression(PREFIX)
		prefixExpression, ok := result.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("parsePrefixExpression() returned wrong type. Expected *ast.PrefixExpression, got %T", result)
		}

		if prefixExpression.Token != tt.token {
			t.Errorf("prefixExpression.Token = %v, want %v", prefixExpression.Token, tt.token)
		}

		if prefixExpression.Operator != tt.operator {
			t.Errorf("prefixExpression.Operator = %v, want %v", prefixExpression.Operator, tt.operator)
		}

		if prefixExpression.Right.TokenLiteral() != tt.rightLiteral {
			t.Errorf("prefixExpression..Right.TokenLiteral() = %v, want %v", prefixExpression.Operator, tt.rightLiteral)
		}

		if prefixExpression.String() != tt.expectedString {
			t.Errorf("prefixExpression.String() = %v, want %v", prefixExpression.String(), tt.expectedString)
		}
	}
}

// 测试 noPrefixParseFnError 函数
func TestNoPrefixParseFnError(t *testing.T) {
	parser := Parser{}
	parser.noPrefixParseFnError(token.BAND)
	if len(parser.errors) != 1 {
		t.Errorf("parser.errors = %v, want 1", parser.errors)
	}
	if parser.errors[0] != "no prefix parse function for ! found" {
		t.Errorf("parser.errors[0] = %v, want 'no prefix parse function for BAND found'", parser.errors[0])
	}
}

// 测试 peekTokenIs 函数
func TestPeekTokenIs(t *testing.T) {
	l := lexer.New("!true")
	parser := New(l)
	if !parser.peekTokenIs(token.TRUE) {
		t.Errorf("parser.peekTokenIs(token.TRUE) = true, want false")
	}
}

// 测试 currPrecedence 函数
func TestCurrPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"=", LOWEST},
		{"==", EQUALS},
		{"!=", EQUALS},
		{"+", SUM},
		{"-", SUM},
		{"*", PRODUCT},
		{"/", PRODUCT},
		{"<", LESSGREATER},
		{">", LESSGREATER},
		{"(", CALL},
		{"[", INDEX},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		parser := New(l)

		actual := parser.currPrecedence()
		if actual != tt.expected {
			t.Errorf("currPrecedence() for input %q = %d, want %d", tt.input, actual, tt.expected)
		}
	}
}

// 测试 parseBoolean 函数
func TestParseBoolean(t *testing.T) {
	tests := []struct {
		tokenLiteral   string
		expectedValue  bool
		expectedString string
	}{
		{
			tokenLiteral:   "true",
			expectedValue:  true,
			expectedString: "true",
		},
		{
			tokenLiteral:   "false",
			expectedValue:  false,
			expectedString: "false",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.tokenLiteral)
		parser := New(l)
		result := parser.parseBoolean()
		booleanLiteral, ok := result.(*ast.Boolean)
		if !ok {
			t.Fatalf("parseBoolean() returned wrong type. Expected *ast.Boolean, got %T", result)
		}
		if booleanLiteral.TokenLiteral() != tt.tokenLiteral {
			t.Errorf("booleanLiteral.TokenLiteral() = %v, want %v", booleanLiteral.TokenLiteral(), tt.tokenLiteral)
		}
		if booleanLiteral.Value != tt.expectedValue {
			t.Errorf("integerLiteral.Value = %v, want %v", booleanLiteral.Value, tt.expectedValue)
		}
		if booleanLiteral.String() != tt.expectedString {
			t.Errorf("booleanLiteral.String() = %v, want %v", booleanLiteral.String(), tt.expectedString)
		}
	}
}

// 测试 parseGroupedExpression 函数
func TestParseGroupedExpression(t *testing.T) {
	tests := []struct {
		input          string
		expectedString string
	}{
		{
			input:          "(1)",
			expectedString: "1",
		},
		{
			input:          "(x)",
			expectedString: "x",
		},
		{
			input:          "(true)",
			expectedString: "true",
		},
		{
			input:          "(false)",
			expectedString: "false",
		},
		{
			input:          "(1 + 2)",
			expectedString: "(1 + 2)",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		parser := New(l)
		result := parser.parseGroupedExpression()
		if result == nil {
			t.Errorf("parseGroupedExpression() returned nil for input %q", tt.input)
			continue
		}
		if result.String() != tt.expectedString {
			t.Errorf("parseGroupedExpression() for input %q = %v, want %v", 
				tt.input, result.String(), tt.expectedString)
		}
	}
}

// 测试 parseIfExpression 函数
func TestParseIfExpression(t *testing.T) {
	tests := []struct {
		input               string
		expectedCondition   string
		expectedConsequence string
		hasAlternative      bool
		expectedAlternative string
		expectError         bool
	}{
		{
			input:               `if (x < y) { x; }`,
			expectedCondition:   "(x < y)",
			expectedConsequence: "x",
			hasAlternative:      false,
			expectedAlternative: "",
			expectError:         false,
		},
		{
			input:               `if (true) { 10; }`,
			expectedCondition:   "true",
			expectedConsequence: "10",
			hasAlternative:      false,
			expectedAlternative: "",
			expectError:         false,
		},
		{
			input:               `if (false) { 10; } else { 20; }`,
			expectedCondition:   "false",
			expectedConsequence: "10",
			hasAlternative:      true,
			expectedAlternative: "20",
			expectError:         false,
		},
		{
			input:               `if (10 > 5) { return true; }`,
			expectedCondition:   "(10 > 5)",
			expectedConsequence: "return true;",
			hasAlternative:      false,
			expectedAlternative: "",
			expectError:         false,
		},
		{
			input:               `if (10 > 5) { let x = 10; } else { let y = 20; }`,
			expectedCondition:   "(10 > 5)",
			expectedConsequence: "let x = 10;",
			hasAlternative:      true,
			expectedAlternative: "let y = 20;",
			expectError:         false,
		},
		{
			input:               `if (10 > 5) { return true; } else { return false; }`,
			expectedCondition:   "(10 > 5)",
			expectedConsequence: "return true;",
			hasAlternative:      true,
			expectedAlternative: "return false;",
			expectError:         false,
		},
		// 错误情况测试 - 缺少条件括号
		{
			input:               `if x < y { x; }`,
			expectedCondition:   "",
			expectedConsequence: "",
			hasAlternative:      false,
			expectedAlternative: "",
			expectError:         true,
		},
		// 错误情况测试 - 缺少条件右括号
		{
			input:               `if (x < y { x; }`,
			expectedCondition:   "",
			expectedConsequence: "",
			hasAlternative:      false,
			expectedAlternative: "",
			expectError:         true,
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		parser := New(l)
		result := parser.parseIfExpression()

		if tt.expectError {
			if len(parser.errors) == 0 {
				t.Errorf("expected error for input %s, but got none", tt.input)
			}
			if result != nil {
				t.Errorf("expected nil result for invalid input %s, but got %T", tt.input, result)
			}
			continue
		}

		// 验证没有错误
		if len(parser.errors) > 0 {
			t.Errorf("unexpected error for input %s: %v", tt.input, parser.errors)
			continue
		}

		ifResult, ok := result.(*ast.IfExpression)
		if !ok {
			t.Fatalf("parseIfExpression() returned wrong type. Expected *ast.IfExpression, got %T", result)
		}

		if ifResult.Condition.String() != tt.expectedCondition {
			t.Errorf("ifResult.Condition.String() = %v, want %v", ifResult.Condition.String(), tt.expectedCondition)
		}

		if ifResult.Consequence.String() != tt.expectedConsequence {
			t.Errorf("ifResult.Consequence.String() = %v, want %v", ifResult.Consequence.String(), tt.expectedConsequence)
		}

		if tt.hasAlternative {
			if ifResult.Alternative == nil {
				t.Errorf("expected alternative block, but got nil")
			}
			if ifResult.Alternative.String() != tt.expectedAlternative {
				t.Errorf("ifResult.Alternative.String() = %v, want %v", ifResult.Alternative.String(), tt.expectedAlternative)
			}
		} else {
			if ifResult.Alternative != nil {
				t.Errorf("expected no alternative block, but got %v", ifResult.Alternative.String())
			}
		}
	}
}

// 测试 parseBlockStatement 函数
func TestParseBlockStatement(t *testing.T) {
	tests := []struct {
		input              string
		expectedStatements []string
		expectError        bool
	}{
		{
			input:              `{ x; }`,
			expectedStatements: []string{"x"},
			expectError:        false,
		},
		{
			input:              `{ 10; }`,
			expectedStatements: []string{"10"},
			expectError:        false,
		},
		{
			input:              `{ true; }`,
			expectedStatements: []string{"true"},
			expectError:        false,
		},
		{
			input:              `{ x; y; z; }`,
			expectedStatements: []string{"x", "y", "z"},
			expectError:        false,
		},
		{
			input:              `{ 10; 20; 30; }`,
			expectedStatements: []string{"10", "20", "30"},
			expectError:        false,
		},
		{
			input:              `{ return 10; }`,
			expectedStatements: []string{"return 10;"},
			expectError:        false,
		},
		{
			input:              `{ return 10; return 20; }`,
			expectedStatements: []string{"return 10;", "return 20;"},
			expectError:        false,
		},
		{
			input:              `{ let x = 10; x; }`,
			expectedStatements: []string{"let x = 10;", "x"},
			expectError:        false,
		},
		{
			input:              `{ let x = 10; let y = 20; x + y; }`,
			expectedStatements: []string{"let x = 10;", "let y = 20;", "(x + y)"},
			expectError:        false,
		},
		// 空块语句测试
		{
			input:              `{ }`,
			expectedStatements: []string{},
			expectError:        false,
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		parser := New(l)

		// 手动设置currToken为LBRACE，因为parseBlockStatement期望当前token是{
		if !parser.currTokenIs(token.LBRACE) {
			// 找到第一个{
			for !parser.currTokenIs(token.LBRACE) && !parser.currTokenIs(token.EOF) {
				parser.nextToken()
			}
		}

		blockStatement := parser.parseBlockStatement()

		if tt.expectError {
			if len(parser.errors) == 0 {
				t.Errorf("expected error for input %s, but got none", tt.input)
			}
			if blockStatement != nil {
				t.Errorf("expected nil blockStatement for invalid input %s, but got %T", tt.input, blockStatement)
			}
			continue
		}

		// 验证没有错误
		if len(parser.errors) > 0 {
			t.Errorf("unexpected error for input %s: %v", tt.input, parser.errors)
			continue
		}

		if len(blockStatement.Statements) != len(tt.expectedStatements) {
			t.Errorf("blockStatement.Statements length = %d, want %d", len(blockStatement.Statements), len(tt.expectedStatements))
			continue
		}

		for i, expectedStmt := range tt.expectedStatements {
			actualStmt := blockStatement.Statements[i].String()
			if actualStmt != expectedStmt {
				t.Errorf("blockStatement.Statements[%d].String() = %v, want %v", i, actualStmt, expectedStmt)
			}
		}
	}
}

// 测试 parseStatement 函数
func TestParseStatement(t *testing.T) {
	tests := []struct {
		input          string
		expectedType   string
		expectedString string
		expectError    bool
	}{
		// 测试let语句
		{
			input:          "let x = 5;",
			expectedType:   "let",
			expectedString: "let x = 5;",
			expectError:    false,
		},
		{
			input:          "let y = 10;",
			expectedType:   "let",
			expectedString: "let y = 10;",
			expectError:    false,
		},
		{
			input:          "let foobar = 123;",
			expectedType:   "let",
			expectedString: "let foobar = 123;",
			expectError:    false,
		},
		// 测试return语句
		{
			input:          "return 5;",
			expectedType:   "return",
			expectedString: "return 5;",
			expectError:    false,
		},
		{
			input:          "return x;",
			expectedType:   "return",
			expectedString: "return x;",
			expectError:    false,
		},
		{
			input:          "return 10 + 5;",
			expectedType:   "return",
			expectedString: "return (10 + 5);",
			expectError:    false,
		},
		// 测试表达式语句
		{
			input:          "x;",
			expectedType:   "expression",
			expectedString: "x",
			expectError:    false,
		},
		{
			input:          "123;",
			expectedType:   "expression",
			expectedString: "123",
			expectError:    false,
		},
		{
			input:          "x + y;",
			expectedType:   "expression",
			expectedString: "(x + y)",
			expectError:    false,
		},
		{
			input:          "10 + 5;",
			expectedType:   "expression",
			expectedString: "(10 + 5)",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		parser := New(l)
		statement := parser.parseStatement()

		if tt.expectError {
			if len(parser.errors) == 0 {
				t.Errorf("expected error for input %s, but got none", tt.input)
			}
			if statement != nil {
				t.Errorf("expected nil statement for invalid input %s, but got %T", tt.input, statement)
			}
			continue
		}

		// 验证没有错误
		if len(parser.errors) > 0 {
			t.Errorf("unexpected error for input %s: %v", tt.input, parser.errors)
			continue
		}

		if statement == nil {
			t.Errorf("parseStatement() returned nil for input %s", tt.input)
			continue
		}

		switch tt.expectedType {
		case "let":
			if _, ok := statement.(*ast.LetStatement); !ok {
				t.Errorf("parseStatement() returned wrong type. Expected *ast.LetStatement, got %T", statement)
			}
		case "return":
			if _, ok := statement.(*ast.ReturnStatement); !ok {
				t.Errorf("parseStatement() returned wrong type. Expected *ast.ReturnStatement, got %T", statement)
			}
		case "expression":
			if _, ok := statement.(*ast.ExpressionStatement); !ok {
				t.Errorf("parseStatement() returned wrong type. Expected *ast.ExpressionStatement, got %T", statement)
			}
		}

		if statement.String() != tt.expectedString {
			t.Errorf("statement.String() = %v, want %v", statement.String(), tt.expectedString)
		}
	}
}

// 测试 parseLetStatement 函数
func TestParseLetStatement(t *testing.T) {
	tests := []struct {
		input         string
		expectedName  string
		expectedValue string
		expectError   bool
	}{
		{
			input:         "let x = 5;",
			expectedName:  "x",
			expectedValue: "5",
			expectError:   false,
		},
		{
			input:         "let y = true;",
			expectedName:  "y",
			expectedValue: "true",
			expectError:   false,
		},
		{
			input:         "let foobar = y;",
			expectedName:  "foobar",
			expectedValue: "y",
			expectError:   false,
		},
		{
			input:         "let x = 10 + 5;",
			expectedName:  "x",
			expectedValue: "(10 + 5)",
			expectError:   false,
		},
		{
			input:         "let x = 10 * 5;",
			expectedName:  "x",
			expectedValue: "(10 * 5)",
			expectError:   false,
		},
		{
			input:         "let x = 10 / 5;",
			expectedName:  "x",
			expectedValue: "(10 / 5)",
			expectError:   false,
		},
		{
			input:         "let x = 10 - 5;",
			expectedName:  "x",
			expectedValue: "(10 - 5)",
			expectError:   false,
		},
		{
			input:         "let x = 10 + 5;",
			expectedName:  "x",
			expectedValue: "(10 + 5)",
			expectError:   false,
		},
		{
			input:         "let x = -5;",
			expectedName:  "x",
			expectedValue: "(-5)",
			expectError:   false,
		},
		{
			input:         "let x = !true;",
			expectedName:  "x",
			expectedValue: "(!true)",
			expectError:   false,
		},
		// 错误情况测试 - 缺少标识符
		{
			input:         "let = 5;",
			expectedName:  "",
			expectedValue: "",
			expectError:   true,
		},
		// 错误情况测试 - 缺少赋值操作符
		{
			input:         "let x 5;",
			expectedName:  "",
			expectedValue: "",
			expectError:   true,
		},
		// 错误情况测试 - 缺少分号
		{
			input:         "let x = 5",
			expectedName:  "",
			expectedValue: "",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		parser := New(l)
		letStatement := parser.parseLetStatement()

		if tt.expectError {
			if len(parser.errors) == 0 {
				t.Errorf("expected error for input %s, but got none", tt.input)
			}
			if letStatement != nil {
				t.Errorf("expected nil letStatement for invalid input %s, but got %T", tt.input, letStatement)
			}
			continue
		}

		// 验证没有错误
		if len(parser.errors) > 0 {
			t.Errorf("unexpected error for input %s: %v", tt.input, parser.errors)
			continue
		}

		if letStatement.Name.String() != tt.expectedName {
			t.Errorf("letStatement.Name.String() = %v, want %v", letStatement.Name.String(), tt.expectedName)
		}

		if letStatement.Value.String() != tt.expectedValue {
			t.Errorf("letStatement.Value.String() = %v, want %v", letStatement.Value.String(), tt.expectedValue)
		}

		if letStatement.String() != fmt.Sprintf("let %s = %s;", tt.expectedName, tt.expectedValue) {
			t.Errorf("letStatement.String() = %v, want %v", letStatement.String(),
				fmt.Sprintf("let %s = %s;", tt.expectedName, tt.expectedValue))
		}
	}
}

// 测试 parseReturnStatement 函数
func TestParseReturnStatement(t *testing.T) {
	tests := []struct {
		input               string
		expectedReturnValue string
		expectError         bool
	}{
		{
			input:               "return 5;",
			expectedReturnValue: "5",
			expectError:         false,
		},
		{
			input:               "return true;",
			expectedReturnValue: "true",
			expectError:         false,
		},
		{
			input:               "return x;",
			expectedReturnValue: "x",
			expectError:         false,
		},
		{
			input:               "return 10 + 5;",
			expectedReturnValue: "(10 + 5)",
			expectError:         false,
		},
		{
			input:               "return 10 * 5;",
			expectedReturnValue: "(10 * 5)",
			expectError:         false,
		},
		{
			input:               "return 10 / 5;",
			expectedReturnValue: "(10 / 5)",
			expectError:         false,
		},
		{
			input:               "return 10 - 5;",
			expectedReturnValue: "(10 - 5)",
			expectError:         false,
		},
		{
			input:               "return -5;",
			expectedReturnValue: "(-5)",
			expectError:         false,
		},
		{
			input:               "return !true;",
			expectedReturnValue: "(!true)",
			expectError:         false,
		},
		{
			input:               "return 10 + 5;",
			expectedReturnValue: "(10 + 5)",
			expectError:         false,
		},
		// 错误情况测试 - 缺少分号
		{
			input:               "return 5",
			expectedReturnValue: "",
			expectError:         true,
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		parser := New(l)
		returnStatement := parser.parseReturnStatement()

		if tt.expectError {
			if len(parser.errors) == 0 {
				t.Errorf("expected error for input %s, but got none", tt.input)
			}
			if returnStatement != nil {
				t.Errorf("expected nil returnStatement for invalid input %s, but got %T", tt.input, returnStatement)
			}
			continue
		}

		// 验证没有错误
		if len(parser.errors) > 0 {
			t.Errorf("unexpected error for input %s: %v", tt.input, parser.errors)
			continue
		}

		if returnStatement.ReturnValue.String() != tt.expectedReturnValue {
			t.Errorf("returnStatement.ReturnValue.String() = %v, want %v", returnStatement.ReturnValue.String(), tt.expectedReturnValue)
		}

		if returnStatement.String() != fmt.Sprintf("return %s;", tt.expectedReturnValue) {
			t.Errorf("returnStatement.String() = %v, want %v", returnStatement.String(),
				fmt.Sprintf("return %s;", tt.expectedReturnValue))
		}
	}
}

// 测试 parseExpressionStatement 函数
func TestParseExpressionStatement(t *testing.T) {
	tests := []struct {
		input              string
		expectedExpression string
		expectError        bool
	}{
		{
			input:              "x;",
			expectedExpression: "x",
			expectError:        false,
		},
		{
			input:              "123;",
			expectedExpression: "123",
			expectError:        false,
		},
		{
			input:              "true;",
			expectedExpression: "true",
			expectError:        false,
		},
		{
			input:              "false;",
			expectedExpression: "false",
			expectError:        false,
		},
		{
			input:              "x + y;",
			expectedExpression: "(x + y)",
			expectError:        false,
		},
		{
			input:              "10 + 5;",
			expectedExpression: "(10 + 5)",
			expectError:        false,
		},
		{
			input:              "10 * 5;",
			expectedExpression: "(10 * 5)",
			expectError:        false,
		},
		{
			input:              "10 / 5;",
			expectedExpression: "(10 / 5)",
			expectError:        false,
		},
		{
			input:              "10 - 5;",
			expectedExpression: "(10 - 5)",
			expectError:        false,
		},
		{
			input:              "-5;",
			expectedExpression: "(-5)",
			expectError:        false,
		},
		{
			input:              "!true;",
			expectedExpression: "(!true)",
			expectError:        false,
		},
		{
			input:              "10 + 5;",
			expectedExpression: "(10 + 5)",
			expectError:        false,
		},
		{
			input:              "x * y + z;",
			expectedExpression: "((x * y) + z)",
			expectError:        false,
		},
		{
			input:              "x + y * z;",
			expectedExpression: "(x + (y * z))",
			expectError:        false,
		},
		// 错误情况测试 - 缺少分号
		{
			input:              "x",
			expectedExpression: "",
			expectError:        true,
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		parser := New(l)
		expressionStatement := parser.parseExpressionStatement()

		if tt.expectError {
			if len(parser.errors) == 0 {
				t.Errorf("expected error for input %s, but got none", tt.input)
			}
			if expressionStatement != nil {
				t.Errorf("expected nil expressionStatement for invalid input %s, but got %T", tt.input, expressionStatement)
			}
			continue
		}

		// 验证没有错误
		if len(parser.errors) > 0 {
			t.Errorf("unexpected error for input %s: %v", tt.input, parser.errors)
			continue
		}

		if expressionStatement.Expression.String() != tt.expectedExpression {
			t.Errorf("expressionStatement.Expression.String() = %v, want %v", expressionStatement.Expression.String(), tt.expectedExpression)
		}

		if expressionStatement.String() != tt.expectedExpression {
			t.Errorf("expressionStatement.String() = %v, want %v", expressionStatement.String(), tt.expectedExpression)
		}
	}
}

// 测试 expectPeek 函数
func TestExpectPeek(t *testing.T) {
	l := lexer.New("!true")
	parser := New(l)
	if !parser.expectPeek(token.TRUE) {
		t.Errorf("parser.expectPeek(token.TRUE) = true, want false")
	}
}

// 测试 peekError 函数
func TestPeekError(t *testing.T) {
	l := lexer.New("!true")
	parser := New(l)
	parser.peekError(token.FALSE)
	if len(parser.errors) == 0 {
		t.Error("expected 1 error, but got 0")
	}
	if parser.errors[0] != "expected next token to be FALSE, got TRUE instead" {
		t.Errorf("parser.errors[0] = %v, want 'expected next token to be FALSE, got TRUE instead'", parser.errors[0])
	}
}