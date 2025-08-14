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
