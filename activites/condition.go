package activites

import "github.com/codebdy/minions-go/runtime"

type ConditionConfig struct {
	TrueExpression string `json:"trueExpression"`
}

type ConditionActivity struct {
	BaseActivity runtime.BaseActivity[ConditionConfig]
}

func (c ConditionActivity) Input(inputValue interface{}) {
	//config := d.BaseActivity.GetConfig()
}
