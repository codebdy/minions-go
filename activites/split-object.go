package activites

import (
	"context"

	"github.com/codebdy/minions-go/runtime"
)

type SplitObjectActivity struct {
	Activity runtime.Activity[any]
}

func init() {
	runtime.RegisterActivity(
		"splitObject",
		SplitObjectActivity{},
	)
	//参数拆分同样的处理
	runtime.RegisterActivity(
		"splitArgs",
		SplitObjectActivity{},
	)
}

func (s SplitObjectActivity) Input(inputValue any, ctx context.Context) {
	if inputValue != nil {
		valueMap := inputValue.(map[string]any)

		for _, output := range s.Activity.BaseActivity.Meta.OutPorts {
			if output.Name != "" {
				s.Activity.Next(valueMap[output.Name], output.Name, ctx)
			}
		}
	}
}
