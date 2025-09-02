// environment_test.go
package object

import (
	"testing"
)

// 测试 NewEnvironment 函数
func TestNewEnvironment(t *testing.T) {
	env := NewEnvironment()

	if env.store == nil {
		t.Error("NewEnvironment() did not initialize store")
	}
	if env.outer != nil {
		t.Error("NewEnvironment() should create environment with nil outer")
	}

	env.Set("test", &Integer{Value: 1})
	obj, ok := env.Get("test")
	if !ok {
		t.Error("Failed to get value from newly created environment")
	}

	integer, ok := obj.(*Integer)
	if !ok {
		t.Fatalf("Expected Integer, got %T", obj)
	}

	if integer.Value != 1 {
		t.Errorf("Expected value 1, got %d", integer.Value)
	}
}

// 测试 NewEnclosedEnvironment 函数
func TestNewEnclosedEnvironment(t *testing.T) {
	outer := NewEnvironment()
	outer.Set("x", &Integer{Value: 10})

	inner := NewEnclosedEnvironment(outer)

	// 检查内部环境是否正确创建
	if inner.store == nil {
		t.Error("NewEnclosedEnvironment() did not initialize store")
	}

	// 检查外部环境是否正确设置
	if inner.outer != outer {
		t.Error("NewEnclosedEnvironment() did not set outer environment correctly")
	}

	// 检查是否可以从外部环境获取值
	obj, ok := inner.Get("x")
	if !ok {
		t.Error("Failed to get value from outer environment")
	}

	integer, ok := obj.(*Integer)
	if !ok {
		t.Fatalf("Expected Integer, got %T", obj)
	}

	if integer.Value != 10 {
		t.Errorf("Expected value 10, got %d", integer.Value)
	}
}

// 测试 Get 方法
func TestEnvironmentGet(t *testing.T) {
	env := NewEnvironment()

	obj, ok := env.Get("nonexistent")
	if ok {
		t.Error("Get() should return false for nonexistent key")
	}
	if obj != nil {
		t.Error("Get() should return nil for nonexistent key")
	}

	expected := &Integer{Value: 42}
	env.Set("test", expected)

	obj, ok = env.Get("test")
	if !ok {
		t.Error("Get() should return true for existing key")
	}

	if obj != expected {
		t.Error("Get() should return the correct object")
	}

	str := &String{Value: "hello"}
	env.Set("str", str)

	obj, ok = env.Get("str")
	if !ok {
		t.Error("Get() should return true for existing string key")
	}

	if obj != str {
		t.Error("Get() should return the correct string object")
	}
}

// 测试 Set 方法
func TestEnvironmentSet(t *testing.T) {
	env := NewEnvironment()

	integer := &Integer{Value: 100}
	result := env.Set("x", integer)

	if result != integer {
		t.Error("Set() should return the value that was set")
	}

	obj, ok := env.Get("x")
	if !ok {
		t.Error("Get() should return true for key that was set")
	}

	retrieved, ok := obj.(*Integer)
	if !ok {
		t.Fatalf("Expected Integer, got %T", obj)
	}

	if retrieved.Value != 100 {
		t.Errorf("Expected value 100, got %d", retrieved.Value)
	}

	newInteger := &Integer{Value: 200}
	env.Set("x", newInteger)

	obj, ok = env.Get("x")
	if !ok {
		t.Error("Get() should return true for key that was set")
	}

	retrieved, ok = obj.(*Integer)
	if !ok {
		t.Fatalf("Expected Integer, got %T", obj)
	}

	if retrieved.Value != 200 {
		t.Errorf("Expected value 200, got %d", retrieved.Value)
	}

	str := &String{Value: "world"}
	env.Set("y", str)

	obj, ok = env.Get("y")
	if !ok {
		t.Error("Get() should return true for string key that was set")
	}

	retrievedStr, ok := obj.(*String)
	if !ok {
		t.Fatalf("Expected String, got %T", obj)
	}

	if retrievedStr.Value != "world" {
		t.Errorf("Expected value 'world', got '%s'", retrievedStr.Value)
	}
}
