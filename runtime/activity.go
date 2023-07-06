package runtime

import (
	"github.com/codebdy/minions-go/dsl"
	"github.com/mitchellh/mapstructure"
)

type ActivityJointers struct {
	inputs  []*Jointer
	outputs []*Jointer
}

func (j *ActivityJointers) GetOutput(name string) *Jointer {
	for _, jointer := range j.outputs {
		if jointer.Name == name {
			return jointer
		}
	}
	return nil
}

func (j *ActivityJointers) GetInput(name string) *Jointer {
	for _, jointer := range j.inputs {
		if jointer.Name == name {
			return jointer
		}
	}
	return nil
}

type BaseActivity[Config any] struct {
	Id       string
	Jointers *ActivityJointers
	Meta     *dsl.ActivityDefine
}

// type Activity[Config any] interface {
// 	GetBaseActivity() *BaseActivity[Config]
// }

// func NewActivity[Config any, T Activity[Config]](meta *dsl.ActivityDefine) *T {
// 	var activity T
// 	activity.GetBaseActivity().Init(meta)
// 	return &activity
// }

func (a BaseActivity[T]) Init(meta *dsl.ActivityDefine) {
	a.Meta = meta
	a.Id = meta.Id
	a.Jointers = &ActivityJointers{}
}

func (b BaseActivity[Config]) Next(inputValue interface{}, outputName string) {
	if outputName == "" {
		outputName = "output"
	}
	nextJointer := b.Jointers.GetOutput(outputName)
	if nextJointer != nil {
		nextJointer.Push(inputValue)
	}
}

func (b BaseActivity[Config]) GetConfig() Config {
	var config Config

	if b.Meta.Config != nil {
		mapstructure.Decode(b.Meta.Config, &config)
	}
	return config
}
