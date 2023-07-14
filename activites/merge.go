package activites

import (
	"context"

	"github.com/codebdy/minions-go/runtime"
)

type MergeActivity struct {
	Activity   runtime.Activity[any]
	inputCount int
	values     map[string]any
}

func init() {
	runtime.RegisterActivity(
		"merge",
		MergeActivity{},
	)
}

func (s MergeActivity) DynamicInput(name string, inputValue any, ctx context.Context) {
	s.values[name] = inputValue
	s.inputCount++
	if len(s.Activity.BaseActivity.Meta.InPorts) == s.inputCount {
		s.Activity.Output(s.values, ctx)
		s.inputCount = 0
	}
}
