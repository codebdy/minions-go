package activites

import (
	minions "github.com/codebdy/minions-go"
	"github.com/codebdy/minions-go/dsl"
	"github.com/codebdy/minions-go/runtime"
)

type SubLogicFlowConfig struct {
	subLogicFlowId string
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
	metas := (&s.BaseActivity).Ctx.Value(minions.CONTEXT_KEY_SUBMETAS)
	if metas != nil {
		flowMeta := s.GetFlowMeta()
		if flowMeta != nil {
			logicFlow := runtime.NewLogicflow(*flowMeta, s.BaseActivity.Ctx)
			s.BaseActivity.Jointers = logicFlow.Jointers
		}
	}
}

func (s SubLogicFlowActivity) GetFlowMeta() *dsl.LogicFlowDefine {
	metas := (&s.BaseActivity).Ctx.Value(minions.CONTEXT_KEY_SUBMETAS)
	if metas != nil {
		logicFlowMetas := metas.([]dsl.LogicFlowDefine)
		for i := range logicFlowMetas {
			flowMeta := logicFlowMetas[i]
			if flowMeta.Id == s.BaseActivity.GetConfig().subLogicFlowId {
				return &flowMeta
			}
		}
	}

	return nil
}
