package activites

import (
	"fmt"

	"github.com/codebdy/minions-go/runtime"
)

type DebugConfig struct {
	Tip    string `json:"tip"`
	Closed bool   `json:"closed"`
}
type DeubugActivity struct {
	BaseActivity runtime.BaseActivity[DebugConfig]
}

func init() {
	runtime.RegisterActivity(
		"debug",
		runtime.NewActivity[DebugConfig, DeubugActivity],
	)
}

func (d DeubugActivity) GetBaseActivity() *runtime.BaseActivity[DebugConfig] {
	return &d.BaseActivity
}

func (d DeubugActivity) Input(inputValue interface{}) {
	config := d.BaseActivity.GetConfig()
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
