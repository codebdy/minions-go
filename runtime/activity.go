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

type BaseActivity[T any] struct {
	Id       string
	Jointers *ActivityJointers
	Meta     *dsl.ActivityDefine
}

func (a BaseActivity[T]) Init(meta *dsl.ActivityDefine) {
	a.Meta = meta
	a.Id = meta.Id
}

func (a BaseActivity[T]) Next(inputValue interface{}, outputName string) {
	if outputName == "" {
		outputName = "output"
	}
	nextJointer := a.Jointers.GetOutput(outputName)
	if nextJointer != nil {
		nextJointer.Push(inputValue)
	}
}

func (d BaseActivity[T]) GetConfig() T {
	var config T

	if d.Meta.Config != nil {
		mapstructure.Decode(d.Meta.Config, &config)
	}
	return config
}
