package main

import "github.com/SagarBhattacharya/parze/core"

// "string:hello"
// "number:48"

func main() {
	parser := core.And([]core.Parser{
		core.And([]core.Parser{
			core.Letters(),
			core.String(":"),
		}).Map(func(result any) any {
			return result.([]any)[0]
		}).Then(func(result any) core.Parser {
			switch result.(string) {
			case "string":
				return core.Letters()
			case "number":
				return core.Digits()
			default:
				return core.Parser{}
			}
		}),
	})

	state := parser.Run("number:48")
	state.Display()
}
