package evaluator

import (
	"fmt"
	"holiya/ast"
	"holiya/object"
	"math"
)

var (
	// null 实例，null是唯一的，所以可以先初始化
	NULL = &object.Null{}
	// true 实例，true也是唯一的，也先初始化
	TRUE = &object.Boolean{Value: true}
	// false 实例，false也是唯一的，也先初始化
	FALSE = &object.Boolean{Value: false}
)

// 递归函数，用于对 AST 节点进行求值
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		// 处理整个程序节点
		return evalProgram(node, env)
	case *ast.LetStatement:
		// 处理 let 语句，将变量绑定到环境中
		value := Eval(node.Value, env)
		if isError(value) {
			return value
		}
		env.Set(node.Name.Value, value)
	case *ast.BlockStatement:
		// 处理代码块语句
		return evalBlockStatement(node, env)
	case *ast.ExpressionStatement:
		// 处理表达式语句
		return Eval(node.Expression, env)
	case *ast.ReturnStatement:
		// 处理 return 语句
		value := Eval(node.ReturnValue, env)
		if isError(value) {
			return value
		}
		return &object.ReturnValue{Value: value}
	case *ast.Identifier:
		// 处理标识符
		return evalIdentifier(node, env)
	case *ast.IntegerLiteral:
		// 处理整数字面量
		return &object.Integer{Value: node.Value}
	case *ast.FloatLiteral:
		// 处理浮点数字面量
		return &object.Float{Value: node.Value}
	case *ast.StringLiteral:
		// 处理字符串字面量
		return &object.String{Value: node.Value}
	case *ast.PrefixExpression:
		// 处理前缀表达式（如 -1, !true）
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.Boolean:
		// 处理布尔值字面量
		return nativeBoolToBooleanObject(node.Value)
	case *ast.IfExpression:
		// 处理 if 表达式
		return evalIfExpression(node, env)
	case *ast.FunctionLiteral:
		// 处理函数字面量
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Body: body, Env: env}
	case *ast.ArrayLiteral:
		// 处理数组字面量
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.HashLiteral:
		// 处理哈希表字面量
		return evalHashLiteral(node, env)
	case *ast.InfixExpression:
		// 处理中缀表达式（如 1 + 2, a == b）
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.CallExpression:
		// 处理函数调用表达式
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
	case *ast.IndexExpression:
		// 处理索引表达式（如 array[0], hash["key"], string[1]）
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	}
	// 对于未处理的节点类型，返回 nil
	return nil
}

// 用于计算整个程序的执行结果
// 它会按顺序执行程序中的每条语句，并返回最终结果
func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	// 初始化结果变量
	var result object.Object

	// 遍历程序中的所有语句
	for _, statement := range program.Statements {
		// 执行当前语句
		result = Eval(statement, env)

		// 类型断言检查执行结果
		switch result := result.(type) {
		case *object.ReturnValue:
			// 如果遇到返回值，直接返回其包装的值
			return result.Value
		case *object.Error:
			// 如果遇到错误，直接返回错误
			return result
		}
	}

	// 返回最后一条语句的执行结果
	return result
}

// 用于计算代码块中所有语句的执行结果
// 它会顺序执行代码块中的每条语句，直到遇到返回语句或错误
func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	// 初始化结果变量，用于存储执行结果
	var result object.Object

	// 遍历代码块中的所有语句
	for _, statement := range block.Statements {
		// 执行当前语句
		result = Eval(statement, env)
		// 检查执行结果是否为空
		if result != nil {
			// 获取结果类型
			resultType := result.Type()
			// 如果是返回值或错误，则立即返回，不执行后续语句
			if resultType == object.RETURN_VALUE_OBJ || resultType == object.ERROR_OBJ {
				return result
			}
		}
	}

	// 返回最后一个语句的执行结果
	return result
}

// 用于计算标识符节点的值
// 它首先在当前环境中查找标识符，如果找不到，则检查是否为内置函数
// 如果都找不到，则返回一个错误
func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	// 首先尝试在当前环境中查找标识符
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	// 如果在环境中找不到，则检查是否为内置函数
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	// 如果标识符在任何地方都找不到，则返回错误
	return newError("identifier not found: " + node.Value)
}

// 用于计算前缀表达式的值
func evalPrefixExpression(operator string, right object.Object) object.Object {
	// 根据运算符类型进行分支处理
	switch operator {
	case "!":
		// 处理逻辑非运算符
		return evalBandOperatorExpression(right)
	case "-":
		// 处理负号运算符
		return evalMinusPrefixOperatorExpression(right)
	default:
		// 未知运算符，返回错误信息
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

// 计算逻辑非运算符(!)的结果
// 该函数实现以下逻辑：
// - 当右操作数为 TRUE 时，返回 FALSE
// - 当右操作数为 FALSE 时，返回 TRUE
// - 当右操作数为 NULL 时，返回 TRUE
// - 其他所有情况返回 FALSE
func evalBandOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

// 处理 - 前缀操作符表达式
func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	// 检查操作数类型是否为整数或浮点数
	if right.Type() != object.INTEGER_OBJ && right.Type() != object.FLOAT_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	// 根据操作数类型进行相应的负号操作
	if right.Type() == object.INTEGER_OBJ {
		value := right.(*object.Integer).Value
		return &object.Integer{Value: -value}
	}
	value := right.(*object.Float).Value
	return &object.Float{Value: -value}
}

// 用于计算 if 表达式的值
// 该函数接收一个 if 表达式节点和当前环境，根据条件判断执行相应的分支
func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	// 计算 if 条件表达式的值
	condition := Eval(ie.Condition, env)

	// 如果条件计算出错，直接返回错误
	if isError(condition) {
		return condition
	}

	// 如果条件为真，执行 consequence 分支
	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		// 如果条件为假但存在 else 分支，则执行 else 分支
		return Eval(ie.Alternative, env)
	} else {
		// 如果条件为假且不存在 else 分支，则返回 NULL
		return NULL
	}
}

// 判断给定的对象是否为真值
func isTruthy(obj object.Object) bool {
	// 根据对象的具体值进行真值判断
	switch obj {
	case TRUE:
		return true
	case FALSE:
		return false
	case NULL:
		// NULL 对象视为 false
		return false
	default:
		// 其他所有对象都视为 true
		return true
	}
}

// 对表达式切片进行求值，返回对应的对象切片
func evalExpressions(expressions []ast.Expression, env *object.Environment) []object.Object {
	var results []object.Object

	// 遍历所有表达式并求值
	for _, expression := range expressions {
		value := Eval(expression, env)
		// 如果求值过程中出现错误，立即返回包含错误的切片
		if isError(value) {
			return []object.Object{value}
		}
		results = append(results, value)
	}

	return results
}

// 计算哈希表表达式的值
func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	// 创建存储哈希键值对的映射
	pairs := make(map[object.HashKey]object.HashPair)

	// 遍历所有键值对
	for keyNode, valueNode := range node.Pairs {
		// 计算键的值
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		// 检查键是否可哈希
		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}

		// 计算值的值
		value := Eval(valueNode, env)
		if isError(value) {
			return newError("error in hash literal: %s", value)
		}

		// 将键值对存储到哈希表中
		hashIndex := hashKey.GetHashKey()
		pairs[hashIndex] = object.HashPair{Key: key, Value: value}
	}

	// 返回构造好的哈希表对象
	return &object.Hashmap{Pairs: pairs}
}

// 计算中缀表达式的值（如 a + b, c == d 等）
func evalInfixExpression(operator string, left, right object.Object) object.Object {
	// 使用 switch 语句根据不同的情况处理表达式
	switch {
	// 当左右操作数都是整数时，调用整数专用处理函数
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)

	// 当左右操作数都是浮点数时，调用浮点数专用处理函数
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(operator, left, right)

	// 当左操作数都是整数，右操作数是浮点数时，调用数字专用处理函数
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalNumberInfixExpression(operator, left, right)

	// 当左操作数都是浮点数，右操作数是整数时，调用数字专用处理函数
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalNumberInfixExpression(operator, right, left)

	// 当左右操作数都是字符串时，调用字符串专用处理函数
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)

	// 处理相等性比较操作符 "=="
	case operator == "==":
		// 直接比较两个对象是否是同一个对象实例
		return nativeBoolToBooleanObject(left == right)

	// 处理不等性比较操作符 "!="
	case operator == "!=":
		// 直接比较两个对象是否不是同一个对象实例
		return nativeBoolToBooleanObject(left != right)

	// 当左右操作数类型不匹配时，返回类型不匹配错误
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())

	// 其他未处理的情况，返回未知操作符错误
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

// 计算两个整数对象之间的中缀表达式
func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	// 提取左右操作数的值
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	// 根据运算符进行相应的运算
	switch operator {
	case "+":
		// 加法运算
		return &object.Integer{Value: leftValue + rightValue}
	case "-":
		// 减法运算
		return &object.Integer{Value: leftValue - rightValue}
	case "*":
		// 乘法运算
		return &object.Integer{Value: leftValue * rightValue}
	case "/":
		// 除法运算，检查除零错误
		if rightValue == 0 {
			return newError("Division by zero")
		}
		return &object.Integer{Value: leftValue / rightValue}
	case "%":
		// 取模运算
		return &object.Integer{Value: leftValue % rightValue}
	case ">":
		// 大于比较运算
		return nativeBoolToBooleanObject(leftValue > rightValue)
	case ">=":
		// 大于等于比较运算
		return nativeBoolToBooleanObject(leftValue >= rightValue)
	case "<":
		// 小于比较运算
		return nativeBoolToBooleanObject(leftValue < rightValue)
	case "<=":
		// 小于等于比较运算
		return nativeBoolToBooleanObject(leftValue <= rightValue)
	case "==":
		// 等于比较运算
		return nativeBoolToBooleanObject(leftValue == rightValue)
	case "!=":
		// 不等于比较运算
		return nativeBoolToBooleanObject(leftValue != rightValue)
	default:
		// 未知运算符错误
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

// 计算两个浮点数对象之间的中缀表达式
func evalFloatInfixExpression(operator string, left, right object.Object) object.Object {
	// 提取左右操作数的浮点数值
	leftValue := left.(*object.Float).Value
	rightValue := right.(*object.Float).Value

	// 根据运算符进行相应的运算
	switch operator {
	case "+":
		// 加法运算
		return &object.Float{Value: leftValue + rightValue}
	case "-":
		// 减法运算
		return &object.Float{Value: leftValue - rightValue}
	case "*":
		// 乘法运算
		return &object.Float{Value: leftValue * rightValue}
	case "/":
		// 除法运算，检查除零错误
		if rightValue == 0 {
			return newError("Division by zero")
		}
		return &object.Float{Value: leftValue / rightValue}
	case "%":
		// 取模运算（使用 math.Mod 处理浮点数取模）
		return &object.Float{Value: math.Mod(leftValue, rightValue)}
	case ">":
		// 大于比较运算
		return nativeBoolToBooleanObject(leftValue > rightValue)
	case ">=":
		// 大于等于比较运算
		return nativeBoolToBooleanObject(leftValue >= rightValue)
	case "<":
		// 小于比较运算
		return nativeBoolToBooleanObject(leftValue < rightValue)
	case "<=":
		// 小于等于比较运算
		return nativeBoolToBooleanObject(leftValue <= rightValue)
	case "==":
		// 等于比较运算
		return nativeBoolToBooleanObject(leftValue == rightValue)
	case "!=":
		// 不等于比较运算
		return nativeBoolToBooleanObject(leftValue != rightValue)
	default:
		// 未知运算符错误
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

// 处理整数对象和浮点数对象之间的中缀表达式
func evalNumberInfixExpression(operator string, left, right object.Object) object.Object {
	leftValue := &object.Float{Value: float64(left.(*object.Integer).Value)}

	return evalFloatInfixExpression(operator, leftValue, right)
}

// 将原生布尔值转换为对象系统的布尔对象
func nativeBoolToBooleanObject(input bool) object.Object {
	// 根据输入的布尔值返回对应的对象系统布尔对象
	if input {
		return TRUE
	}
	return FALSE
}

// 处理字符串类型的中缀表达式运算
func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	// 检查是否为支持的运算符，目前只支持字符串连接运算符"+"
	if operator != "+" {
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}

	// 获取左右操作数的字符串值
	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value

	// 执行字符串连接运算并返回结果
	return &object.String{Value: leftValue + rightValue}
}

// 将函数对象应用于给定的参数并返回结果
func applyFunction(fn object.Object, args []object.Object) object.Object {
	// 使用类型断言来处理不同类型的函数
	switch fn := fn.(type) {
	case *object.Function:
		// 处理用户定义的函数
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		// 处理内置函数，直接调用其Fn字段
		return fn.Fn(args...)
	default:
		// 如果对象不是函数，则返回错误
		return newError("not a function: %s", fn.Type())
	}
}

// 创建函数执行环境，将函数参数绑定到封闭环境中的变量
func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}
	return env
}

// 从返回值对象中提取实际的值
// 如果传入的对象是 ReturnValue 类型，则返回其内部包装的值
func unwrapReturnValue(obj object.Object) object.Object {
	// 检查对象是否为 ReturnValue 类型，如果是则返回其内部值
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	// 如果不是 ReturnValue 类型，直接返回原对象
	return obj
}

// 计算索引表达式的值（如 array[0], hash["key"], string[1]）
// 根据左操作数和索引的类型，调用相应的处理函数
func evalIndexExpression(left, index object.Object) object.Object {
	// 使用 switch 语句根据不同的情况处理索引表达式
	switch {
	// 当左操作数是数组且索引是整数时，调用数组索引处理函数
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	// 当左操作数是哈希表时，调用哈希表索引处理函数
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	// 当左操作数是字符串且索引是整数时，调用字符串索引处理函数
	case left.Type() == object.STRING_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalStringIndexExpression(left, index)
	// 其他情况返回错误信息
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

// 计算数组索引表达式的值
// 该函数接收一个数组对象和一个索引对象，返回数组中对应索引位置的元素
// 如果索引超出数组范围，则返回 NULL
func evalArrayIndexExpression(array, index object.Object) object.Object {
	// 类型断言获取数组对象
	arrayObject := array.(*object.Array)

	// 类型断言获取索引值
	idx := index.(*object.Integer).Value

	// 检查索引是否越界（小于0或大于等于数组长度）
	if idx < 0 || idx > int64(len(arrayObject.Elements)-1) {
		return NULL
	}

	// 返回数组中指定索引位置的元素
	return arrayObject.Elements[idx]
}

// 计算哈希表索引表达式的值
// 该函数接收一个哈希表对象和一个索引对象，返回哈希表中对应键的值
// 如果键不可哈希，则返回错误信息，键不存在则返回 NULL
func evalHashIndexExpression(hash, index object.Object) object.Object {
	// 类型断言获取哈希表对象
	hashObject := hash.(*object.Hashmap)

	// 检查索引是否可哈希（只有可哈希的对象才能作为哈希表的键）
	key, ok := index.(object.Hashable)
	if !ok {
		// 如果索引不可哈希，返回错误信息
		return newError("unusable as hash key: %s", index.Type())
	}

	// 在哈希表中查找对应的键值对
	pair, ok := hashObject.Pairs[key.GetHashKey()]
	if !ok {
		// 如果键不存在，返回 NULL
		return NULL
	}

	// 返回找到的值
	return pair.Value
}

// 计算字符串索引表达式，获取字符串中指定位置的字符
func evalStringIndexExpression(stringObject, index object.Object) object.Object {
	stringValue := stringObject.(*object.String).Value
	idx := index.(*object.Integer).Value

	// 检查索引是否越界
	if idx < 0 || idx > int64(len(stringValue)-1) {
		return NULL
	}

	// 将字符串转换为 rune 切片以正确处理 Unicode 字符，然后返回指定索引的字符
	runes := []rune(stringValue)
	return &object.String{Value: string(runes[idx])}

}

// 检查给定的对象是否为错误类型
func isError(value object.Object) bool {
	return value != nil && value.Type() == object.ERROR_OBJ
}

// 创建一个新的 Error 对象，包含格式化的错误信息。
// 它接收一个格式字符串和可变参数，用于构造错误消息。
func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}
