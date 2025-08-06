package parser

import (
	"holiya/token"
)

const (
	// 占位
	_ int = iota
	// LOWEST 最低的优先级
	LOWEST
	// EQUALS =，!=
	EQUALS
	// SUM +，-
	SUM
	// LESSGREATER <，>
	LESSGREATER
	// PRODUCT *，/，%
	PRODUCT
	// PREFIX !，-
	PREFIX
	// CALL (，函数调用
	CALL
	// INDEX [，索引
	INDEX
)

// 定义所有的token类型对应的整数值
var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NEQ:      EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.DIV:      PRODUCT,
	token.MUL:      PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
}
