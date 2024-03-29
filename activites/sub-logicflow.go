package activites

import (
	"github.com/codebdy/minions-go/dsl"
	"github.com/codebdy/minions-go/runtime"
)

type LogicFlowConfigParam struct {
	LogicFlowId string `json:"logicFlowId"`
}

type SubLogicFlowConfig struct {
	Param LogicFlowConfigParam `json:"param"`
}

type SubLogicFlowActivity struct {
	Activity runtime.Activity[SubLogicFlowConfig]
}

func init() {
	runtime.RegisterActivity(
		"subLogicFlow",
		SubLogicFlowActivity{},
	)
}

//该方法如果存在，会通过反射被自动调用
func (s *SubLogicFlowActivity) Init() {
	metas := (&s.Activity).BaseActivity.Ctx.Value(runtime.CONTEXT_KEY_SUBMETAS)
	if metas != nil {
		flowMeta := s.GetFlowMeta()
		if flowMeta != nil {
			logicFlow := runtime.NewLogicflow(*flowMeta, s.Activity.BaseActivity.Ctx)
			s.Activity.BaseActivity.Jointers = logicFlow.Jointers
		}
	}
}

func (s *SubLogicFlowActivity) GetFlowMeta() *dsl.LogicFlowMeta {
	metas := (&s.Activity.BaseActivity).Ctx.Value(runtime.CONTEXT_KEY_SUBMETAS)
	if metas != nil {
		logicFlowMetas := metas.(*[]dsl.SubLogicFlowMeta)
		for i := range *logicFlowMetas {
			flowMeta := (*logicFlowMetas)[i]
			if flowMeta.Id == s.Activity.GetConfig().Param.LogicFlowId {
				return &flowMeta.LogicFlowMeta
			}
		}
	}

	return nil
}
