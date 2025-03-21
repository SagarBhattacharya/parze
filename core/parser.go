package core

type StateTransformerFn = func(State) State

type Parser struct {
	StateTransformerFn
}

func (p Parser) Map(fn func(result any) any) Parser {
	return Parser{
		StateTransformerFn: func(state State) State {
			nextState := p.StateTransformerFn(state)
			if nextState.IsError {
				return nextState
			}
			return nextState.SetResult(fn(nextState.Result))
		},
	}
}

func (p Parser) ErrorMap(fn func(message string, index int) string) Parser {
	return Parser{
		StateTransformerFn: func(state State) State {
			nextState := p.StateTransformerFn(state)
			if !nextState.IsError {
				return nextState
			}
			return nextState.SetError(fn(nextState.Error, nextState.Index))
		},
	}
}

func (p *Parser) Run(target string) State {
	initialState := State{
		IsError: false,
		Target:  target,
		Index:   0,
		Result:  nil,
		Error:   "",
	}

	return p.StateTransformerFn(initialState)
}
