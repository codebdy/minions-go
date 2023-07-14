package activites

import (
	"context"

	"github.com/codebdy/minions-go/runtime"
)

type ContextWriteConfig struct {
	Name string `json:"name"`
}

type ContextWriteActivity struct {
	Activity runtime.Activity[ContextWriteConfig]
}

func init() {
	runtime.RegisterActivity(
		"debug",
		DebugActivity{},
	)
}

func (r ContextWriteActivity) Input(inputValue any, ctx context.Context) {
	config := r.Activity.GetConfig()
	r.Activity.Output(inputValue, context.WithValue(ctx, config.Name, inputValue))
}
