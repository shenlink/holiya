package evaluator

import (
	"fmt"
	"holiya/object"
)

// builtins 定义了所有内置函数的映射表
var builtins = map[string]*object.Builtin{
	// len 函数：返回数组或字符串的长度
	"len": {
		Fn: func(args ...object.Object) object.Object {
			// 检查参数数量是否正确（必须是1个）
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			// 根据参数类型执行不同操作
			switch arg := args[0].(type) {
			case *object.Array:
				// 如果是数组，返回其元素个数
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.String:
				// 如果是字符串，返回其字符个数
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				// 不支持的类型报错
				return newError("argument to `len` not supported, got %s", arg.Type())
			}
		}},

	// puts 函数：打印所有参数并返回 NULL
	"puts": {
		Fn: func(args ...object.Object) object.Object {
			// 遍历所有参数并打印它们的字符串表示
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			// 返回 nil，避免删除多余的 null 字符串
			return nil
		},
	},

	// first 函数：返回数组的第一个元素
	"first": {
		Fn: func(args ...object.Object) object.Object {
			// 检查参数数量是否正确（必须是1个）
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			// 检查参数是否为数组类型
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
			}
			// 获取数组对象
			arr := args[0].(*object.Array)
			// 如果数组不为空，返回第一个元素
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}
			// 空数组返回 NULL
			return NULL
		},
	},

	// last 函数：返回数组的最后一个元素
	"last": {
		Fn: func(args ...object.Object) object.Object {
			// 检查参数数量是否正确（必须是1个）
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			// 检查参数是否为数组类型
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `last` must be ARRAY, got %s", args[0].Type())
			}
			// 获取数组对象
			arr := args[0].(*object.Array)
			// 如果数组不为空，返回最后一个元素
			if len(arr.Elements) > 0 {
				return arr.Elements[len(arr.Elements)-1]
			}
			// 空数组返回 NULL
			return NULL
		},
	},

	// rest 函数：返回除第一个元素外的其余元素组成的新数组
	"rest": {
		Fn: func(args ...object.Object) object.Object {
			// 检查参数数量是否正确（必须是1个）
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			// 检查参数是否为数组类型
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `rest` must be ARRAY, got %s", args[0].Type())
			}
			// 获取数组对象
			arr := args[0].(*object.Array)
			// 如果数组不为空，创建新数组包含除第一个元素外的所有元素
			if len(arr.Elements) > 0 {
				newElements := make([]object.Object, len(arr.Elements)-1)
				copy(newElements, arr.Elements[1:])
				return &object.Array{Elements: newElements}
			}
			// 空数组返回 NULL
			return NULL
		},
	},

	// push 函数：向数组末尾添加一个元素并返回新数组（注意：函数名可能应为 "push"）
	"push": {
		Fn: func(args ...object.Object) object.Object {
			// 检查参数数量是否正确（必须是2个）
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			// 检查第一个参数是否为数组类型
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `push` must be ARRAY, got %s", args[0].Type())
			}
			// 获取原数组对象
			arr := args[0].(*object.Array)
			// 计算原数组长度
			length := len(arr.Elements)
			// 创建新数组，容量比原数组多1
			newElements := make([]object.Object, length+1)
			// 复制原数组元素到新数组
			copy(newElements, arr.Elements)
			// 将第二个参数添加到新数组末尾
			newElements[length] = args[1]
			// 返回新数组
			return &object.Array{Elements: newElements}
		},
	},
}
