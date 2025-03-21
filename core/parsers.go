package core

import (
	"fmt"
	"regexp"
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

func Letters() Parser {
	return Parser{
		StateTransformerFn: func(state State) State {
			if state.IsError {
				return state
			}

			target := state.Target[state.Index:]

			if len(target) == 0 {
				return state.SetError(fmt.Sprintf("ParseError: Expected 'letter' but input ended unexpectedly"))
			}

			re := regexp.MustCompile(`^[A-Za-z]+`)
			result := re.FindAllString(target, -1)
			if result != nil {
				return state.SetResult(result[0]).SetIndex(state.Index + len(result[0]))
			}

			return state.SetError(fmt.Sprintf("ParseError: Expected letter at %d", state.Index))
		},
	}
}

func Digits() Parser {
	return Parser{
		StateTransformerFn: func(state State) State {
			if state.IsError {
				return state
			}

			target := state.Target[state.Index:]

			if len(target) == 0 {
				return state.SetError(fmt.Sprintf("ParseError: Expected 'digit' but input ended unexpectedly"))
			}

			re := regexp.MustCompile(`^[0-9]+`)
			result := re.FindAllString(target, -1)
			if result != nil {
				return state.SetResult(result[0]).SetIndex(state.Index + len(result[0]))
			}

			return state.SetError(fmt.Sprintf("ParseError: Expected digit at %d", state.Index))
		},
	}
}

func And(parsers []Parser) Parser {
	return Parser{
		StateTransformerFn: func(state State) State {
			if state.IsError {
				return state
			}

			var results []any
			nextState := state

			for _, parser := range parsers {
				nextState = parser.StateTransformerFn(nextState)
				results = append(results, nextState.Result)
			}
			return nextState.SetResult(results)
		},
	}
}

func Or(parsers []Parser) Parser {
	return Parser{
		StateTransformerFn: func(state State) State {
			if state.IsError {
				return state
			}
			var nextState State

			for _, parser := range parsers {
				nextState = parser.StateTransformerFn(state)
				if !nextState.IsError {
					return nextState
				}
			}
			return nextState.SetError(fmt.Sprintf("ParseError: No Matches Found at %d", nextState.Index))
		},
	}
}

func Many(parser Parser) Parser {
	return Parser{
		StateTransformerFn: func(state State) State {
			if state.IsError {
				return state
			}
			done := false
			var results []any
			nextState := state

			for !done {
				testState := parser.StateTransformerFn(nextState)
				if !testState.IsError {
					results = append(results, testState.Result)
					nextState = testState
				} else {
					done = true
				}
			}

			return nextState.SetResult(results)
		},
	}
}

func ManyOne(parser Parser) Parser {
	return Parser{
		StateTransformerFn: func(state State) State {
			if state.IsError {
				return state
			}
			done := false
			var results []any
			nextState := state

			for !done {
				testState := parser.StateTransformerFn(nextState)
				if !testState.IsError {
					results = append(results, testState.Result)
					nextState = testState
				} else {
					done = true
				}
			}

			if len(results) < 1 {
				return nextState.SetError(fmt.Sprintf("ParseError: At least one match expected but found none at %d", nextState.Index))
			}

			return nextState.SetResult(results)
		},
	}
}

func Between(left Parser, content Parser, right Parser) Parser {
	return And([]Parser{
		left,
		content,
		right,
	}).Map(func(result any) any {
		return result.([]any)[1]
	})
}
