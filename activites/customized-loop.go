package activites

import (
	"fmt"

	"github.com/codebdy/minions-go/runtime"
)

type CustomizedLoopConfig struct {
	FromInput bool `json:"fromInput"`
	Times     int  `json:"times"`
}
type CustomizedLoopActivity struct {
	BaseActivity runtime.BaseActivity[CustomizedLoopConfig]
}

const CUSTOMIZED_LOOP_PORT_OUTPUT = "output"
const CUSTOMIZED_LOOP_PORT_FINISHED = "finished"

func init() {
	runtime.RegisterActivity(
		"debug",
		CustomizedLoopActivity{},
	)
}

func (l CustomizedLoopActivity) Input(inputValue any) {
	if l.BaseActivity.GetConfig().FromInput {
		if inputValue == nil {
			fmt.Println("CustomizedLoop input is nil")
		} else {
			for _, one := range inputValue.([]any) {
				l.Output(one)
			}
		}
	} else if l.BaseActivity.GetConfig().Times > 0 {
		for i := 0; i < l.BaseActivity.GetConfig().Times; i++ {
			l.Output(i)
		}
	}
	l.BaseActivity.Next(inputValue, CUSTOMIZED_LOOP_PORT_FINISHED)
}

func (l CustomizedLoopActivity) Output(value any) {
	l.BaseActivity.Next(value, CUSTOMIZED_LOOP_PORT_OUTPUT)
}
