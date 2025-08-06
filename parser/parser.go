package parser

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
