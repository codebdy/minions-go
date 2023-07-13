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
	Activity runtime.Activity[CustomizedLoopConfig]
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
	config := l.Activity.GetConfig()
	if config.FromInput {
		if inputValue == nil {
			fmt.Println("CustomizedLoop input is nil")
		} else {
			for _, one := range inputValue.([]any) {
				l.Output(one)
			}
		}
	} else if config.Times > 0 {
		for i := 0; i < config.Times; i++ {
			l.Output(i)
		}
	}
	l.Activity.Next(inputValue, CUSTOMIZED_LOOP_PORT_FINISHED)
}

func (l CustomizedLoopActivity) Output(value any) {
	l.Activity.Next(value, CUSTOMIZED_LOOP_PORT_OUTPUT)
}
