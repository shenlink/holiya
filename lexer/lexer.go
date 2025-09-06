package lexer

import (
	"holiya/token"
)

// 结束时返回 0，表示输入的内容已经结束了
const END byte = 0

// 词法结构体
type Lexer struct {
	// 输入的字符串
	input string
	// 当前字符的位置，默认为 0
	position int
	// 下一个字符的位置，默认为 0
	nextPosition int
	// 当前字符
	ch byte
}

// New 创建一个新的 Lexer 实例。
// 该函数接收一个字符串输入，用作 Lexer 的输入数据。
// 函数内部会初始化 Lexer 结构体，并调用 readChar 方法读取输入的第一个字符。
// 返回值是 Lexer 类型的指针，用于进一步处理输入数据。
func New(input string) *Lexer {
	l := &Lexer{
		input: input,
	}
	l.readChar()
	return l
}

// NextToken 获取下一个词法单元
// 该方法用于解析输入字符流，生成一个接一个的词法单元（Token）
// 它通过读取当前字符并根据字符的类型返回相应的Token
// 返回值：
//
//	token.Token: 词法单元
func (l *Lexer) NextToken() token.Token {
	// token变量
	var tok token.Token
	// 跳过空白符
	l.skipWhitespace()
	// 跳过注释
	l.skipComment()
	// 查看当前字符
	switch l.ch {
	// 数学运算符
	// 当前字符是 "="，可能是 "=="
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '*':
		tok = newToken(token.MUL, l.ch)
	case '/':
		tok = newToken(token.DIV, l.ch)
	case '%':
		tok = newToken(token.MOD, l.ch)

	// 逻辑运算符
	// !是取非运算符，之后可能是 "="，构成 "!="
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NEQ, Literal: literal}
		} else {
			tok = newToken(token.BAND, l.ch)
		}
	// &是逻辑与运算符&&的一部分，之后只能是 "&"，构成 "&&"
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.AND, Literal: literal}
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	// |是逻辑或运算符||的一部分，之后只能是 "|"，构成 "||"
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.OR, Literal: literal}
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	// 比较运算符
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.LTE, Literal: literal}
		} else {
			tok = newToken(token.LT, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.GTE, Literal: literal}
		} else {
			tok = newToken(token.GT, l.ch)
		}

	// 分隔符
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case '"':
		tok = token.Token{Type: token.STRING, Literal: l.readString()}
	// 结束符
	case END:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok = l.readIdentifier()
			return tok
		} else if isDigit(l.ch) {
			tok = l.readNumber()
			return tok
		} else {
			// 非法此法单元
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

// skipWhitespace 跳过输入中的空白字符
// 该函数通过移动读取指针，跳过如空格、制表符、换行符和回车符等空白字符
func (l *Lexer) skipWhitespace() {
	// 如果当前字符是空白字符，就移动position和readPosition指针到下一位
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// skipComment 跳过源代码中的注释，支持中文注释
// 该函数识别并跳过单行注释（以'//'开头）
// 对于单行注释，它会跳到行末。
func (l *Lexer) skipComment() {
	// 跳过单行注释
	for l.ch == '/' && l.peekChar() == '/' {
		// 读取并跳过两个'/'
		l.readChar()
		l.readChar()
		// 跳过注释内容直到行末或文件结束
		for l.ch == '\r' || l.ch != '\n' {
			// 防止EOF时无限循环
			if l.ch == END {
				break
			}
			l.readChar()
		}
		// 跳过空白字符符
		l.skipWhitespace()
	}
}

// readChar 读取当前字符并为下一个字符做准备
// 如果当前读取位置已经到达或超过输入字符串的长度，则将当前字符标记为END
// 否则，将当前位置的字符赋值给当前字符变量ch
// 随后，将当前位置指针移动到下一个字符位置，准备下一次读取
func (l *Lexer) readChar() {
	if l.nextPosition >= len(l.input) {
		l.ch = END
	} else {
		l.ch = (l.input[l.nextPosition])
	}
	// 注意这里：每次读取字符后，将当前位置指针移动到下一个字符位置
	// 也就是说读取完一个字符之后，再次读取的话，会读取到下一个字符
	l.position = l.nextPosition
	l.nextPosition++
}

// peekChar 返回当前位置之后的一个字符，而不改变当前的位置。
// 这个函数用于在不移动 Lexer position 位置的情况下预览下一个字符。
// 如果下一个位置超出了输入字符串的长度，则返回END标志。
// 返回值:
//
// byte 类型的字符，表示当前位置之后的字符或是END。
func (l *Lexer) peekChar() byte {
	// 检查当前位置之后的位置是否超出了输入字符串的长度
	if l.nextPosition >= len(l.input) {
		// 如果超出，返回END标志
		return END
	} else {
		// 如果没有超出，将当前位置的字符返回
		return l.input[l.nextPosition]
	}
}

// newToken 创建并返回一个新的 token.Token 实例。
// 该函数接受一个token类型（tokenType）和一个字符（ch）作为参数，
// 并使用这些参数构建一个token.Token结构体。
//
// 参数:
//
//	tokenType - token的类型，例如标识符、关键字等。
//	ch - token的字面值，由单个字符组成。
//
// 返回值:
//
//	返回一个token.Token结构体，其中包含了传入的类型和字符。
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}

// readString 读取一个字符串字面量。
// 该方法从当前字符位置开始，读取直到遇到结束引号 '"' 或输入结束标志 END。
// 它主要用于解析字符串常量，跳过引号内的内容，并返回解析的字符串。
// 该方法不处理转义字符或其他复杂的字符串内部逻辑。
func (l *Lexer) readString() string {
	// 记录字符串开始的位置，跳过开头的引号。
	position := l.position + 1
	for {
		// 逐字符读取输入。
		l.readChar()
		// 遇到结束引号或输入结束标志时停止读取。
		if l.ch == '"' || l.ch == END {
			break
		}
	}
	// 返回读取到的字符串内容，不包括引号。
	// 注意：这里的[]是左闭右开的
	return l.input[position:l.position]
}

// isLetter 判断给定的字符是否为字母或下划线。
// 参数:
//
//	ch byte: 需要判断的字符。
//
// 返回值:
//
//	bool: 如果字符是字母或下划线，则返回true，否则返回false。
func isLetter(ch byte) bool {
	// 判断字符是否在'a'到'z'、'A'到'Z'的范围内，或者是否为下划线。
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// readIdentifier 读取一个标识符或关键字
// 该方法从当前的字符开始，读取一系列连续的字母或下划线，直到遇到非字母且非下划线的字符为止
// 返回读取到的标识符或关键字字符串
func (l *Lexer) readIdentifier() token.Token {
	// 记录当前的位置，作为标识符的起始点
	position := l.position
	if !isLetter(l.ch) {
		l.skipNotEndSeparator()
		return token.Token{Type: token.ILLEGAL, Literal: l.input[position:l.position]}
	}
	// 读取当前的字母或下划线
	l.readChar()
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	// 这里不是字母，下划线和数字，如果不是分隔符或运算符，则是非法标识符
	if !isSeparator(l.ch) && !isOperator(l.ch) {
		l.skipNotEndSeparatorAndNotOperator()
		return token.Token{
			Type:    token.ILLEGAL,
			Literal: l.input[position:l.position],
		}
	}
	// 返回从起始点到当前位置（不包括当前位置）的字符串作为标识符或关键字
	// 这里使用了字符串切片的语法来获取标识符或关键字
	return token.Token{
		Type:    token.LookupIdentifier(l.input[position:l.position]),
		Literal: l.input[position:l.position],
	}
}

// isDigit 检查一个字符是否为数字。
// 参数:
//
//	ch byte: 需要检查的字符。
//
// 返回值:
//
//	bool: 如果字符是数字（0-9），则返回 true，否则返回 false。
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// readNumber 读取一个数字，可能是整数、浮点数或非法的数字
// 该函数通过读取字符来构建数字的字符串表示，并根据内容决定数字的类型
func (l *Lexer) readNumber() token.Token {
	// 记录数字开始的位置
	position := l.position
	// 是数字就一直读取
	for isDigit(l.ch) {
		l.readChar()
	}
	if l.ch == '.' {
		l.readChar()
		if !isDigit(l.ch) {
			l.skipNotEndSeparator()
			return token.Token{
				Type:    token.ILLEGAL,
				Literal: l.input[position:l.position],
			}
		}
		for isDigit(l.ch) {
			l.readChar()
		}
		if !isEndSeparator(l.ch) {
			l.skipNotEndSeparator()
			return token.Token{
				Type:    token.ILLEGAL,
				Literal: l.input[position:l.position],
			}
		}
		return token.Token{
			Type:    token.FLOAT,
			Literal: l.input[position:l.position],
		}
	} else if !isEndSeparator(l.ch) {
		l.skipNotEndSeparator()
		return token.Token{
			Type:    token.ILLEGAL,
			Literal: l.input[position:l.position],
		}
	}
	return token.Token{
		Type:    token.INT,
		Literal: l.input[position:l.position],
	}
}

// isEndSeparator 是否是结束分隔符
// 参数：
//
//	ch byte: 需要判断的字符
func isEndSeparator(ch byte) bool {
	return ch == ')' || ch == '}' || ch == ']' || ch == ',' || ch == ';' || ch == ':' || ch == ' ' || ch == END || ch == '\n'
}

// isSeparator 是否是分隔符
// 参数：
//
//	ch byte: 需要判断的字符
func isSeparator(ch byte) bool {
	return isEndSeparator(ch) || ch == '(' || ch == '{' || ch == '['
}

// isOperator 是否是运算符
// 参数：
//
//	ch byte: 需要判断的字符
func isOperator(ch byte) bool {
	return ch == '=' || ch == '+' || ch == '-' || ch == '*' || ch == '/' || ch == '%' || ch == '!' || ch == '&' || ch == '|' || ch == '<' || ch == '>'
}

// skipNotEndSeparator 跳过非结束分隔符
// 参数：
//
//	ch byte: 需要判断的字符
func (l *Lexer) skipNotEndSeparator() {
	for !isEndSeparator(l.ch) {
		l.readChar()
	}
}

// skipNotEndSeparatorAndNotOperator 跳过非结束分隔符和非运算符
func (l *Lexer) skipNotEndSeparatorAndNotOperator() {
	for !isEndSeparator(l.ch) && !isOperator(l.ch) {
		l.readChar()
	}
}
