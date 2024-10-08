package evaluator

import (
	"fmt"
	"limLang/object"
)

var buildtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			// return NULL
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},

	"print": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Print(arg.Inspect())
			}
			return nil
		},
	},

	// "push": &object.Builtin{
	// 	Fn: func(args ...object.Object) object.Object {
	// 		fmt.Println("wtf is in push ")
	// 		if len(args) != 2 {
	// 			return newError("wrong number of arguments. got=%d, want=2",
	// 				len(args))
	// 		}
	// 		if args[0].Type() != object.ARRAY_OBJ {
	// 			return newError("argument to `push` must be ARRAY, got %s", args[0].Type())
	// 		}
	// 		arr := args[0].(*object.Array)
	// 		length := len(arr.Elements)
	// 		newElements := make([]object.Object, length+1, length+1)
	// 		copy(newElements, arr.Elements)
	// 		newElements[length] = args[1]
	// 		// return &object.Array{Elements: newElements}
	// 		return &object.Array{Elements: []object.Object{&object.Integer{Value: 1}}}
	// 		// return &object.Array{}
	// 	},
	// },
}
