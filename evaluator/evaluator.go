package evaluator

import (
	"holiya/object"
)

var (
	// null 实例，null是唯一的，所以可以先初始化
	NULL = &object.Null{}
	// true 实例，true也是唯一的，也先初始化
	TRUE = &object.Boolean{Value: true}
	// false 实例，false也是唯一的，也先初始化
	FALSE = &object.Boolean{Value: false}
)
