package lexer

import (
	"holiya/token"
	"reflect"
	"testing"
)

// TestNew 测试 New 函数
func TestNew(t *testing.T) {
	// 测试用例 1: 空输入
	l := New("")
	expected := &Lexer{
		input:        "",
		position:     0,
		nextPosition: 1,
		ch:           0, // EOF
	}
	if !reflect.DeepEqual(l, expected) {
		t.Errorf("Expected %+v, got %+v", expected, l)
	}

	// 测试用例 2: 单个字符输入
	l = New("a")
	expected = &Lexer{
		input:        "a",
		position:     0,
		nextPosition: 1,
		ch:           'a',
	}
	if !reflect.DeepEqual(l, expected) {
		t.Errorf("Expected %+v, got %+v", expected, l)
	}

	// 测试用例 3: 多字符输入
	l = New("abc")
	expected = &Lexer{
		input:        "abc",
		position:     0,
		nextPosition: 1,
		ch:           'a',
	}
	if !reflect.DeepEqual(l, expected) {
		t.Errorf("Expected %+v, got %+v", expected, l)
	}
}

// TestNextToken 测试获取下一个 token 词法单元的功能
func TestNextToken(t *testing.T) {
	input := `
// 定义变量 
let five = 5;
// 定义变量
let ten1 = 10;

// 定义函数
let add = fn(x, y) {
  x + y;
};

let result = add(five, ten1);
! - / * 5;
5 < 10 > 5;

if (5 < 10) {
	return true;
} else {
	return false;
}

10 == 10;
10 != 9;
>= 10.5 [ ] && || <= : % "hello world"
& 
`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENTIFIER, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIER, "ten1"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIER, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "x"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENTIFIER, "x"},
		{token.PLUS, "+"},
		{token.IDENTIFIER, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIER, "result"},
		{token.ASSIGN, "="},
		{token.IDENTIFIER, "add"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "five"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "ten1"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.BAND, "!"},
		{token.MINUS, "-"},
		{token.DIV, "/"},
		{token.MUL, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.NEQ, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},
		{token.GTE, ">="},
		{token.FLOAT, "10.5"},
		{token.LBRACKET, "["},
		{token.RBRACKET, "]"},
		{token.AND, "&&"},
		{token.OR, "||"},
		{token.LTE, "<="},
		{token.COLON, ":"},
		{token.MOD, "%"},
		{token.STRING, "hello world"},
		{token.ILLEGAL, "&"},
		{token.EOF, ""},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - token type wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// TestSkipWhitespace 测试跳过空白字符的功能
func TestSkipWhitespace(t *testing.T) {
	tests := []struct {
		input                string
		expectedPosition     int
		expectedNextPosition int
		expectedChar         byte
	}{
		{
			// 多个空格后接 'a'
			input:                "   a",
			expectedPosition:     3,
			expectedNextPosition: 4,
			expectedChar:         'a',
		},
		{
			// 制表符后接 'a'
			input:                "\tb",
			expectedPosition:     1,
			expectedNextPosition: 2,
			expectedChar:         'b',
		},
		{
			// 两个换行符后接 'a'
			input:                "\n\nc",
			expectedPosition:     2,
			expectedNextPosition: 3,
			expectedChar:         'c',
		},
		{
			// Windows风格换行后接空格和 'a'
			input:                "\r\n d",
			expectedPosition:     3,
			expectedNextPosition: 4,
			expectedChar:         'd',
		},
		{
			// 全部是空格
			input:                "     ",
			expectedPosition:     5,
			expectedNextPosition: 6,
			// EOF
			expectedChar: END,
		},
	}

	for i, tt := range tests {
		l := New(tt.input)
		l.skipWhitespace()
		if l.position != tt.expectedPosition || l.nextPosition != tt.expectedNextPosition || l.ch != tt.expectedChar {
			t.Errorf("Test %d failed. Expected position=%d, nextPosition=%d, ch='%c'; But got position=%d, nextPosition=%d, ch='%c'",
				i, tt.expectedPosition, tt.expectedNextPosition, tt.expectedChar, l.position, l.nextPosition, l.ch)
		}
	}
}

// TestSkipComment 测试跳过注释的功能
func TestSkipComment(t *testing.T) {
	tests := []struct {
		input        string
		expectedChar byte
	}{
		{
			// 单行注释到行尾
			input:        "// 这是一个注释\nx",
			expectedChar: 'x',
		},
		{
			// 单行注释后接单行注释
			input:        "// 注释1\n// 注释2\n x",
			expectedChar: 'x',
		},
		{
			// 文件末尾的单行注释
			input:        "// 末尾注释",
			expectedChar: END, // EOF
		},
		{
			// 注释中包含特殊字符
			input:        "// 注释!@#$%^&*()\nx",
			expectedChar: 'x',
		},
		{
			// 注释前后有空格
			input:        "//    注释    \ny",
			expectedChar: 'y',
		},
	}

	for i, tt := range tests {
		l := New(tt.input)
		l.skipComment()
		if l.ch != tt.expectedChar {
			t.Errorf("Test %d failed. Expected ch='%c'; But got ch='%c'",
				i, tt.expectedChar, l.ch)
		}
	}
}

// TestReadChar 测试 readChar 函数
func TestReadChar(t *testing.T) {
	// 测试用例 1: 空输入
	l := New("")
	if l.ch != 0 { // EOF
		t.Errorf("Expected ch to be EOF (0), got %v", l.ch)
	}

	// 测试用例 2: 单个字符输入
	l = New("a")
	expected := byte('a')
	if l.ch != expected {
		t.Errorf("Expected ch to be '%c', got '%c'", expected, l.ch)
	}

	// 测试用例 3: 多字符输入，逐步读取
	input := "abc"
	l = New(input)
	// 第一次读取应为 'a'
	if l.ch != 'a' {
		t.Errorf("Expected first char to be 'a', got '%c'", l.ch)
	}
	if l.position != 0 || l.nextPosition != 1 {
		t.Errorf("Expected position=0 and nextPosition=1 after reading first char")
	}
	// 第二次读取应为 'b'
	l.readChar()
	if l.ch != 'b' {
		t.Errorf("Expected second char to be 'b', got '%c'", l.ch)
	}
	if l.position != 1 || l.nextPosition != 2 {
		t.Errorf("Expected position=1 and nextPosition=2 after reading second char")
	}
	// 第三次读取应为 'c'
	l.readChar()
	if l.ch != 'c' {
		t.Errorf("Expected third char to be 'c', got '%c'", l.ch)
	}
	if l.position != 2 || l.nextPosition != 3 {
		t.Errorf("Expected position=2 and nextPosition=3 after reading third char")
	}
	// 第四次读取应为 EOF
	l.readChar()
	if l.ch != 0 {
		t.Errorf("Expected EOF (0) after input ends, got %v", l.ch)
	}
	if l.position != 3 || l.nextPosition != 4 {
		t.Errorf("Expected position=3 and nextPosition=4 after reaching EOF")
	}
}

// TestPeekChar 测试 peekChar 函数
func TestPeekChar(t *testing.T) {
	tests := []struct {
		input           string
		initialPosition int
		expectedPeek    byte
	}{
		{
			// 空输入，nextPosition 超出范围
			input:           "",
			initialPosition: 0,
			expectedPeek:    END,
		},
		{
			// 当前位置是 'a'，下一个是 'b'
			input:           "ab",
			initialPosition: 0,
			expectedPeek:    'b',
		},
		{
			// 当前位置是 'b'，下一个是 EOF
			input:           "b",
			initialPosition: 1,
			expectedPeek:    END,
		},
		{
			// 多字符输入，在中间位置 peek
			input:           "abcdef",
			initialPosition: 3,
			expectedPeek:    'e',
		},
		{
			// 最后一个字符前 peek EOF
			input:           "abcde",
			initialPosition: 4,
			expectedPeek:    END,
		},
	}

	for i, tt := range tests {
		l := &Lexer{
			input:        tt.input,
			position:     tt.initialPosition,
			nextPosition: tt.initialPosition + 1,
			ch:           0,
		}
		peeked := l.peekChar()
		if peeked != tt.expectedPeek {
			t.Errorf("Test case %d failed. Expected peek='%c', got='%c'", i, tt.expectedPeek, peeked)
		}
	}
}

// TestNewToken 测试 newToken 函数
func TestNewToken(t *testing.T) {
	tests := []struct {
		tokenType token.TokenType
		ch        byte
		expected  token.Token
	}{
		// EOF
		{token.EOF, END, token.Token{Type: token.EOF, Literal: ""}},
		{token.IDENTIFIER, 'f', token.Token{Type: token.IDENTIFIER, Literal: "f"}},

		// 数据类型
		{token.INT, '5', token.Token{Type: token.INT, Literal: "5"}},
		{token.STRING, '"', token.Token{Type: token.STRING, Literal: "\""}},

		// 数学运算符
		{token.ASSIGN, '=', token.Token{Type: token.ASSIGN, Literal: "="}},
		{token.PLUS, '+', token.Token{Type: token.PLUS, Literal: "+"}},
		{token.MINUS, '-', token.Token{Type: token.MINUS, Literal: "-"}},
		{token.MUL, '*', token.Token{Type: token.MUL, Literal: "*"}},
		{token.DIV, '/', token.Token{Type: token.DIV, Literal: "/"}},
		{token.MOD, '%', token.Token{Type: token.MOD, Literal: "%"}},

		// 逻辑运算符
		{token.BAND, '!', token.Token{Type: token.BAND, Literal: "!"}},
		// 非法输入，需要组合成 '&&'
		{token.ILLEGAL, '&', token.Token{Type: token.ILLEGAL, Literal: "&"}},
		// 非法输入，需要组合成 '||'
		{token.ILLEGAL, '|', token.Token{Type: token.ILLEGAL, Literal: "|"}},

		// 比较运算符
		{token.LT, '<', token.Token{Type: token.LT, Literal: "<"}},
		{token.GT, '>', token.Token{Type: token.GT, Literal: ">"}},

		// 分隔符
		{token.LPAREN, '(', token.Token{Type: token.LPAREN, Literal: "("}},
		{token.RPAREN, ')', token.Token{Type: token.RPAREN, Literal: ")"}},
		{token.LBRACE, '{', token.Token{Type: token.LBRACE, Literal: "{"}},
		{token.RBRACE, '}', token.Token{Type: token.RBRACE, Literal: "}"}},
		{token.LBRACKET, '[', token.Token{Type: token.LBRACKET, Literal: "["}},
		{token.RBRACKET, ']', token.Token{Type: token.RBRACKET, Literal: "]"}},
		{token.COMMA, ',', token.Token{Type: token.COMMA, Literal: ","}},
		{token.SEMICOLON, ';', token.Token{Type: token.SEMICOLON, Literal: ";"}},
		{token.COLON, ':', token.Token{Type: token.COLON, Literal: ":"}},
	}

	for i, tt := range tests {
		tok := newToken(tt.tokenType, tt.ch)
		if tt.ch == END {
			tok.Literal = ""
		}
		if tok != tt.expected {
			t.Errorf("Test case %d failed. Expected %+v, got %+v", i, tt.expected, tok)
		}
	}
}

// TestReadString 测试 readString 方法
func TestReadString(t *testing.T) {
	tests := []struct {
		input           string
		expectedLiteral string
		// 期望读取结束后当前字符（应为结束引号后的字符或 EOF）
		expectedChar byte
	}{
		{
			// 正常字符串
			input:           `"hello world"`,
			expectedLiteral: "hello world",
			// 读完字符串后应到达 EOF
			expectedChar: END,
		},
		{
			// 空字符串
			input:           `""`,
			expectedLiteral: "",
			// 读完空字符串后应到达 EOF
			expectedChar: END,
		},
		{
			// 包含特殊字符的字符串
			input:           `"abc!@#$%^&*()"`,
			expectedLiteral: "abc!@#$%^&*()",
			expectedChar:    END,
		},
		{
			// 字符串后跟其他字符
			input:           `"test"abc`,
			expectedLiteral: "test",
			// 读完字符串后当前字符应为 'a'
			expectedChar: 'a',
		},
		{
			// 未闭合的字符串
			input: `"unclosed`,
			// 应读取到 EOF 前的内容
			expectedLiteral: "unclosed",
			// 当前字符应为 EOF
			expectedChar: END,
		},
		{
			// 多个连续字符串
			input:           `"foo" "bar"`,
			expectedLiteral: "foo",
			// 第一次读完字符串后当前字符应为 ' '
			expectedChar: ' ',
		},
	}

	for i, tt := range tests {
		l := New(tt.input)
		str := l.readString()
		if str != tt.expectedLiteral {
			t.Errorf("Test case %d failed. Expected literal=%q, got=%q", i, tt.expectedLiteral, str)
		}
		l.readChar()
		if l.ch != tt.expectedChar {
			t.Errorf("Test case %d failed. Expected current char='%c', got='%c'", i, tt.expectedChar, l.ch)
		}
	}
}

// TestIsLetter 测试 isLetter 函数
func TestIsLetter(t *testing.T) {
	// 生成所有小写字母和大写字母的测试用例
	var tests []struct {
		ch       byte
		expected bool
	}

	// 添加 a-z 小写字母
	for ch := 'a'; ch <= 'z'; ch++ {
		tests = append(tests, struct {
			ch       byte
			expected bool
		}{byte(ch), true})
	}

	// 添加 A-Z 大写字母
	for ch := 'A'; ch <= 'Z'; ch++ {
		tests = append(tests, struct {
			ch       byte
			expected bool
		}{byte(ch), true})
	}

	// 手动添加其他特殊情况
	tests = append(tests,
		struct {
			ch       byte
			expected bool
		}{'_', true},
		struct {
			ch       byte
			expected bool
		}{'0', false},
		struct {
			ch       byte
			expected bool
		}{'9', false},
		struct {
			ch       byte
			expected bool
		}{'!', false},
		struct {
			ch       byte
			expected bool
		}{'@', false},
		struct {
			ch       byte
			expected bool
		}{'#', false},
		struct {
			ch       byte
			expected bool
		}{' ', false},
		struct {
			ch       byte
			expected bool
		}{'\t', false},
		struct {
			ch       byte
			expected bool
		}{'\n', false},
		struct {
			ch       byte
			expected bool
		}{0, false},
	)

	for i, tt := range tests {
		result := isLetter(tt.ch)
		if result != tt.expected {
			t.Errorf("Test case %d failed. Expected isLetter('%c')=%v, got %v", i, tt.ch, tt.expected, result)
		}
	}
}

// TestReadIdentifier 测试 readIdentifier 方法
func TestReadIdentifier(t *testing.T) {
	identifierType := string(token.IDENTIFIER)
	tests := []struct {
		input            string
		expectedLiteral  string
		expectedType     string
		expectedPosition int
	}{
		// 标识符测试用例
		{"abc", "abc", identifierType, 3},
		{"abc123", "abc123", identifierType, 6},
		{"a1b2c3", "a1b2c3", identifierType, 6},
		{"_test", "_test", identifierType, 5},
		{"test_", "test_", identifierType, 5},
		{"x", "x", identifierType, 1},
		{"xyz", "xyz", identifierType, 3},
		{"helloWorld", "helloWorld", identifierType, 10},
		{"test123abc", "test123abc", identifierType, 10},
		{"var1 var2", "var1", identifierType, 4},
		{"funcName()", "funcName", identifierType, 8},
		{"a)", "a", identifierType, 1},

		// 测试标识符跟着运算符
		{"a+", "a", identifierType, 1},
		{"a-", "a", identifierType, 1},
		{"a*", "a", identifierType, 1},
		{"a/", "a", identifierType, 1},
		{"a%", "a", identifierType, 1},

		// 测试标识符跟着逻辑运算符
		{"a!", "a", identifierType, 1},
		{"a&&", "a", identifierType, 1},
		{"a||", "a", identifierType, 1},

		// 测试标识符跟着比较运算符
		{"a<", "a", identifierType, 1},
		{"a>", "a", identifierType, 1},
		{"a<=", "a", identifierType, 1},
		{"a>=", "a", identifierType, 1},
		{"a==", "a", identifierType, 1},
		{"a!=", "a", identifierType, 1},

		// 关键字测试
		{"fn", "fn", string(token.FUNCTION), 2},
		{"let", "let", string(token.LET), 3},
		{"true", "true", string(token.TRUE), 4},
		{"false", "false", string(token.FALSE), 5},
		{"if", "if", string(token.IF), 2},
		{"else", "else", string(token.ELSE), 4},
		{"return", "return", string(token.RETURN), 6},

		// 特殊边界测试
		{"a=", "a", identifierType, 1},
		{"z,", "z", identifierType, 1},
		{"_", "_", identifierType, 1},
		{"__", "__", identifierType, 2},
		{"a_b_c", "a_b_c", identifierType, 5},

		// 非法输入测试（应不读取）
		{"1abc", "1abc", string(token.ILLEGAL), 4},
		{"!abc", "!abc", string(token.ILLEGAL), 4},
		{"", "", string(token.ILLEGAL), 0},
		{"   ", "", string(token.ILLEGAL), 0},
		{"123", "123", string(token.ILLEGAL), 3},
		{"a@", "a@", string(token.ILLEGAL), 2},
	}

	for i, tt := range tests {
		l := New(tt.input)
		identifier := l.readIdentifier()
		if identifier.Literal != tt.expectedLiteral || l.position != tt.expectedPosition || string(identifier.Type) != tt.expectedType {
			t.Errorf("Test case %d failed. Expected literal=%q, token type=%q and position=%d, got literal=%q, token type=%q and position=%d",
				i, tt.expectedLiteral, tt.expectedType, tt.expectedPosition, identifier.Literal, identifier.Type, l.position)
		}
	}
}

// TestIsDigit 测试 isDigit 函数
func TestIsDigit(t *testing.T) {
	tests := []struct {
		ch       byte
		expected bool
	}{
		{'0', true},
		{'1', true},
		{'2', true},
		{'3', true},
		{'4', true},
		{'5', true},
		{'6', true},
		{'7', true},
		{'8', true},
		{'9', true},
		{'a', false},
		{'A', false},
		{'_', false},
		{'!', false},
		{'@', false},
		{'#', false},
		{' ', false},
		{'\t', false},
		{'\n', false},
		{END, false},
	}

	for i, tt := range tests {
		result := isDigit(tt.ch)
		if result != tt.expected {
			t.Errorf("Test case %d failed. Expected isDigit('%c')=%v, got %v", i, tt.ch, tt.expected, result)
		}
	}
}

// TestReadNumber 测试 readNumber 方法
func TestReadNumber(t *testing.T) {
	tests := []struct {
		input            string
		expectedLiteral  string
		expectedType     token.TokenType
		expectedPosition int
	}{
		// 整数测试用例
		{"123", "123", token.INT, 3},
		{"456x", "456x", token.ILLEGAL, 4},
		{"0", "0", token.INT, 1},

		// 浮点数测试用例
		{"123.456", "123.456", token.FLOAT, 7},
		{"789.012x", "789.012x", token.ILLEGAL, 8},
		{"3.141592653589793", "3.141592653589793", token.FLOAT, 17},

		// 非法数字测试用例
		{"123.", "123.", token.ILLEGAL, 4},
		{"1abc", "1abc", token.ILLEGAL, 4},
		{"123.abc", "123.abc", token.ILLEGAL, 7},

		// 边界情况测试
		{"9999999999999999999999999999999999999999", "9999999999999999999999999999999999999999", token.INT, 40},
	}

	for i, tt := range tests {
		l := New(tt.input)
		tok := l.readNumber()
		if tok.Literal != tt.expectedLiteral || tok.Type != tt.expectedType || l.position != tt.expectedPosition {
			t.Errorf("Test case %d failed. Expected literal=%q, type=%v and position=%d, got literal=%q, type=%v and position=%d",
				i, tt.expectedLiteral, tt.expectedType, tt.expectedPosition, tok.Literal, tok.Type, l.position)
		}
	}
}

// TestIsEndSeparator 测试字符是否是结束分隔符
func TestIsEndSeparator(t *testing.T) {
	tests := []struct {
		input    byte
		expected bool
	}{
		{'(', false},
		{')', true},
		{'{', false},
		{'}', true},
		{'[', false},
		{']', true},
		{',', true},
		{';', true},
		{':', true},
		{' ', true},
		{'a', false},
		{'0', false},
		{END, true},
	}

	for i, tt := range tests {
		result := isEndSeparator(tt.input)
		if result != tt.expected {
			t.Errorf("Test case %d failed. Expected char=%c, but result is=%v",
				i, tt.input, result)
		}
	}
}

// TestIsSeparator 测试字符是否为分隔符
func TestIsSeparator(t *testing.T) {
	tests := []struct {
		input    byte
		expected bool
	}{
		{'(', true},
		{')', true},
		{'{', true},
		{'}', true},
		{'[', true},
		{']', true},
		{',', true},
		{';', true},
		{':', true},
		{' ', true},
		{'a', false},
		{'0', false},
		{END, true},
	}

	for i, tt := range tests {
		result := isSeparator(tt.input)
		if result != tt.expected {
			t.Errorf("Test case %d failed. Expected char=%c, but result is=%v",
				i, tt.input, result)
		}
	}
}

// TestIsOperator 测试字符是否为运算符
func TestIsOperator(t *testing.T) {
	tests := []struct {
		input    byte
		expected bool
	}{
		{'=', true},
		{'+', true},
		{'-', true},
		{'*', true},
		{'/', true},
		{'%', true},
		{'!', true},
		{'&', true},
		{'|', true},
		{'<', true},
		{'>', true},
		{' ', false},
		{'a', false},
		{'0', false},
		{END, false},
	}

	for i, tt := range tests {
		result := isOperator(tt.input)
		if result != tt.expected {
			t.Errorf("Test case %d failed. Expected char=%c, but result is=%v",
				i, tt.input, result)
		}
	}
}

// TestSkipNotEndSeparator 测试跳过非结束分隔符
func TestSkipNotEndSeparator(t *testing.T) {
	tests := []struct {
		input           string
		expectedLiteral string
		expectedChar    byte // 跳过后当前字符
	}{
		{"abc", "", END},
		{"x+123", "", END},
		{"123.", "", END},
		{"a+b*c", "", END},
		{"a,b", ",b", ','},
		{"test)", ")", ')'},
		{"hello:", ":", ':'},
		{"world;", ";", ';'},
		{" ", " ", ' '},
		{"", "", END},
	}

	for i, tt := range tests {
		l := New(tt.input)
		l.skipNotEndSeparator()
		readStr := l.input[l.position:]

		if readStr != tt.expectedLiteral {
			t.Errorf("Test case %d failed. Expected literal=%s, char=%c, got literal=%s, char=%c",
				i, tt.expectedLiteral, tt.expectedChar, readStr, l.ch)
		}
	}
}

// TestSkipNotEndSeparatorAndOperator 测试跳过非结束分隔符和非运算符
func TestSkipNotEndSeparatorAndOperator(t *testing.T) {
	tests := []struct {
		input           string
		expectedLiteral string
		expectedChar    byte // 跳过后当前字符
	}{
		{"abc", "", END},
		{"x123", "", END},
		{"123.", "", END},
		{"a,b", ",b", ','},
		{"test)", ")", ')'},
		{"hello:", ":", ':'},
		{"world;", ";", ';'},
		{" ", " ", ' '},
		{"", "", END},

		{"abc@", "", END},
		{"x123\\", "", END},
	}

	for i, tt := range tests {
		l := New(tt.input)
		l.skipNotEndSeparatorAndNotOperator()
		readStr := l.input[l.position:]

		if readStr != tt.expectedLiteral {
			t.Errorf("Test case %d failed. Expected literal=%s, char=%c, got literal=%s, char=%c",
				i, tt.expectedLiteral, tt.expectedChar, readStr, l.ch)
		}
	}
}
