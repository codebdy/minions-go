package activites

import (
	"context"

	"github.com/codebdy/minions-go/runtime"
)

type ContextReadConfig struct {
	Name string `json:"name"`
}

type ContextReadActivity struct {
	Activity runtime.Activity[ContextReadConfig]
}

func init() {
	runtime.RegisterActivity(
		"contextRead",
		ContextReadActivity{},
	)
}

func (r ContextReadActivity) Input(inputValue any, ctx context.Context) {
	config := r.Activity.GetConfig()
	r.Activity.Output(ctx.Value(config.Name), ctx)
}
