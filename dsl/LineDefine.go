package dsl

type LineDefine struct {
	Id     string     `json:"id"`
	Source PortDefine `json:"source"`
	Target PortDefine `json:"target"`
}
