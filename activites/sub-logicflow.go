package activites

import "github.com/codebdy/minions-go/runtime"

type SubLogicFlowConfig struct {
	Tip    string `json:"tip"`
	Closed bool   `json:"closed"`
}

type SubLogicFlowActivity struct {
	BaseActivity runtime.BaseActivity[SubLogicFlowConfig]
}

func init() {
	runtime.RegisterActivity(
		"debug",
		SubLogicFlowActivity{},
	)
}

//该方法如果存在，会通过反射被自动调用
func (s SubLogicFlowActivity) Init() {

}
