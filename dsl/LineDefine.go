package dsl

//连线定义
type LineDefine struct {
	Id string `json:"id"`
	//源节点
	Source PortRef `json:"source"`
	//目标节点
	Target PortRef `json:"target"`
}
