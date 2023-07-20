package dsl

//一段逻辑编排
type LogicFlowMeta struct {
	//所有节点
	Nodes []NodeDefine `json:"nodes"`
	//所有连线
	Lines []LineDefine `json:"lines"`
}

//子编排，可以被其它编排调用
type SubLogicFlowMeta struct {
	LogicFlowMeta
	Id string
}
