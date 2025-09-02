package object

// 存储环境中的变量，分为内部变量和外部变量
// 内部变量是函数内部的变量，外部变量是函数
// 外部的变量
type Environment struct {
	store map[string]Object
	outer *Environment
}

// 创建一个环境
func NewEnvironment() *Environment {
	store := make(map[string]Object)
	return &Environment{
		store: store,
		outer: nil,
	}
}

// 创建一个环境，关闭外部环境
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// 获取变量
// 获取变量是先获取自己作用域内的变量，然后才是获取外部的变量
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// 设置变量
// 设置变量时，只能设置自己作用域范围内的变量，确保不会影响外部环境
func (e *Environment) Set(name string, value Object) Object {
	e.store[name] = value
	return value
}
