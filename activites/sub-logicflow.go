package activites

import (
	"github.com/codebdy/minions-go/dsl"
	"github.com/codebdy/minions-go/runtime"
)

type SubLogicFlowConfig struct {
	subLogicFlowId string
}

type SubLogicFlowActivity struct {
	BaseActivity runtime.BaseActivity
}

func init() {
	runtime.RegisterActivity(
		"subLogicFlow",
		SubLogicFlowActivity{},
	)
}

//该方法如果存在，会通过反射被自动调用
func (s SubLogicFlowActivity) Init(meta *dsl.ActivityDefine) {
	metas := (&s.BaseActivity).Ctx.Value(runtime.CONTEXT_KEY_SUBMETAS)
	if metas != nil {
		flowMeta := s.GetFlowMeta()
		if flowMeta != nil {
			logicFlow := runtime.NewLogicflow(*flowMeta, s.BaseActivity.Ctx)
			s.BaseActivity.Jointers = logicFlow.Jointers
		}
	}
}

func (s SubLogicFlowActivity) GetFlowMeta() *dsl.LogicFlowMeta {
	metas := (&s.BaseActivity).Ctx.Value(runtime.CONTEXT_KEY_SUBMETAS)
	if metas != nil {
		logicFlowMetas := metas.(*[]dsl.SubLogicFlowMeta)
		for i := range *logicFlowMetas {
			flowMeta := (*logicFlowMetas)[i]
			if flowMeta.Id == s.GetConfig().subLogicFlowId {
				return &flowMeta.LogicFlowMeta
			}
		}
	}

	return nil
}

func (s SubLogicFlowActivity) GetConfig() SubLogicFlowConfig {
	return runtime.GetActivityConfig[SubLogicFlowConfig](&s.BaseActivity)
}
