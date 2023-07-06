package runtime

import "github.com/codebdy/minions-go/dsl"

type LogicFlow struct {
	Id       string
	Jointers *ActivityJointers
}

func NewLogicflow(flowMeta dsl.LogicFlowDefine) *LogicFlow {
	var logicFlow LogicFlow
	logicFlow.Id = flowMeta.Id
	logicFlow.Jointers = &ActivityJointers{}
	return &logicFlow
}
