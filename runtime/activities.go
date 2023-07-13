package runtime

import (
	"reflect"
)

var activitiesMap map[string]reflect.Type = make(map[string]reflect.Type)

//func RegisterActivity[Config any, T Activity[Config]](name string, factory func(meta *dsl.ActivityDefine) *T) {
func RegisterActivity(name string, activity interface{}) {
	if activitiesMap[name] != nil {
		panic("Activity dupliated:" + name)
	}
	activitiesMap[name] = reflect.TypeOf(activity)
}
