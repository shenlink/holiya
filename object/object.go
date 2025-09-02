package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"holiya/ast"
	"strings"
)

const (
	// null
	NULL_OBJ = "NULL"
	// error
	ERROR_OBJ = "ERROR"

	// 整数
	INTEGER_OBJ = "INTEGER"
	// 浮点数
	FLOAT_OBJ = "FLOAT"
	// 布尔值
	BOOLEAN_OBJ = "BOOLEAN"
	// 字符串
	STRING_OBJ = "STRING"

	// 返回值
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	// 函数
	FUNCTION_OBJ = "FUNCTION"
	// 内置函数
	BUILTIN_OBJ = "BUILTIN"

	// 数组
	ARRAY_OBJ = "ARRAY"
	// 哈希表
	HASH_OBJ = "HASH"
)

// 对象类型
type ObjectType string

// 对象接口
type Object interface {
	Type() ObjectType
	Inspect() string
}

// 内置函数类型
type BuiltinFunction func(args ...Object) Object

// 哈希键结构体
type HashKey struct {
	// 类型
	Type ObjectType
	// 哈希键的值
	Value uint64
}

// 哈希接口，实现该接口的对象都可以作为哈希键
type Hashable interface {
	// 获取哈希键
	GetHashKey() HashKey
}

// 整数
type Integer struct {
	Value int64
}

// 返回整数类型
func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

// 返回整数的字符串表示
func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

// 获取整数的哈希键对象
func (i *Integer) GetHashKey() HashKey {
	return HashKey{
		Type:  i.Type(),
		Value: uint64(i.Value)}
}

// 浮点数
type Float struct {
	Value float64
}

// 返回浮点数类型
func (f *Float) Type() ObjectType {
	return FLOAT_OBJ
}

// 返回浮点数的字符串表示
func (f *Float) Inspect() string {
	return fmt.Sprintf("%f", f.Value)
}

// 返回浮点数的哈希键对象
func (f *Float) GetHashKey() HashKey {
	return HashKey{
		Type:  f.Type(),
		Value: uint64(f.Value)}
}

// 布尔
type Boolean struct {
	Value bool
}

// 返回布尔类型
func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

// 获取布尔的字符串表示
func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

// 返回布尔的哈希键对象
func (b *Boolean) GetHashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{
		Type:  b.Type(),
		Value: value,
	}
}

// null，使用空结构体，使用了
type Null struct{}

// 返回null类型
func (n *Null) Type() ObjectType {
	return NULL_OBJ
}

// 返回null的字符串表示
func (n *Null) Inspect() string {
	return "null"
}

// 返回值
type ReturnValue struct {
	Value Object
}

// 返回返回值类型
func (rv *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJ
}

// 返回返回值的字符串表示
func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

// 错误
type Error struct {
	Message string
}

// 返回错误类型
func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}

// 返回错误的字符串表示
func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}

// 函数
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	// 函数的运行环境
	Env *Environment
}

// 返回函数类型
func (f *Function) Type() ObjectType {
	return FUNCTION_OBJ
}

// 返回函数的字符串表示
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}\n")

	return out.String()
}

// 字符串
type String struct {
	Value string
}

// 返回字符串类型
func (s *String) Type() ObjectType {
	return STRING_OBJ
}

// 直接返回字符串
func (s *String) Inspect() string {
	return s.Value
}

// 返回字符串的哈希键对象
func (s *String) GetHashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

// 内置函数
type Builtin struct {
	Fn BuiltinFunction
}

// 返回内置函数类型
func (b *Builtin) Type() ObjectType {
	return BUILTIN_OBJ
}

// 返回内置函数的字符串表示
func (b *Builtin) Inspect() string {
	return "builtin function"
}

// 数组
type Array struct {
	Elements []Object
}

// 返回数组类型
func (a *Array) Type() ObjectType {
	return ARRAY_OBJ
}

// 返回数组的字符串表示
func (a *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, element := range a.Elements {
		elements = append(elements, element.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// 哈希键值对
type HashPair struct {
	Key   Object
	Value Object
}

// 哈希表
type Hashmap struct {
	Pairs map[HashKey]HashPair
}

// 返回哈希表类型
func (h *Hashmap) Type() ObjectType {
	return HASH_OBJ
}

// 返回哈希表的字符串表示
func (h *Hashmap) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
