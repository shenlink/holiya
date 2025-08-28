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

// ParseProgram 解析Program
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.currTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
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

// 解析字符串表达式
func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.currToken, Value: p.currToken.Literal}
}

// 解析前缀表达式
func (p *Parser) parsePrefixExpression() ast.Expression {
	// 当前token是前缀运算符 '-' 和 '!'
	expression := &ast.PrefixExpression{Token: p.currToken, Operator: p.currToken.Literal}

	// 跳过 '!' 和 '-'
	p.nextToken()

	// 解析 '!' 和 '-' 后面的表达式
	expression.Right = p.parseExpression(LOWEST)

	return expression
}

// 解析表达式
func (p *Parser) parseExpression(precedence int) ast.Expression {
	// 获取前缀表达式的解析函数
	// 注意：这里的表达式对应的token不只是只有 '-' 和 '!'
	prefix := p.prefixParseFns[p.currToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.currToken.Type)
		return nil
	}

	// 解析前缀表达式
	leftExpression := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExpression
		}
		// 跳到中缀运算符
		p.nextToken()
		// 解析中缀表达式、
		// 这里为什么要赋值给leftExpression呢？
		// 有这种情况：5 + 3 - 2
		// + 和 - 是同优先级的运算符，当当前token是5时，会调用前缀表达式解析
		// 函数 parseIntegerLiteral 解析token 5，接着是判断
		// !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence()
		// 如果符合的，那获取中缀运算符解析函数，这里是有的，+ 运算符的解析函数是
		// parseInfixExpression，解析完成后之后，那就获取解析结果，并赋值给 leftExpression，
		// 然后是 - 运算符，还是符合 !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence()，
		// 所以，还是会进入这一段循环，解析出 x - 2 表达式
		leftExpression = infix(leftExpression)
	}
	return leftExpression
}

// 没有前缀表达式解析函数错误
func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// 获取当前token的优先级，在precedences中匹配不到的话就是最低优先级
func (p *Parser) currPrecedence() int {
	if p, ok := precedences[p.currToken.Type]; ok {
		return p
	}
	return LOWEST
}

// 辅助函数，判断当前token的类型与给出的token的类型是否相等
func (p *Parser) currTokenIs(t token.TokenType) bool {
	return p.currToken.Type == t
}

// 获取下一个token的优先级，在precedences中匹配不到的话就是最低优先级
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// 辅助函数，判断下一个token的类型与给出的token的类型是否相等
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return t == p.peekToken.Type
}

// 解析bool值
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.currToken, Value: p.currTokenIs(token.TRUE)}
}

// 解析()
func (p *Parser) parseGroupedExpression() ast.Expression {
	// 跳过(
	p.nextToken()

	// 解析括号中的表达式
	expression := p.parseExpression(LOWEST)

	// 如果不是以)结尾，则报错
	if !p.peekTokenIs(token.RPAREN) {
		p.errors = append(p.errors, "expected )")
		return nil
	}

	return expression
}

// 解析if
func (p *Parser) parseIfExpression() ast.Expression {
	// 当前token是if
	expression := &ast.IfExpression{Token: p.currToken}

	// 验证下一个token是不是(
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	// 解析if()中的条件表达式
	expression.Condition = p.parseExpression(LOWEST)
	// 验证下一个token是不是)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	// 验证下一个token是不是{
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	// 解析块语句
	expression.Consequence = p.parseBlockStatement()
	// 验证当前token是不是}
	if !p.currTokenIs(token.RBRACE) {
		p.appendError(fmt.Sprintf("expected }, but got %s", p.currToken.Type))
		return nil
	}

	// 验证下一个token是不是else
	if p.peekTokenIs(token.ELSE) {
		// 跳过}，跳过之后，当前token就是else了
		p.nextToken()
		// 验证下一个token是不是{
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		// 解析块语句
		expression.Alternative = p.parseBlockStatement()
		// 验证当前token是不是}
		if !p.currTokenIs(token.RBRACE) {
			p.appendError(fmt.Sprintf("expected }, but got %s", p.currToken.Type))
			return nil
		}
	}

	return expression
}

// 解析块语句，块语句就是在if-else和函数体中的语句
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	// 创建块语句
	block := &ast.BlockStatement{Token: p.currToken}

	// 初始化块语句的语句列表，go语言中，切片不初始化，直接append会panic
	block.Statements = []ast.Statement{}

	// 跳过{
	p.nextToken()

	// 如果当前token不是}或是EOF，则继续解析语句
	for !p.currTokenIs(token.RBRACE) && !p.currTokenIs(token.EOF) {
		// 解析语句，在holiya中，语句只有let语句，return语句和表达式语句
		statement := p.parseStatement()
		// statement可能会返回空
		if statement != nil {
			// 添加语句
			block.Statements = append(block.Statements, statement)
		}
		// 跳过每个语句的结尾，来到下一个语句的开头
		p.nextToken()
	}
	return block
}

// 解析语句
func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.LET:
		// 解析let语句
		return p.parseLetStatement()
	case token.RETURN:
		// 解析return语句
		return p.parseReturnStatement()
	default:
		// 解析表达式语句
		return p.parseExpressionStatement()
	}
}

// 解析let语句
func (p *Parser) parseLetStatement() *ast.LetStatement {
	// 当前token是let
	statement := &ast.LetStatement{Token: p.currToken}

	// 如果下一个token不是标识符，则记录错误并返回nil
	if !p.expectPeek(token.IDENTIFIER) {
		return nil
	}

	// 创建标识符
	statement.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}

	// 如果下一个token不是=，则记录错误并返回nil
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// 跳过=
	p.nextToken()

	// 解析表达式，并保存到statement中
	statement.Value = p.parseExpression(LOWEST)

	// 如果下一个token不是;，则记录错误并返回
	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}

	return statement
}

// 解析return语句
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	// 创建return语句
	statement := &ast.ReturnStatement{Token: p.currToken}

	// 跳过return关键字
	p.nextToken()

	// 解析表达式
	statement.ReturnValue = p.parseExpression(LOWEST)

	// 如果下一个token不是;，则记录错误并返回nil
	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}

	return statement
}

// 表达式语句，类似foobar;
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	// 创建表达式语句
	statement := &ast.ExpressionStatement{Token: p.currToken}

	// 解析表达式
	statement.Expression = p.parseExpression(LOWEST)

	// 如果下一个token不是;，则记录错误并返回nil
	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}

	return statement
}

// 解析函数
func (p *Parser) parseFunctionLiteral() ast.Expression {
	// 创建函数表达式
	fnExpression := &ast.FunctionLiteral{Token: p.currToken}

	// 如果下一个token不是(，则记录错误并返回nil
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	// 解析函数的参数列表
	parameters := p.parseFunctionParameters()

	// 如果参数列表为nil，则返回nil
	if parameters == nil {
		return nil
	}

	fnExpression.Parameters = parameters

	// 如果下一个token不是{，则记录错误并返回nil
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	// 解析函数的函数体
	fnExpression.Body = p.parseBlockStatement()

	return fnExpression
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	// 创建一个空数组，用于存储参数列表
	identifiers := []*ast.Identifier{}

	// 如果没有参数，则返回空数组
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	// 跳过(
	p.nextToken()

	// 第一个参数
	identifier := &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
	identifiers = append(identifiers, identifier)

	// 如果下一个token是,，则继续解析下一个参数
	for p.peekTokenIs(token.COMMA) {
		// 跳过当前参数所在的token
		p.nextToken()
		// 跳过,
		p.nextToken()
		// 添加参数
		identifier := &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
		identifiers = append(identifiers, identifier)
	}

	// 如果下一个token不是)，则记录错误并返回nil
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

// 解析数组
func (p *Parser) parseArrayLiteral() ast.Expression {
	// 创建数组表达式
	arrayExpression := &ast.ArrayLiteral{Token: p.currToken}

	// 解析表达式列表，以参数end作为标识
	list := p.parseExpressionList(token.RBRACKET)
	if list == nil {
		return nil
	}
	arrayExpression.Elements = list

	return arrayExpression
}

// 解析表达式列表，有数组和调用表达式的实参列表
func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	// 创建一个空数组，用于存储表达式列表
	list := []ast.Expression{}

	// 下一个token是结束符，那就是没有[]或()中没有表达式，返回空表达式数组
	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	// 跳过[或(
	p.nextToken()

	// 第一个表达式
	expression := p.parseExpression(LOWEST)
	list = append(list, expression)

	// 如果下一个token是,，则继续解析下一个表达式
	for p.peekTokenIs(token.COMMA) {
		// 跳到,
		p.nextToken()
		// 调到下一个表达式所在的token
		p.nextToken()
		// 解析表达式
		expression := p.parseExpression(LOWEST)
		// 添加表达式
		list = append(list, expression)
	}

	// 如果下一个token不是end，则记录错误并返回nil
	if !p.expectPeek(end) {
		return nil
	}

	return list
}

// 解析哈希（键值对）
func (p *Parser) parseHashLiteral() ast.Expression {
	// 创建哈希表达式
	hashExpression := &ast.HashLiteral{Token: p.currToken}

	// 创建一个空map，用于存储键值对
	hashExpression.Pairs = make(map[ast.Expression]ast.Expression)

	// 如果下一个token不是}，继续解析
	for !p.peekTokenIs(token.RBRACE) {
		// 跳过第一个token { 或,
		p.nextToken()
		// 解析键
		key := p.parseExpression(LOWEST)
		// 如果下一个token不是:，则记录错误并返回nil
		if !p.expectPeek(token.COLON) {
			return nil
		}
		// 跳过:
		p.nextToken()
		// 解析值
		value := p.parseExpression(LOWEST)
		// 添加键值对
		hashExpression.Pairs[key] = value
		// 如果下一个token不是},，则记录错误并返回nil
		if !p.peekTokenIs(token.RBRACE) && !p.peekTokenIs(token.COMMA) {
			p.appendError("Expected comma or right brace after hash pair")
			return nil
		}
		// 跳过,
		if p.peekTokenIs(token.COMMA) {
			p.nextToken()
		}
	}

	return hashExpression
}

// 注册中缀表达式
// 注册 =，+，-，*，/，%，&&，||，<，>，<=，>=，==，!=，(，[
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// 解析中缀表达式
// 参数是中缀运算符的左边的表达式
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.currToken,
		Left:     left,
		Operator: p.currToken.Literal,
	}

	// 获取当前token的优先级
	precedence := p.currPrecedence()
	// 跳过中缀运算符
	p.nextToken()
	// 解析中缀运算符右边的表达式
	expression.Right = p.parseExpression(precedence)

	return expression
}

// 解析调用表达式
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	// 创建调用表达式
	callExpression := &ast.CallExpression{Token: p.currToken, Function: function}

	// 解析参数列表
	arguments := p.parseExpressionList(token.RPAREN)
	if arguments == nil {
		return nil
	}
	callExpression.Arguments = arguments

	return callExpression
}

// 解析index表达式
func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	// 创建索引表达式
	indexExpression := &ast.IndexExpression{Token: p.currToken, Left: left}

	// 跳过[
	p.nextToken()

	// 解析索引
	index := p.parseExpression(LOWEST)
	if index == nil {
		return nil
	}
	indexExpression.Index = index

	// 如果下一个token不是]，则记录错误并返回nil
	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return indexExpression
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

// 判断下一个token与给出的token的类型是否相等，相等的话，跳过当前token
// 如果不相等，则调价错误信息，并返回false
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

// append错误信息
// 注意：词法解析应该要尽量得搜集更多的错误
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.appendError(msg)
}
