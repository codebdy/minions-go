package activites

import (
	"context"
	"fmt"

	"github.com/codebdy/minions-go/dsl"
	"github.com/codebdy/minions-go/runtime"
)

type CustomizedLoopConfig struct {
	FromInput bool `json:"fromInput"`
	Times     int  `json:"times"`
}
type CustomizedLoopActivity struct {
	Activity  runtime.Activity[CustomizedLoopConfig]
	count     int
	logicFlow *runtime.LogicFlow
	finished  bool
}

const CUSTOMIZED_LOOP_PORT_OUTPUT = "output"
const CUSTOMIZED_LOOP_PORT_FINISHED = "finished"

func init() {
	runtime.RegisterActivity(
		"customizedLoop",
		CustomizedLoopActivity{},
	)
}

//该方法如果存在，会通过反射被自动调用
func (l *CustomizedLoopActivity) Init() {
	if l.Activity.BaseActivity.Meta != nil {
		l.logicFlow = runtime.NewLogicflow(l.Activity.BaseActivity.Meta.Children, l.Activity.BaseActivity.Ctx)
		outputMeta := getOutputByName(CUSTOMIZED_LOOP_PORT_OUTPUT, l.Activity.BaseActivity.Meta.OutPorts)
		if outputMeta == nil {
			panic("No output port in CustomizedLoop")
		}
		jointer := l.logicFlow.Jointers.GetOutputById(outputMeta.Id)
		if jointer == nil {
			panic("No output jointer in CustomizedLoop")
		}
		jointer.Connect(l.Output)

		finishedMeta := getOutputByName(CUSTOMIZED_LOOP_PORT_FINISHED, l.Activity.BaseActivity.Meta.OutPorts)
		if finishedMeta == nil {
			panic("No finished port in CustomizedLoop")
		}

		jointer = l.logicFlow.Jointers.GetOutputById(finishedMeta.Id)
		if jointer == nil {
			panic("No finished jointer in CustomizedLoop")
		}
		jointer.Connect(l.Finshed)
	}
}

func (l *CustomizedLoopActivity) Input(inputValue any, ctx context.Context) {
	config := l.Activity.GetConfig()
	if config.FromInput {
		if inputValue == nil {
			fmt.Println("CustomizedLoop input is nil")
		} else {
			for _, one := range inputValue.([]any) {
				l.logicFlow.Jointers.GetSingleInput().Push(one, ctx)
				l.count++
				if l.finished {
					break
				}
			}
		}
	} else if config.Times > 0 {
		for i := 0; i < config.Times; i++ {
			l.logicFlow.Jointers.GetSingleInput().Push(i, ctx)
			l.count++
			if l.finished {
				break
			}
		}
	}
	if !l.finished {
		l.Activity.Next(l.count, CUSTOMIZED_LOOP_PORT_FINISHED, ctx)
	}
}

func (l *CustomizedLoopActivity) Output(value any, ctx context.Context) {
	l.Activity.Next(value, CUSTOMIZED_LOOP_PORT_OUTPUT, ctx)
}

func (l *CustomizedLoopActivity) Finshed(value any, ctx context.Context) {
	l.finished = true
	//启动一个新的协程去处理
	//go l.Activity.Next(value, CUSTOMIZED_LOOP_PORT_FINISHED, ctx)
	l.Activity.Next(value, CUSTOMIZED_LOOP_PORT_FINISHED, ctx)
}

func getOutputByName(name string, ports []dsl.PortDefine) *dsl.PortDefine {
	for _, port := range ports {
		if port.Name == name {
			return &port
		}
	}

	return nil
}
