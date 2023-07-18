package activites

import (
	"context"
	"fmt"
	"reflect"

	"github.com/codebdy/minions-go/runtime"
)

type DebugConfig struct {
	Tip    string `json:"tip"`
	Closed bool   `json:"closed"`
}
type DebugActivity struct {
	Activity runtime.Activity[DebugConfig]
}

func init() {
	runtime.RegisterActivity(
		"debug",
		DebugActivity{},
	)
}

func (d *DebugActivity) Input(inputValue any, ctx context.Context) {
	config := d.Activity.GetConfig()
	if !config.Closed {
		tip := "Debug"
		if config.Tip != "" {
			tip = config.Tip
		}
		text := ""
		if inputValue != nil {
			if reflect.TypeOf(inputValue).Kind() == reflect.String {
				text = inputValue.(string)
			} else if reflect.TypeOf(inputValue).Kind() == reflect.Int ||
				reflect.TypeOf(inputValue).Kind() == reflect.Int16 ||
				reflect.TypeOf(inputValue).Kind() == reflect.Int32 ||
				reflect.TypeOf(inputValue).Kind() == reflect.Int64 {
				text = fmt.Sprintf("%d", inputValue)
			} else if reflect.TypeOf(inputValue).Kind() == reflect.Float32 ||
				reflect.TypeOf(inputValue).Kind() == reflect.Float64 {
				text = fmt.Sprintf("%f", inputValue)
			} else if reflect.TypeOf(inputValue).Kind() == reflect.Bool {
				if inputValue.(bool) {
					text = "true"
				} else {
					text = "false"
				}
			} else {
				text = "input is type:" + reflect.TypeOf(inputValue).String()
			}
		}
		fmt.Println("🪲" + tip + ":" + text)
	}
}
