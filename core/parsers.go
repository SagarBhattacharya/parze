package core

import (
	"fmt"
	"regexp"
	"strings"
)

func Lazy(thunk func() Parser) Parser {
	return Parser{
		StateTransformerFn: func(state State) State {
			parser := thunk()
			return parser.StateTransformerFn(state)
		},
	}
}

func String(pattern string) Parser {
	return Lazy(func() Parser {
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
	})
}

func Letters() Parser {
	return Lazy(func() Parser {
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
	})
}

func Digits() Parser {
	return Lazy(func() Parser {
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
	})
}

func And(parsers []Parser) Parser {
	return Lazy(func() Parser {
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
	})
}

func Or(parsers []Parser) Parser {
	return Lazy(func() Parser {
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
	})
}

func Many(parser Parser) Parser {
	return Lazy(func() Parser {
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
	})
}

func ManyOne(parser Parser) Parser {
	return Lazy(func() Parser {
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
	})
}

func Seperated(seperatorParser Parser, valueParser Parser) Parser {
	return Lazy(func() Parser {
		return Parser{
			StateTransformerFn: func(state State) State {
				if state.IsError {
					return state
				}

				var results []any
				nextState := state

				for {
					contentState := valueParser.StateTransformerFn(nextState)
					if contentState.IsError {
						break
					}
					results = append(results, contentState.Result)
					nextState = contentState

					seperatorState := seperatorParser.StateTransformerFn(nextState)
					if seperatorState.IsError {
						break
					}
					nextState = seperatorState
				}

				return nextState.SetResult(results)
			},
		}
	})
}

func SeperatedOne(seperatorParser Parser, valueParser Parser) Parser {
	return Lazy(func() Parser {
		return Parser{
			StateTransformerFn: func(state State) State {
				if state.IsError {
					return state
				}

				var results []any
				nextState := state

				for {
					contentState := valueParser.StateTransformerFn(nextState)
					if contentState.IsError {
						break
					}
					results = append(results, contentState.Result)
					nextState = contentState

					seperatorState := seperatorParser.StateTransformerFn(nextState)
					if seperatorState.IsError {
						break
					}
					nextState = seperatorState
				}

				if len(results) < 1 {
					return nextState.SetError(fmt.Sprintf("ParseError: At least one match expected but found none at %d", nextState.Index))
				}

				return nextState.SetResult(results)
			},
		}
	})
}

func WhiteSpace() Parser {
	return Lazy(func() Parser {
		return Many(
			Or([]Parser{
				String(" "),
				String("\t"),
				String("\n"),
				String("\r"),
			}),
		).Map(func(result any) any {
			return ""
		})
	})
}

func Between(brackets string, content Parser) Parser {
	return Lazy(func() Parser {
		return And([]Parser{
			WhiteSpace(),
			String(string(brackets[0])),
			WhiteSpace(),
			content,
			WhiteSpace(),
			String(string(brackets[1])),
			WhiteSpace(),
		}).Map(func(result any) any {
			return result.([]any)[3]
		})
	})
}

func Optional(parser Parser) Parser {
	return Lazy(func() Parser {
		return Parser{
			StateTransformerFn: func(state State) State {
				nextState := parser.StateTransformerFn(state)
				if nextState.IsError {
					return state.SetResult(nil)
				}
				return nextState
			},
		}
	})
}

func fromAnyArraytoStringArray(array any) []string {
	str, _ := array.([]any)
	var results []string
	for _, val := range str {
		results = append(results, val.(string))
	}
	return results
}

func Number() Parser {
	return Lazy(func() Parser {
		return And([]Parser{
			Optional(String("-")),
			Digits(),
			Optional(And([]Parser{
				String("."),
				Digits(),
			})).Map(func(result any) any {
				if result == nil {
					return result
				}
				return strings.Join(fromAnyArraytoStringArray(result.([]any)), "")
			}),
		}).Map(func(result any) any {
			var value string
			res := result.([]any)
			if res[0] != nil {
				value += res[0].(string)
			}
			value += res[1].(string)
			if res[2] != nil {
				value += res[2].(string)
			}
			return value
		})
	})
}
