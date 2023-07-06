package runtime

import (
	"reflect"

	"github.com/codebdy/minions-go/dsl"
)

var activitiesMap map[string]reflect.Type

//func RegisterActivity[Config any, T Activity[Config]](name string, factory func(meta *dsl.ActivityDefine) *T) {
func RegisterActivity(name string, activity interface{}) {
	activitiesMap[name] = reflect.TypeOf(activity)
}

func makeActivity(meta dsl.ActivityDefine) {

}
