package runtime

import (
	"context"
)

type InputHandler = func(inputValue any, ctx context.Context)

type Jointer struct {
	Id      string
	Name    string
	outlets []InputHandler
}

func (j *Jointer) Push(inputValue any, ctx context.Context) {
	for _, outputHandler := range j.outlets {
		outputHandler(inputValue, ctx)
	}
}

func (j *Jointer) Connect(inuptHandler InputHandler) {
	j.outlets = append(j.outlets, inuptHandler)
}
