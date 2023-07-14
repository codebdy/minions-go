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

func (m *MergeActivity) Init() {
	m.values = map[string]any{}
}

func (m *MergeActivity) DynamicInput(name string, inputValue any, ctx context.Context) {
	m.values[name] = inputValue
	m.inputCount++
	if len(m.Activity.BaseActivity.Meta.InPorts) == m.inputCount {
		m.Activity.Output(m.values, ctx)
		m.inputCount = 0
	}
}
