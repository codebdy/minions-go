package runtime

type InputHandler = func(inputValue interface{})

type Jointer struct {
	Id      string
	Name    string
	outlets []InputHandler
}

func (j *Jointer) Push(inputValue interface{}) {
	for _, outputHandler := range j.outlets {
		outputHandler(inputValue)
	}
}

func (j *Jointer) Connect(inuptHandler InputHandler) {
	j.outlets = append(j.outlets, inuptHandler)
}
