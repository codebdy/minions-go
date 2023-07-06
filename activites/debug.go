package activites

import (
	"fmt"

	"github.com/codebdy/minions-go/runtime"
)

type DebugConfig struct {
	Tip    string `json:"tip"`
	Closed bool   `json:"closed"`
}
type DebugActivity struct {
	BaseActivity runtime.BaseActivity[DebugConfig]
}

// func NewDebug() *DebugActivity {
// 	var debug DebugActivity

// 	return &debug
// }

func init() {
	runtime.RegisterActivity(
		"debug",
		DebugActivity{},
	)
}

// func (d DeubugActivity) GetBaseActivity() *runtime.BaseActivity[DebugConfig] {
// 	return &d.BaseActivity
// }

func (d DebugActivity) Input(inputValue interface{}) {
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
