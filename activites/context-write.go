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
		"contextWrite",
		DebugActivity{},
	)
}

func (w *ContextWriteActivity) Input(inputValue any, ctx context.Context) {
	config := w.Activity.GetConfig()
	w.Activity.Output(inputValue, context.WithValue(ctx, config.Name, inputValue))
}
