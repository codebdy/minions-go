package activites

import (
	"context"

	"github.com/codebdy/minions-go/runtime"
	"github.com/dop251/goja"
)

type ConditionConfig struct {
	TrueExpression string `json:"trueExpression"`
}

type ConditionActivity struct {
	Activity runtime.Activity[ConditionConfig]
}

func init() {
	runtime.RegisterActivity(
		"condition",
		ConditionActivity{},
	)
}

func (c ConditionActivity) Input(inputValue interface{}, ctx context.Context) {
	config := c.Activity.GetConfig()
	if inputValue != nil && config.TrueExpression == "" {
		c.Activity.Output(inputValue, ctx)
	} else {
		vm := goja.New()
		vm.Set("inputValue", inputValue)
		v, err := vm.RunString(config.TrueExpression)
		if err != nil {
			panic(err)
		}
		if result := v.Export().(bool); result {
			c.Activity.Next(inputValue, "true", ctx)
		} else {
			c.Activity.Next(inputValue, "false", ctx)
		}
	}
}
