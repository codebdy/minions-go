package activites

import (
	"context"

	"github.com/codebdy/minions-go/runtime"
)

type ConstValueConfig struct {
	Value any `json:"value"`
}
type ConstValueActivity struct {
	Activity runtime.Activity[ConstValueConfig]
}

func init() {
	runtime.RegisterActivity(
		"constValue",
		ConstValueActivity{},
	)
}

func (d ConstValueActivity) Input(inputValue any, ctx context.Context) {
	config := d.Activity.GetConfig()
	d.Activity.Output(config.Value, ctx)
}
