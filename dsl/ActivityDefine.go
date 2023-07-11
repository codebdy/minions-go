package dsl

type LogicFlowMetas struct {
	Nodes []ActivityDefine `json:"nodes"`
	Lines []LineDefine     `json:"lines"`
}

type ActivityDefine struct {
	Id           string                 `json:"id"`
	Name         string                 `json:"name"` //嵌入编排，端口转换成子节点时使用
	Type         string                 `json:"type"`
	ActivityName string                 `json:"activityName"`
	Label        string                 `json:"label"`
	Config       map[string]interface{} `json:"config"`
	InPorts      []PortDefine           `json:"inPorts"`
	OutPorts     []PortDefine           `json:"outPorts"`
	Children     LogicFlowMetas         `json:"children"`
}
