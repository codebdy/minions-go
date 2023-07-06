package dsl

type LogicFlowDefine struct {
	Id    string           `json:"id"`
	Name  string           `json:"name"`
	Label string           `json:"label"`
	Nodes []ActivityDefine `json:"nodes"`
	Lines []LineDefine     `json:"lines"`
}
