package dsl

type LogicFlowMeta struct {
	Nodes []NodeDefine `json:"nodes"`
	Lines []LineDefine `json:"lines"`
}

type SubLogicFlowMeta struct {
	LogicFlowMeta
	Id string
}
