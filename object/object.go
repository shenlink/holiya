package object

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
