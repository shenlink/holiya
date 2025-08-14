package parser

import (
	"fmt"
	"holiya/ast"
	"holiya/lexer"
	"holiya/token"
	"strconv"
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

type (
	// 解析前缀表达式
	prefixParseFn func() ast.Expression
	// 解析中缀表达式
	infixParseFn func(ast.Expression) ast.Expression
)

// 解析器
type Parser struct {
	// Lexer指针
	l *lexer.Lexer
	// 错误信息
	errors []string

	// 当前指针指向的token
	currToken token.Token
	// 指针指向的下一个token
	peekToken token.Token

	// 前缀解析函数map
	prefixParseFns map[token.TokenType]prefixParseFn
	// 中缀解析函数map
	infixParseFns map[token.TokenType]infixParseFn
}

// New 实例化Parser
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	// 注册标识符的前缀表达式的解析函数
	p.registerPrefix(token.IDENTIFIER, p.parseIdentifier)
	// 注册整数的前缀表达式的解析函数
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	// 注册浮点数的前缀表达式的解析函数
	p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
	// 注册字符串前缀表达式的解析函数
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	// 注册!的前缀表达式的解析函数
	p.registerPrefix(token.BAND, p.parsePrefixExpression)
	// 注册-的前缀表达式的解析函数
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	// 注册true的前缀表达式的解析函数
	p.registerPrefix(token.TRUE, p.parseBoolean)
	// 注册false的前缀表达式的解析函数
	p.registerPrefix(token.FALSE, p.parseBoolean)
	// 注册(的前缀表达式的解析函数
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	// 注册if的前缀表达式的解析函数
	p.registerPrefix(token.IF, p.parseIfExpression)
	// 注册函数的前缀表达式的解析函数
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	// 注册[的前缀表达式的解析函数
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	// 注册{的前缀表达式的解析函数
	p.registerPrefix(token.LBRACE, p.parseHashLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	// 注册+的中缀表达式的解析函数
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	// 注册-的中缀表达式的解析函数
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	// 注册*的中缀表达式的解析函数
	p.registerInfix(token.MUL, p.parseInfixExpression)
	// 注册/的中缀表达式的解析函数
	p.registerInfix(token.DIV, p.parseInfixExpression)
	// 注册/的中缀表达式的解析函数
	p.registerInfix(token.MOD, p.parseInfixExpression)
	// 注册<的中缀表达式的解析函数
	p.registerInfix(token.LT, p.parseInfixExpression)
	// 注册>的中缀表达式的解析函数
	p.registerInfix(token.GT, p.parseInfixExpression)
	// 注册<=的中缀表达式的解析函数
	p.registerInfix(token.LTE, p.parseInfixExpression)
	// 注册>=的中缀表达式的解析函数
	p.registerInfix(token.GTE, p.parseInfixExpression)
	// 注册==的中缀表达式的解析函数
	p.registerInfix(token.EQ, p.parseInfixExpression)
	// 注册!=的中缀表达式的解析函数
	p.registerInfix(token.NEQ, p.parseInfixExpression)

	// 注册(的中缀表达式的解析函数
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	// 注册[的中缀表达式的解析函数
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)

	// 这里调用了2次nextToken，第一次currToken = nil，peekToken = 第一个token，
	// 第二次，currToken = 第一个token，peekToken = 第二个token
	p.nextToken()
	p.nextToken()

	return p
}

// 注册前缀表达式
// 注册identifier，int，float，string，!，-，true，false，(，if，fn，[，{
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// 解析标识符，直接返回标识符的表达式节点
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
}

// 解析整数
// 注意：解析整数时，使用int64存储整数值
func (p *Parser) parseIntegerLiteral() ast.Expression {
	value, err := strconv.ParseInt(p.currToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %s as integer", p.currToken.Literal)
		p.appendError(msg)
		return nil
	}

	literal := &ast.IntegerLiteral{Token: p.currToken}
	literal.Value = value

	return literal
}

// 解析浮点数
// 注意：解析浮点数时，使用float64存储浮点数数值
func (p *Parser) parseFloatLiteral() ast.Expression {
	value, err := strconv.ParseFloat(p.currToken.Literal, 64)
	if err != nil { 
		msg := fmt.Sprintf("could not parse %s as float", p.currToken.Literal)
		p.appendError(msg)
		return nil
	}

	literal := &ast.FloatLiteral{Token: p.currToken}
	literal.Value = value

	return literal
}

// 实现Expression接口的方法
func (p *Parser) parseStringLiteral() ast.Expression {
	return nil
}

// 解析前缀表达式
func (p *Parser) parsePrefixExpression() ast.Expression {
	return nil
}

// 解析bool值
func (p *Parser) parseBoolean() ast.Expression {
	return nil
}

// 解析()
func (p *Parser) parseGroupedExpression() ast.Expression {
	return nil
}

// 解析if
func (p *Parser) parseIfExpression() ast.Expression {
	return nil
}

// 解析函数
func (p *Parser) parseFunctionLiteral() ast.Expression {
	return nil
}

// 解析数组
func (p *Parser) parseArrayLiteral() ast.Expression {
	return nil
}

// 解析哈希（键值对）
func (p *Parser) parseHashLiteral() ast.Expression {
	return nil
}

// 注册中缀表达式
// 注册 =，+，-，*，/，%，&&，||，<，>，<=，>=，==，!=，(，[
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// 解析中缀表达式
// 参数是中缀运算符的左边的表达式
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	return nil
}

// 解析调用表达式
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	return nil
}

// 解析index表达式
func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	return nil
}

// 获取下一个token，currToken指向下一个token，peekToken指向下两个token
func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// 添加错误信息到错误列表
func (p *Parser) appendError(errorMessage string) {
	p.errors = append(p.errors, errorMessage)
}