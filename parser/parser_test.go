package parser

import (
	"holiya/ast"
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
		tokenLiteral  string
		expectedValue int64
		expectedString string
		expectError   bool
	}{
		{
			tokenLiteral: "123",
			expectedValue: 123,
			expectedString: "123",
			expectError:   false,
		},
		{
			tokenLiteral:  "0",
			expectedValue: 0,
			expectedString: "0",
			expectError:   false,
		},
		{
			tokenLiteral:  "-456",
			expectedValue: -456,
			expectedString: "-456",
			expectError:   false,
		},
		{
			// int64最大值
			tokenLiteral:  "9223372036854775807",
			expectedValue: 9223372036854775807,
			expectedString: "9223372036854775807",
			expectError:   false,
		},
		{
			tokenLiteral:  "abc",
			expectedValue: 0,
			expectedString: "",
			expectError:   true,
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
		tokenLiteral  string
		expectedValue float64
		expectedString string
		expectError   bool
	}{
		{
			tokenLiteral:  "123.456",
			expectedValue: 123.456,
			expectedString: "123.456",
			expectError:   false,
		},
		{
			tokenLiteral:  "0.0",
			expectedValue: 0.0,
			expectedString: "0.0",
			expectError:   false,
		},
		{
			tokenLiteral:  "-456.789",
			expectedValue: -456.789,
			expectedString: "-456.789",
			expectError:   false,
		},
		{
			tokenLiteral:  "3.141592653589793",
			expectedValue: 3.141592653589793,
			expectedString: "3.141592653589793",
			expectError:   false,
		},
		{
			tokenLiteral:  "10000000000",
			expectedValue: 10000000000,
			expectedString: "10000000000",
			expectError:   false,
		},
		{
			tokenLiteral:  "0.000123",
			expectedValue: 0.000123,
			expectedString: "0.000123",
			expectError:   false,
		},
		{
			tokenLiteral:  "abc",
			expectedValue: 0,
			expectedString: "",
			expectError:   true,
		},
		{
			tokenLiteral:  "",
			expectedValue: 0,
			expectedString: "",
			expectError:   true,
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