package activites

import (
	"context"
	"time"

	"github.com/codebdy/minions-go/runtime"
)

type DelayConfig struct {
	//单位是秒
	Time int `json:"time"`
}
type DelayActivity struct {
	Activity runtime.Activity[DelayConfig]
}

func init() {
	runtime.RegisterActivity(
		"delay",
		DelayActivity{},
	)
}

func (d *DelayActivity) Input(inputValue any, ctx context.Context) {
	config := d.Activity.GetConfig()
	time.Sleep(time.Duration(config.Time) * time.Second)
	d.Activity.Output(inputValue, ctx)
}
