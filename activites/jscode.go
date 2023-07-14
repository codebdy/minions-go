package activites

import (
	"context"

	"github.com/codebdy/minions-go/runtime"
	"github.com/dop251/goja"
)

type JsCodeConfig struct {
	Expression string `json:"expression"`
}

type JsCodeActivity struct {
	Activity   runtime.Activity[JsCodeConfig]
	inputCount int
	values     map[string]any
}

func init() {
	runtime.RegisterActivity(
		"jsCode",
		JsCodeActivity{},
	)
}

type PortFunc = func(inputValue any, ctx context.Context)

func (j *JsCodeActivity) Init() {
	j.values = map[string]any{}
}

func (j *JsCodeActivity) DynamicInput(name string, inputValue any, ctx context.Context) {
	j.values[name] = inputValue
	j.inputCount++
	if len(j.Activity.BaseActivity.Meta.InPorts) == j.inputCount {
		config := j.Activity.GetConfig()
		if config.Expression == "" {
			panic("Not set code to jsCode activity")
		}
		vm := goja.New()
		vm.SetFieldNameMapper(goja.UncapFieldNameMapper())
		_, err := vm.RunString("const codeFunc = " + config.Expression)
		if err != nil {
			panic(err)
		}

		var codeFunc func(map[string]any, map[string]PortFunc, context.Context)
		err = vm.ExportTo(vm.Get("codeFunc"), &codeFunc)
		if err != nil {
			panic(err)
		}

		outputHandlers := map[string]PortFunc{}

		for _, output := range j.Activity.BaseActivity.Meta.OutPorts {
			outputHandlers[output.Name] = func(inputValue any, ctx context.Context) {
				j.Activity.Next(inputValue, output.Name, ctx)
			}
		}

		codeFunc(j.values, outputHandlers, ctx)
		j.inputCount = 0
	}

}
