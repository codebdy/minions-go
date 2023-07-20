package dsl

//线的连接点
type PortRef struct {
	//节点ID
	NodeId string `json:"nodeId"`
	//端口ID
	PortId string `json:"portId"`
}
