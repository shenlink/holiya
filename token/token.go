package token

// token 类型
type TokenType string

const (
	// 非法的字符
	ILLEGAL TokenType = "ILLEGAL"
	// 结束符
	EOF TokenType = "EOF"
	// 标识符
	IDENTIFIER TokenType = "IDENTIFIER"

	// 数据类型
	// 整数
	INT TokenType = "INT"
	// 浮点数
	FLOAT TokenType = "FLOAT"
	// 字符串
	STRING TokenType = "STRING"

	// 数学运算符
	// 赋值运算符
	ASSIGN TokenType = "="
	// 加法运算符
	PLUS TokenType = "+"
	// 减法运算符
	MINUS TokenType = "-"
	// 乘法运算符
	MUL TokenType = "*"
	// 除法运算符
	DIV TokenType = "/"
	// 取余运算符
	MOD TokenType = "%"

	// 逻辑运算符
	// 取反运算符
	BAND TokenType = "!"
	// 与运算符
	AND TokenType = "&&"
	// 或运算符
	OR TokenType = "||"

	// 比较运算符
	// 小于
	LT TokenType = "<"
	// 大于
	GT TokenType = ">"
	// 小于等于
	LTE TokenType = "<="
	// 大于等于
	GTE TokenType = ">="
	// 等于
	EQ TokenType = "=="
	// 不等于
	NEQ TokenType = "!="

	// 分隔符
	LPAREN    TokenType = "("
	RPAREN    TokenType = ")"
	LBRACE    TokenType = "{"
	RBRACE    TokenType = "}"
	LBRACKET  TokenType = "["
	RBRACKET  TokenType = "]"
	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"
	COLON     TokenType = ":"

	// 关键字
	// 函数关键字，声明函数
	FUNCTION TokenType = "FUNCTION"
	// 变量关键字，声明变量
	LET TokenType = "LET"
	// bool 关键字，表示 true
	TRUE TokenType = "TRUE"
	// bool 关键字，表示 false
	FALSE TokenType = "FALSE"
	// 流程控制关键字，表示 if
	IF TokenType = "IF"
	// 流程控制关键字，表示 else
	ELSE TokenType = "ELSE"
	// 函数关键字，表示终止执行，如果后面跟着表达式，则表示返回对应表达式的值
	RETURN TokenType = "RETURN"
)

// Token 结构体
type Token struct {
	// token 类型
	Type TokenType
	// token 的字符串表示
	Literal string
}

// 关键字 map
var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

// LookupIdent 查看是不是关键字，不是的话返回标识符
// identifier: 输入的字符串
func LookupIdentifier(identifier string) TokenType {
	if tok, ok := keywords[identifier]; ok {
		return tok
	}
	return IDENTIFIER
}
