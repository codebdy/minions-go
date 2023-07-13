package activites

import (
	"context"
	"fmt"

	"github.com/codebdy/minions-go/runtime"
)

type DebugConfig struct {
	Tip    string `json:"tip"`
	Closed bool   `json:"closed"`
}
type DebugActivity struct {
	BaseActivity runtime.BaseActivity
}

func init() {
	runtime.RegisterActivity(
		"debug",
		DebugActivity{},
	)
}

func (d DebugActivity) Input(inputValue any, ctx context.Context) {
	config := d.GetConfig()
	if !config.Closed {
		tip := "Debug"
		if config.Tip != "" {
			tip = config.Tip
		}
		text := ""
		if inputValue != nil {
			text = inputValue.(string)
		}
		fmt.Print("🪲" + tip + ":" + text)
	}
}

func (l DebugActivity) GetConfig() DebugConfig {
	return runtime.GetActivityConfig[DebugConfig](&l.BaseActivity)
}
