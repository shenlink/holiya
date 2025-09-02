// object_test.go
package object

import (
	"holiya/ast"
	"holiya/token"
	"testing"
)

// 测试 Integer 对象的 Type 方法
func TestIntegerType(t *testing.T) {
	integer := &Integer{Value: 5}
	if integer.Type() != INTEGER_OBJ {
		t.Errorf("Integer.Type() = %s, want %s", integer.Type(), INTEGER_OBJ)
	}
}

// 测试 Integer 对象的 Inspect 方法
func TestIntegerInspect(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
	}{
		{5, "5"},
		{0, "0"},
		{-5, "-5"},
		{9223372036854775807, "9223372036854775807"},
		{-9223372036854775808, "-9223372036854775808"},
	}

	for _, tt := range tests {
		integer := &Integer{Value: tt.input}
		if integer.Inspect() != tt.expected {
			t.Errorf("Integer.Inspect() = %s, want %s", integer.Inspect(), tt.expected)
		}
	}
}

// 测试 Integer 对象的 GetHashKey 方法
func TestIntegerGetHashKey(t *testing.T) {
	int1 := &Integer{Value: 5}
	int2 := &Integer{Value: 5}
	int3 := &Integer{Value: 10}

	key1 := int1.GetHashKey()
	key2 := int2.GetHashKey()
	key3 := int3.GetHashKey()

	if key1.Type != INTEGER_OBJ {
		t.Errorf("Integer.GetHashKey().Type = %s, want %s", key1.Type, INTEGER_OBJ)
	}

	if key1.Value != key2.Value {
		t.Errorf("Integer with same value have different hash keys")
	}

	if key1.Value == key3.Value {
		t.Errorf("Integer with different value have same hash keys")
	}
}

// 测试 Float 对象的 Type 方法
func TestFloatType(t *testing.T) {
	float := &Float{Value: 5.5}
	if float.Type() != FLOAT_OBJ {
		t.Errorf("Float.Type() = %s, want %s", float.Type(), FLOAT_OBJ)
	}
}

// 测试 Float 对象的 Inspect 方法
func TestFloatInspect(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{5.5, "5.500000"},
		{0.0, "0.000000"},
		{-5.5, "-5.500000"},
		{3.141592653589793, "3.141593"},
	}

	for _, tt := range tests {
		float := &Float{Value: tt.input}
		if float.Inspect() != tt.expected {
			t.Errorf("Float.Inspect() = %s, want %s", float.Inspect(), tt.expected)
		}
	}
}

// 测试 Float 对象的 GetHashKey 方法
func TestFloatGetHashKey(t *testing.T) {
	float1 := &Float{Value: 5.5}
	float2 := &Float{Value: 5.5}
	float3 := &Float{Value: 10.10}

	key1 := float1.GetHashKey()
	key2 := float2.GetHashKey()
	key3 := float3.GetHashKey()

	if key1.Type != FLOAT_OBJ {
		t.Errorf("Float.GetHashKey().Type = %s, want %s", key1.Type, FLOAT_OBJ)
	}

	if key1.Value != key2.Value {
		t.Errorf("Float with same value have different hash keys")
	}

	if key1.Value == key3.Value {
		t.Errorf("Float with different value have same hash keys")
	}
}

// 测试 Boolean 对象的 Type 方法
func TestBooleanType(t *testing.T) {
	trueBool := &Boolean{Value: true}
	falseBool := &Boolean{Value: false}

	if trueBool.Type() != BOOLEAN_OBJ {
		t.Errorf("Boolean.Type() = %s, want %s", trueBool.Type(), BOOLEAN_OBJ)
	}
	if falseBool.Type() != BOOLEAN_OBJ {
		t.Errorf("Boolean.Type() = %s, want %s", falseBool.Type(), BOOLEAN_OBJ)
	}
}

// 测试 Boolean 对象的 Inspect 方法
func TestBooleanInspect(t *testing.T) {
	trueBool := &Boolean{Value: true}
	falseBool := &Boolean{Value: false}

	if trueBool.Inspect() != "true" {
		t.Errorf("Boolean.Inspect() = %s, want true", trueBool.Inspect())
	}
	if falseBool.Inspect() != "false" {
		t.Errorf("Boolean.Inspect() = %s, want false", falseBool.Inspect())
	}
}

// 测试 Boolean 对象的 GetHashKey 方法
func TestBooleanGetHashKey(t *testing.T) {
	true1 := &Boolean{Value: true}
	true2 := &Boolean{Value: true}
	false1 := &Boolean{Value: false}
	false2 := &Boolean{Value: false}

	trueKey1 := true1.GetHashKey()
	trueKey2 := true2.GetHashKey()
	falseKey1 := false1.GetHashKey()
	falseKey2 := false2.GetHashKey()

	if trueKey1.Type != BOOLEAN_OBJ {
		t.Errorf("Boolean.GetHashKey().Type = %s, want %s", trueKey1.Type, BOOLEAN_OBJ)
	}

	if trueKey1.Value != trueKey2.Value {
		t.Errorf("Boolean true values have different hash keys")
	}

	if falseKey1.Value != falseKey2.Value {
		t.Errorf("Boolean false values have different hash keys")
	}

	if trueKey1.Value == falseKey1.Value {
		t.Errorf("Boolean true and false have same hash keys")
	}
}

// 测试 Null 对象的 Type 方法
func TestNullType(t *testing.T) {
	null := &Null{}
	if null.Type() != NULL_OBJ {
		t.Errorf("Null.Type() = %s, want %s", null.Type(), NULL_OBJ)
	}
}

// 测试 Null 对象的 Inspect 方法
func TestNullInspect(t *testing.T) {
	null := &Null{}
	if null.Inspect() != "null" {
		t.Errorf("Null.Inspect() = %s, want null", null.Inspect())
	}
}

// 测试 ReturnValue 对象的 Type 方法
func TestReturnValueType(t *testing.T) {
	value := &Integer{Value: 5}
	rv := &ReturnValue{Value: value}
	if rv.Type() != RETURN_VALUE_OBJ {
		t.Errorf("ReturnValue.Type() = %s, want %s", rv.Type(), RETURN_VALUE_OBJ)
	}
}

// 测试 ReturnValue 对象的 Inspect 方法
func TestReturnValueInspect(t *testing.T) {
	value := &Integer{Value: 5}
	rv := &ReturnValue{Value: value}
	if rv.Inspect() != "5" {
		t.Errorf("ReturnValue.Inspect() = %s, want 5", rv.Inspect())
	}
}

// 测试 Error 对象的 Type 方法
func TestErrorType(t *testing.T) {
	err := &Error{Message: "test error"}
	if err.Type() != ERROR_OBJ {
		t.Errorf("Error.Type() = %s, want %s", err.Type(), ERROR_OBJ)
	}
}

// 测试 Error 对象的 Inspect 方法
func TestErrorInspect(t *testing.T) {
	err := &Error{Message: "test error"}
	expected := "ERROR: test error"
	if err.Inspect() != expected {
		t.Errorf("Error.Inspect() = %s, want %s", err.Inspect(), expected)
	}
}

// 测试 Function 对象的 Type 方法
func TestFunctionType(t *testing.T) {
	fn := &Function{
		Parameters: []*ast.Identifier{},
		Body:       &ast.BlockStatement{},
		Env:        NewEnvironment(),
	}

	if fn.Type() != FUNCTION_OBJ {
		t.Errorf("Function.Type() = %s, want %s", fn.Type(), FUNCTION_OBJ)
	}
}

// 测试 Function 对象的 Inspect 方法
func TestFunctionInspect(t *testing.T) {
	fn := &Function{
		Parameters: []*ast.Identifier{
			{Value: "x"},
			{Value: "y"},
		},
		Body: &ast.BlockStatement{
			Statements: []ast.Statement{
				&ast.ExpressionStatement{
					Expression: &ast.IntegerLiteral{
						Token: token.Token{Type: token.INT, Literal: "5"},
						Value: 5,
					},
				},
			},
		},
		Env: NewEnvironment(),
	}

	result := fn.Inspect()
	if len(result) == 0 {
		t.Errorf("Function.Inspect() returned empty string")
	}

	expected := "fn(x, y) {\n5\n}\n"
	if fn.Inspect() != expected {
		t.Errorf("Function.Inspect() = %s, want %s", fn.Inspect(), expected)
	}
}

// 测试 String 对象的 Type 方法
func TestStringType(t *testing.T) {
	str := &String{Value: "hello world"}
	if str.Type() != STRING_OBJ {
		t.Errorf("String.Type() = %s, want %s", str.Type(), STRING_OBJ)
	}
}

// 测试 String 对象的 Inspect 方法
func TestStringInspect(t *testing.T) {
	str := &String{Value: "hello world"}
	if str.Inspect() != "hello world" {
		t.Errorf("String.Inspect() = %s, want hello world", str.Inspect())
	}
}

// 测试 String 对象的 GetHashKey 方法
func TestStringGetHashKey(t *testing.T) {
	str1 := &String{Value: "hello"}
	str2 := &String{Value: "hello"}
	str3 := &String{Value: "world"}

	key1 := str1.GetHashKey()
	key2 := str2.GetHashKey()
	key3 := str3.GetHashKey()

	if key1.Type != STRING_OBJ {
		t.Errorf("String.GetHashKey().Type = %s, want %s", key1.Type, STRING_OBJ)
	}

	if key1.Value != key2.Value {
		t.Errorf("String with same value have different hash keys")
	}

	if key1.Value == key3.Value {
		t.Errorf("String with different value have same hash keys")
	}
}

// 测试 Builtin 对象的 Type 方法
func TestBuiltinType(t *testing.T) {
	builtin := &Builtin{}
	if builtin.Type() != BUILTIN_OBJ {
		t.Errorf("Builtin.Type() = %s, want %s", builtin.Type(), BUILTIN_OBJ)
	}
}

// 测试 Builtin 对象的 Inspect 方法
func TestBuiltinInspect(t *testing.T) {
	builtin := &Builtin{}
	if builtin.Inspect() != "builtin function" {
		t.Errorf("Builtin.Inspect() = %s, want builtin function", builtin.Inspect())
	}
}

// 测试 Array 对象的 Type 方法
func TestArrayType(t *testing.T) {
	array := &Array{
		Elements: []Object{},
	}
	if array.Type() != ARRAY_OBJ {
		t.Errorf("Array.Type() = %s, want %s", array.Type(), ARRAY_OBJ)
	}
}

// 测试 Array 对象的 Inspect 方法
func TestArrayInspect(t *testing.T) {
	array := &Array{
		Elements: []Object{
			&Integer{Value: 1},
			&Integer{Value: 2},
			&Integer{Value: 3},
		},
	}
	expected := "[1, 2, 3]"
	if array.Inspect() != expected {
		t.Errorf("Array.Inspect() = %s, want %s", array.Inspect(), expected)
	}

	emptyArray := &Array{Elements: []Object{}}
	if emptyArray.Inspect() != "[]" {
		t.Errorf("Array.Inspect() = %s, want []", emptyArray.Inspect())
	}
}

// 测试 Hashmap 对象的 Type 方法
func TestHashmapType(t *testing.T) {
	hash := &Hashmap{
		Pairs: map[HashKey]HashPair{},
	}
	if hash.Type() != HASH_OBJ {
		t.Errorf("Hashmap.Type() = %s, want %s", hash.Type(), HASH_OBJ)
	}
}

// 测试 Hashmap 对象的 Inspect 方法
func TestHashmapInspect(t *testing.T) {
	s := String{Value: "name"}
	hash := &Hashmap{
		Pairs: map[HashKey]HashPair{
			s.GetHashKey(): {
				Key:   &String{Value: "name"},
				Value: &String{Value: "Jimmy"},
			},
		},
	}

	result := hash.Inspect()
	if len(result) == 0 {
		t.Errorf("Hashmap.Inspect() returned empty string")
	}
	if len(result) < 2 || result[0] != '{' || result[len(result)-1] != '}' {
		t.Errorf("Hashmap.Inspect() format is incorrect: %s", result)
	}

	emptyHash := &Hashmap{Pairs: map[HashKey]HashPair{}}
	if emptyHash.Inspect() != "{}" {
		t.Errorf("Hashmap.Inspect() = %s, want {}", emptyHash.Inspect())
	}
}
