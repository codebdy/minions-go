package runtime

import (
	"context"

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

func (j *ActivityJointers) GetSingleInput() *Jointer {
	if len(j.inputs) > 0 {
		return j.inputs[0]
	}
	return nil
}

func (j *ActivityJointers) GetSingleOutput() *Jointer {
	if len(j.outputs) > 0 {
		return j.outputs[0]
	}
	return nil
}

func (j *ActivityJointers) GetOutputById(id string) *Jointer {
	for _, jointer := range j.outputs {
		if jointer.Id == id {
			return jointer
		}
	}
	return nil
}

func (j *ActivityJointers) GetInputById(id string) *Jointer {
	for _, jointer := range j.inputs {
		if jointer.Id == id {
			return jointer
		}
	}
	return nil
}

type BaseActivity struct {
	Id       string
	Jointers *ActivityJointers
	Meta     *dsl.ActivityDefine
	Ctx      context.Context
}

//为了处理Config添加的一层，这个在反射时不能被显式类型转换
type Activity[Config any] struct {
	BaseActivity BaseActivity
}

// type Activity[Config any] interface {
// 	GetBaseActivity() *BaseActivity[Config]
// }

// func NewActivity[Config any, T Activity[Config]](meta *dsl.ActivityDefine) *T {
// 	var activity T
// 	activity.GetBaseActivity().Init(meta)
// 	return &activity
// }

func (b *BaseActivity) Init(meta *dsl.ActivityDefine, ctx context.Context) {
	b.Meta = meta
	b.Id = meta.Id
	b.Jointers = &ActivityJointers{}
	b.Ctx = ctx
}

func (b *BaseActivity) Next(inputValue interface{}, outputName string) {
	if outputName == "" {
		outputName = "output"
	}
	nextJointer := b.Jointers.GetOutput(outputName)
	if nextJointer != nil {
		nextJointer.Push(inputValue, b.Ctx)
	}
}

func (a *Activity[Config]) GetConfig() Config {
	var config Config

	if a.BaseActivity.Meta.Config != nil {
		mapstructure.Decode(a.BaseActivity.Meta.Config, &config)
	}
	return config
}

func (a *Activity[Config]) Next(inputValue any, outputName string) {
	a.BaseActivity.Next(inputValue, outputName)
}

func (a *Activity[Config]) Output(inputValue any) {
	a.BaseActivity.Next(inputValue, "output")
}
