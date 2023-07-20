package dsl

//节点定义
type NodeDefine struct {
	Id           string                 `json:"id"`
	//名称
	Name         string                 `json:"name"` //嵌入编排，端口转换成子节点时使用
	//节点类型，对应Typescript的枚举
	Type         string                 `json:"type"`
	//元件对应Activity名称
	ActivityName string                 `json:"activityName"`
	//标题
	Label        string                 `json:"label"`
	//配置
	Config       map[string]interface{} `json:"config"`
	//入端口
	InPorts      []PortDefine           `json:"inPorts"`
	//出端口
	OutPorts     []PortDefine           `json:"outPorts"`
	//子节点，嵌入式节点用，比如自定义循环节点、事务节点
	Children     LogicFlowMeta          `json:"children"`
}
