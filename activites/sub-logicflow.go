package activites

import (
	"github.com/codebdy/minions-go/dsl"
	"github.com/codebdy/minions-go/runtime"
)

type SubLogicFlowConfig struct {
	metas dsl.LogicFlowDefine
}

type SubLogicFlowActivity struct {
	BaseActivity runtime.BaseActivity[SubLogicFlowConfig]
}

func init() {
	runtime.RegisterActivity(
		"subLogicFlow",
		SubLogicFlowActivity{},
	)
}

//该方法如果存在，会通过反射被自动调用
func (s SubLogicFlowActivity) Init() {

}
