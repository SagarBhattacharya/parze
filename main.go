package main

import (
	"fmt"
	"strings"

	"github.com/SagarBhattacharya/parze/core"
)

type Res struct {
	Value string
}

func main() {
	// parser := core.Sequence([]core.Parser{
	// 	core.String("Hello there!"),
	// 	core.String("Goodbye world!"),
	// })

	parser := core.String("Hello").
		Map(func(result any) any {
			return strings.ToUpper(result.(string))
		}).
		ErrorMap(func(message string, index int) string {
			return fmt.Sprintf("Expected a greeting at %d", index)
		})

	state := parser.Run("Hello world!")
	state.Display()
}
