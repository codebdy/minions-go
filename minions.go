package minions

import (
	"context"

	_ "github.com/codebdy/minions-go/activites"
	"github.com/codebdy/minions-go/dsl"
	"github.com/codebdy/minions-go/runtime"
)

//注册元件
func RegisterActivity(name string, activity interface{}) {
	runtime.RegisterActivity(name, activity)
}

func AttachSubFlowsToContext(flowMetas *[]dsl.SubLogicFlowMeta, ctx context.Context) context.Context {
	return runtime.AttachSubFlowsToContext(flowMetas, ctx)
}

// ctx用于传递value，minions.CONTEXT_KEY_SUBMETAS 对应*[]dsl.LogicFlowDefine， 子编排metas
func NewLogicflow(flowMeta dsl.LogicFlowMeta, ctx context.Context) *runtime.LogicFlow {
	return runtime.NewLogicflow(flowMeta, ctx)
}
