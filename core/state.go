package core

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type State struct {
	IsError bool   `json:"isError"`
	Target  string `json:"target"`
	Index   int    `json:"index"`
	Result  any    `json:"result"`
	Error   string `json:"error"`
}

func (state State) SetIndex(index int) State {
	state.Index = index
	return state
}

func (state State) SetResult(result any) State {
	state.Result = result
	return state
}

func (state State) SetError(message string) State {
	state.IsError = true
	state.Error = message
	return state
}

func (state State) Display() {
	result, err := json.Marshal(state)
	if err != nil {
		fmt.Println(err)
		return
	}
	buffer := &bytes.Buffer{}
	err = json.Indent(buffer, result, "", "  ")
	fmt.Println(buffer.String())
}
