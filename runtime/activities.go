package runtime

import (
	"sync"

	"github.com/codebdy/minions-go/dsl"
)

var activitiesMap sync.Map

//func RegisterActivity[Config any, T Activity[Config]](name string, factory func(meta *dsl.ActivityDefine) *T) {
func RegisterActivity(name string, factory interface{}) {
	activitiesMap.Store(name, factory)
}

func makeActivity(meta dsl.ActivityDefine) {

}
