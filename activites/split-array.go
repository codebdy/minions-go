package activites

import (
	"context"
	"fmt"

	"github.com/codebdy/minions-go/runtime"
)

type SplitArrayActivity struct {
	Activity runtime.Activity[any]
}

func init() {
	runtime.RegisterActivity(
		"splitArray",
		SplitArrayActivity{},
	)
}

func (s SplitArrayActivity) Input(inputValue any, ctx context.Context) {
	if inputValue != nil {
		valueArray, ok := inputValue.([]any)
		if !ok {
			fmt.Println("SplitArrayActivity input value is not []any")
		}

		for i := range valueArray {
			if i < len(s.Activity.BaseActivity.Meta.OutPorts) {
				output := s.Activity.BaseActivity.Meta.OutPorts[i]
				s.Activity.Next(valueArray[i], output.Name, ctx)
			}
		}
	}
}
