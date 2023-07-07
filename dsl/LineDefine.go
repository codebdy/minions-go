package dsl

type LineDefine struct {
	Id     string  `json:"id"`
	Source PortRef `json:"source"`
	Target PortRef `json:"target"`
}
