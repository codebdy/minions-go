package activites

import (
	"context"
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
		"loop",
		LoopActivity{},
	)
}

func (l *LoopActivity) Input(inputValue any, ctx context.Context) {
	config := l.Activity.GetConfig()
	if config.FromInput {
		if inputValue == nil {
			fmt.Println("Loop input is nil")
		} else {
			for _, one := range inputValue.([]any) {
				l.Output(one, ctx)
			}
		}
	} else if config.Times > 0 {
		for i := 0; i < config.Times; i++ {
			l.Output(i, ctx)
		}
	}
	l.Activity.Next(inputValue, LOOP_PORT_FINISHED, ctx)
}

func (l *LoopActivity) Output(value any, ctx context.Context) {
	l.Activity.Next(value, LOOP_PORT_OUTPUT, ctx)
}
