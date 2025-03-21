package core

import (
	"fmt"
	"strings"
)

func String(pattern string) Parser {
	return Parser{
		StateTransformerFn: func(state State) State {
			if state.IsError {
				return state
			}

			target := state.Target[state.Index:]

			if len(target) == 0 {
				return state.SetError(fmt.Sprintf("ParseError: Expected %s but input ended unexpectedly", pattern))
			}

			if strings.HasPrefix(target, pattern) {
				return state.SetResult(pattern).SetIndex(state.Index + len(pattern))
			}

			return state.SetError(
				fmt.Sprintf("ParseError: Expected %s but got %s...", pattern, state.Target[state.Index:]),
			)
		},
	}
}

func Sequence(parsers []Parser) Parser {
	return Parser{
		StateTransformerFn: func(state State) State {
			if state.IsError {
				return state
			}

			var results []string
			nextState := state

			for _, parser := range parsers {
				nextState = parser.StateTransformerFn(nextState)
				results = append(results, nextState.Result.(string))
			}
			return nextState.SetResult(results)
		},
	}
}
