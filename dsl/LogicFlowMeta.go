package dsl

type LogicFlowMeta struct {
	Nodes []ActivityDefine `json:"nodes"`
	Lines []LineDefine     `json:"lines"`
}

type SubLogicFlowMeta struct {
	LogicFlowMeta
	Id string
}
