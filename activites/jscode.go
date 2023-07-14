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
		fnStr := "const codeFunc = " + config.Expression + `
		  function callFunc() {
				codeFunc(inputs, outputs, context)
			}
		`
		_, err := vm.RunString(fnStr)
		if err != nil {
			panic(err)
		}

		outputHandlers := map[string]PortFunc{}

		for _, output := range j.Activity.BaseActivity.Meta.OutPorts {
			outputHandlers[output.Name] = func(inputValue any, ctx context.Context) {
				if ctx == nil {
					ctx = context.Background()
				}
				j.Activity.Next(inputValue, output.Name, ctx)
			}
		}

		vm.Set("inputs", j.values)
		vm.Set("outputs", outputHandlers)
		vm.Set("context", ctx)

		var callFunc func()
		err = vm.ExportTo(vm.Get("callFunc"), &callFunc)
		if err != nil {
			panic(err)
		}

		callFunc()
		j.inputCount = 0
	}

}
