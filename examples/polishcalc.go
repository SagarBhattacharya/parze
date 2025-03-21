package examples

import (
	"strconv"

	"github.com/SagarBhattacharya/parze/core"
)

type Value struct {
	Type  string `json:"type"`
	Value any    `json:"value"`
}

type Operation struct {
	Op   string `json:"op"`
	ValA any    `json:"a"`
	ValB any    `json:"b"`
}

func Evaluate(node any) int {
	Node := node.(Value)
	switch Node.Type {
	case "number":
		return Node.Value.(int)
	case "operation":
		Op := Node.Value.(Operation)
		switch Op.Op {
		case "+":
			return Evaluate(Op.ValA) + Evaluate(Op.ValB)
		case "-":
			return Evaluate(Op.ValA) - Evaluate(Op.ValB)
		case "*":
			return Evaluate(Op.ValA) * Evaluate(Op.ValB)
		case "/":
			return Evaluate(Op.ValA) / Evaluate(Op.ValB)
		default:
			return 0
		}
	default:
		return 0
	}
}

func Interpreter() {
	var operationParser core.Parser

	numberParser := core.Digits().
		Map(func(result any) any {
			val, _ := strconv.ParseInt(result.(string), 10, 0)
			return Value{
				Type:  "number",
				Value: val,
			}
		})

	operatorParser := core.Or([]core.Parser{
		core.String("+"),
		core.String("-"),
		core.String("/"),
		core.String("*"),
	})

	expr := core.Lazy(func() core.Parser {
		return core.Or([]core.Parser{
			numberParser,
			operationParser,
		})
	})

	operationParser = core.Between(
		core.String("("),
		core.And([]core.Parser{
			operatorParser,
			core.String(" "),
			expr,
			core.String(" "),
			expr,
		}),
		core.String(")"),
	).Map(func(result any) any {
		return Value{
			Type: "operation",
			Value: Operation{
				Op:   result.([]any)[0].(string),
				ValA: result.([]any)[2],
				ValB: result.([]any)[4],
			},
		}
	})

	state := operationParser.Run("(+ 10 (/ 40 20))")
	state.Display()
	// fmt.Printf("Evaluation Result : %d", Evaluate(state.Result))
}
