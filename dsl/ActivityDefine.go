package dsl

type ActivityDefine struct {
	Id           string                 `json:"id"`
	Type         string                 `json:"type"`
	ActivityName string                 `json:"activityName"`
	Label        string                 `json:"label"`
	Config       map[string]interface{} `json:"config"`
	InPorts      []PortDefine           `json:"inPorts"`
	OutPorts     []PortDefine           `json:"outPorts"`
}
