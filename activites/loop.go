package activites

import (
	"fmt"

	"github.com/codebdy/minions-go/runtime"
)

type LoopConfig struct {
	FromInput bool `json:"fromInput"`
	Times     int  `json:"times"`
}
type LoopActivity struct {
	Activity runtime.Activity[LoopConfig]
}

const LOOP_PORT_OUTPUT = "output"
const LOOP_PORT_FINISHED = "finished"

func init() {
	runtime.RegisterActivity(
		"debug",
		LoopActivity{},
	)
}

func (l LoopActivity) Input(inputValue any) {
	config := l.Activity.GetConfig()
	if config.FromInput {
		if inputValue == nil {
			fmt.Println("Loop input is nil")
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
	l.Activity.Next(inputValue, LOOP_PORT_FINISHED)
}

func (l LoopActivity) Output(value any) {
	l.Activity.Next(value, LOOP_PORT_OUTPUT)
}
