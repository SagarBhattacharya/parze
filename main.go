package main

import (
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

func main() {
	var experssion core.Parser
	experssion = core.Lazy(func() core.Parser {
		return core.Or([]core.Parser{
			core.Number().Map(func(result any) any {
				return Value{
					Type:  "number",
					Value: result,
				}
			}),
			core.Between(
				core.String("("),
				core.And([]core.Parser{
					core.Or([]core.Parser{
						core.String("+"),
						core.String("-"),
						core.String("/"),
						core.String("*"),
					}),
					core.WhiteSpace(),
					experssion,
					core.WhiteSpace(),
					experssion,
				}),
				core.String(")"),
			).
				Map(func(result any) any {
					return Value{
						Type: "operation",
						Value: Operation{
							Op:   result.([]any)[0].(string),
							ValA: result.([]any)[2],
							ValB: result.([]any)[4],
						},
					}
				}),
		})
	})

	state := experssion.Run("(+ 10 (/ 40 20))")
	state.Display()
}
