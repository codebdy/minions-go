package runtime

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

type AbstractActivity struct {
	Id       string
	Jointers *ActivityJointers
	Config   map[string]interface{}
}

func (a *AbstractActivity) Next(inputValue interface{}, outputName string) {
	if outputName == "" {
		outputName = "output"
	}
	nextJointer := a.Jointers.GetOutput(outputName)
	if nextJointer != nil {
		nextJointer.Push(inputValue)
	}
}
