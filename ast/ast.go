package ast

import (
	"bytes"
	"holiya/token"
	"strings"
)

// Node 节点接口
// 语法节点
type Node interface {
	TokenLiteral() string
	String() string
}

// Expression 表达式接口
// 表达式有：标识符，整数，浮点数，字符串，中缀表达式，前缀表达式
// bool，if 表达式，函数调用，map，数组，索引表达式
// 表达式可以理解成值，可以复赋值给标识符，if 表达式除外
type Expression interface {
	Node
	expressionNode()
}

// Statement 语句接口
// 语句有let语句，return语句，块语句，表达式语句
// let语句类似let x = 5;
// return语句类似return x;
// 块语句类似{let x = 5; return x;}
// 表达式语句类似x，5
type Statement interface {
	Node
	statementNode()
}

// Program 程序结构体
// 一个 Program 是由多个 Statement 语句构成的
// Program 结构体实现 Node 接口，可以看成是词法树的根节点
type Program struct {
	Statements []Statement
}

// Identifier 标识符节点，实现了 Expression 接口的方法
// 标识符可以看成是变量，单个变量是表达式
type Identifier struct {
	Token token.Token
	Value string
}

// expressionNode 实现了 Expression 接口的方法
func (i *Identifier) expressionNode() {}

// TokenLiteral 实现了 Expression 接口的方法
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

// String 实现了 Expression 接口的方法
func (i *Identifier) String() string {
	return i.Value
}

// IntegerLiteral 结构体，实现了 Expression 接口的方法
// 表示整数字面量节点，如 5
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

// expressionNode 实现了 Expression 接口的方法
func (i *IntegerLiteral) expressionNode() {}

// TokenLiteral 实现了 Expression 接口的方法
func (i *IntegerLiteral) TokenLiteral() string {
	return i.Token.Literal
}

// String 实现了 Expression 接口的方法
func (i *IntegerLiteral) String() string {
	return i.Token.Literal
}

// FloatLiteral 结构体，实现了 Expression 接口的方法
// 表示浮点数字面量节点，如 3.14
type FloatLiteral struct {
	Token token.Token
	Value float64
}

// expressionNode 实现了 Expression 接口的方法
func (i *FloatLiteral) expressionNode() {}

// TokenLiteral 实现了 Expression 接口的方法
func (i *FloatLiteral) TokenLiteral() string {
	return i.Token.Literal
}

// String 实现了 Expression 接口的方法
func (i *FloatLiteral) String() string {
	return i.Token.Literal
}

// StringLiteral 字符串字面量节点，实现 Expression 接口
// 表示字符串字面量，如 "hello"
type StringLiteral struct {
	Token token.Token
	Value string
}

// expressionNode 实现了 Expression 接口的方法
func (sl *StringLiteral) expressionNode() {}

// TokenLiteral 实现了 Expression 接口的方法
func (sl *StringLiteral) TokenLiteral() string {
	return sl.Token.Literal
}

// String 实现了 Expression 接口的方法
func (sl *StringLiteral) String() string {
	return sl.Token.Literal
}

// InfixExpression 中缀表达式节点，如 a + b
// 包含左右操作数和操作符
type InfixExpression struct {
	// 操作符，比如+，-
	Token token.Token
	Left  Expression
	// 操作符
	Operator string
	Right    Expression
}

// expressionNode 实现了 Expression 接口的方法
func (p *InfixExpression) expressionNode() {}

// TokenLiteral 实现了 Expression 接口的方法
func (p *InfixExpression) TokenLiteral() string {
	return p.Token.Literal
}

// String 实现了 Expression 接口的方法
func (p *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(p.Left.String())
	out.WriteString(" " + p.Operator + " ")
	out.WriteString(p.Right.String())
	out.WriteString(")")

	return out.String()
}

// PrefixExpression 前缀表达式节点，如 -a
// 包含前缀操作符和右边的操作数
type PrefixExpression struct {
	// 前缀操作符，比如-，！
	Token token.Token
	// 前缀操作符
	Operator string
	Right    Expression
}

// expressionNode 实现了 Expression 接口的方法
func (p *PrefixExpression) expressionNode() {}

// TokenLiteral 实现了 Expression 接口的方法
func (p *PrefixExpression) TokenLiteral() string {
	return p.Token.Literal
}

// String 实现了 Expression 接口的方法
func (p *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(p.Operator)
	out.WriteString(p.Right.String())
	out.WriteString(")")

	return out.String()
}

// FunctionStatement 函数声明语句，用于声明函数
// 包含函数参数和函数体
type FunctionStatement struct {
	// fn 关键字
	Token      token.Token
	Name       *Identifier
	Parameters []*Identifier
	Body       *BlockStatement
}

// expressionNode 实现了 Expression 接口的方法
func (fl *FunctionStatement) statementNode() {}

// TokenLiteral 实现了 Expression 接口的方法
func (fl *FunctionStatement) TokenLiteral() string {
	return fl.Token.Literal
}

// String 实现了 Expression 接口的方法
func (fl *FunctionStatement) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(fl.Body.String())

	return out.String()
}

// Boolean 布尔字面量节点，表示 true 或 false
// 用于表示布尔类型的值
type Boolean struct {
	Token token.Token
	Value bool
}

// expressionNode 实现了 Expression 接口的方法
func (b *Boolean) expressionNode() {}

// TokenLiteral 实现了 Expression 接口的方法
func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

// String 实现了 Expression 接口的方法
func (b *Boolean) String() string {
	return b.Token.Literal
}

// IfExpression if 表达式节点，支持条件判断，如 if (x > 5) { ... } else { ... }
// 包含条件表达式、条件为真时执行的代码块和可选的条件为假时执行的代码块
type IfExpression struct {
	// if 关键字
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

// expressionNode 实现了 Expression 接口的方法
func (ie *IfExpression) expressionNode() {}

// TokenLiteral 实现了 Expression 接口的方法
func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}

// String 实现了 Expression 接口的方法
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

// CallExpression 函数调用表达式节点，如 myFunction(1, 2)
// 表示对函数的调用，包含函数名和参数列表
type CallExpression struct {
	// (
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

// expressionNode 实现了 Expression 接口的方法
func (ce *CallExpression) expressionNode() {}

// TokenLiteral 实现了 Expression 接口的方法
func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}

// String 实现了 Expression 接口的方法
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

// HashLiteral 哈希（键值对）字面量节点，如 {"a": "b"}
// 表示键值对集合，可用于创建字典或对象
type HashLiteral struct {
	// {
	Token token.Token
	Pairs map[Expression]Expression
}

// expressionNode 实现了 Expression 接口的方法
func (hl *HashLiteral) expressionNode() {}

// TokenLiteral 实现了 Expression 接口的方法
func (hl *HashLiteral) TokenLiteral() string {
	return hl.Token.Literal
}

// String 实现了 Expression 接口的方法
func (hl *HashLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

// ArrayLiteral 数组字面量节点，如 [1, 2, 3]
// 表示数组的字面量形式
type ArrayLiteral struct {
	// [
	Token    token.Token
	Elements []Expression
}

// expressionNode 实现了 Expression 接口的方法
func (al *ArrayLiteral) expressionNode() {}

// TokenLiteral 实现了 Expression 接口的方法
func (al *ArrayLiteral) TokenLiteral() string {
	return al.Token.Literal
}

// String 实现了 Expression 接口的方法
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	var elements []string
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// IndexExpression 索引表达式节点，如 array[index]
// 表示对数组，map 或字符串的索引访问，
type IndexExpression struct {
	// [
	Token token.Token
	Left  Expression
	Index Expression
}

// expressionNode 实现了 Expression 接口的方法
func (ie *IndexExpression) expressionNode() {}

// TokenLiteral 实现了 Expression 接口的方法
func (ie *IndexExpression) TokenLiteral() string {
	return ie.Token.Literal
}

// String 实现了 Expression 接口的方法
func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("]")
	out.WriteString(")")

	return out.String()
}

// LetStatement let语句节点，如 let x = 5;
// 用于声明并初始化变量
type LetStatement struct {
	// let
	Token token.Token
	Name  *Identifier
	Value Expression
}

// statementNode 实现 Statement 接口的方法
func (ls *LetStatement) statementNode() {}

// TokenLiteral 实现 Statement 接口的方法
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

// String 实现 Statement 接口的方法
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

// ReturnStatement return 语句节点，如 return x;
// 用于表示从函数中返回值
type ReturnStatement struct {
	// return
	Token       token.Token
	ReturnValue Expression
}

// statementNode 实现 Statement 接口的方法
func (rs *ReturnStatement) statementNode() {}

// TokenLiteral 实现 Statement 接口的方法
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

// String 实现 Statement 接口的方法
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")

	return out.String()
}

// ExpressionStatement 表达式语句节点，用于包装表达式作为语句使用
// 例如表达式 x + y; 作为语句存在
type ExpressionStatement struct {
	// 表达式语句的第一个 token
	Token      token.Token
	Expression Expression
}

// statementNode 实现 Statement 接口的方法
func (es *ExpressionStatement) statementNode() {}

// TokenLiteral 实现 Statement 接口的方法
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

// String 实现 Statement 接口的方法
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// BlockStatement 块语句节点，由多个语句组成，如 { let x = 5; return x; }
// 用于表示代码块中的多个语句
type BlockStatement struct {
	// {
	Token      token.Token
	Statements []Statement
}

// statementNode 实现 Statement 接口的方法
func (bs *BlockStatement) statementNode() {}

// TokenLiteral 实现 Statement 接口的方法
func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

// String 实现 Statement 接口的方法
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// TokenLiteral 实现了 Expression 接口的方法
// 返回程序第一个语句的字面量，如果不存在语句则返回空字符串
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// String 实现了 Expression 接口的方法
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}
