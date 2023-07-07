package runtime

import (
	"context"
	"reflect"

	"github.com/codebdy/minions-go/dsl"
)

type LogicFlow struct {
	Id       string
	Jointers *ActivityJointers
	flowMeta *dsl.LogicFlowDefine
	//存Activity 指针
	baseActivites []*BaseActivity[any]
	ctx           context.Context
}

func NewLogicflow(flowMeta dsl.LogicFlowDefine, ctx context.Context) *LogicFlow {
	var logicFlow LogicFlow
	logicFlow.Id = flowMeta.Id
	logicFlow.Jointers = &ActivityJointers{}
	logicFlow.flowMeta = &flowMeta
	logicFlow.ctx = ctx
	//第一步，解析节点
	logicFlow.constructActivities()
	//第二步， 构建连接关系
	logicFlow.contructLines()
	return &logicFlow
}

//构建一个图的所有节点
func (l *LogicFlow) constructActivities() {
	for _, activityMeta := range l.flowMeta.Nodes {
		switch activityMeta.Type {
		case dsl.ACTIVITY_TYPE_START:
			name := "input"
			if activityMeta.ActivityName != "" {
				name = activityMeta.ActivityName
			}
			//start只有一个端口，所以name可以跟meta name一样
			l.Jointers.inputs = append(l.Jointers.inputs, &Jointer{Id: activityMeta.Id, Name: name})
		case dsl.ACTIVITY_TYPE_END:
			name := "output"
			if activityMeta.ActivityName != "" {
				name = activityMeta.ActivityName
			}
			//end 只有一个端口，所以name可以跟meta name一样
			l.Jointers.outputs = append(l.Jointers.outputs, &Jointer{Id: activityMeta.Id, Name: name})
		case dsl.ACTIVITY_TYPE_ACTIVITY, dsl.ACTIVITY_TYPE_LOGICFLOWACTIVITY:
			if activityMeta.ActivityName != "" {
				activityType := activitiesMap[activityMeta.ActivityName]
				if activityType == nil {
					panic("Can not find activity by name:" + activityMeta.ActivityName)
				}
				if activityType.Kind() != reflect.Struct {
					panic("expect struct")
				}

				activity := reflect.New(activityType)

				rv := reflect.ValueOf(activity)
				f := rv.FieldByName("BaseActivity")
				if f.IsValid() {
					v := f.Addr().Interface().(*BaseActivity[any])
					//构造Jointers
					for _, out := range activityMeta.OutPorts {
						v.Jointers.outputs = append(v.Jointers.outputs, &Jointer{Id: out.Id, Name: out.Name})
					}

					for _, input := range activityMeta.InPorts {
						v.Jointers.inputs = append(v.Jointers.inputs, &Jointer{Id: input.Id, Name: input.Name})
					}

					for i := range v.Jointers.inputs {
						input := v.Jointers.inputs[i]

						m := f.MethodByName(input.Name)

						//如果Inputhandler存在
						if m.IsValid() {
							mt := m.Type()
							//输入处理函数不能有返回值
							if mt.NumOut() != 0 {
								panic("Input handler can not be have return values")
							}

							input.Connect(func(inputValue any, ctx context.Context) {
								inputs := make([]reflect.Value, mt.NumIn())

								if mt.NumIn() > 0 {
									inputs[0] = reflect.ValueOf(inputValue)
								}
								if mt.NumIn() > 1 {
									inputs[1] = reflect.ValueOf(ctx)
								}
								m.Call(inputs)
							})
						} else {
							//处理动态端口
							m = f.MethodByName("DynamicInput")
							if m.IsValid() {
								mt := m.Type()
								input.Connect(func(inputValue any, ctx context.Context) {
									inputs := make([]reflect.Value, mt.NumIn())
									if mt.NumIn() == 0 {
										panic("DynamicInput must have one args")
									}
									inputs[0] = reflect.ValueOf(input.Name)
									if mt.NumIn() > 1 {
										inputs[1] = reflect.ValueOf(inputValue)
									}
									if mt.NumIn() > 2 {
										inputs[2] = reflect.ValueOf(ctx)
									}
									m.Call(inputs)
								})
							} else {
								panic("Can not find input handler:" + input.Name)
							}
						}
					}

					l.baseActivites = append(l.baseActivites, v)
				}
			}
		}
	}
}

//连接一个图的所有节点，把所有的jointer连起来
func (l *LogicFlow) contructLines() {
	for _, lineMeta := range l.flowMeta.Lines {
		//先判断是否连接本编排的input
		sourceJointer := l.Jointers.GetInputById(lineMeta.Source.NodeId)

		//如果不是，连接Source activity的output
		if sourceJointer == nil && lineMeta.Source.PortId != "" {
			sourceJointer = l.getSourceJointerByPortRef(lineMeta.Source)
		}
		if sourceJointer == nil {
			panic("Can find source jointer")
		}

		//先判断是否连接本编排的output
		targetJointer := l.Jointers.GetOutputById(lineMeta.Target.NodeId)

		//如果不是，连接Targe activity的Input
		if targetJointer == nil && lineMeta.Target.PortId != "" {
			targetJointer = l.geTargetJointerByPortRef(lineMeta.Target)
		}
		if targetJointer == nil {
			panic("Can find target jointer")
		}

		sourceJointer.Connect(targetJointer.Push)
	}
}
func (l *LogicFlow) geTargetJointerByPortRef(portRef dsl.PortRef) *Jointer {
	targetActivity := l.getActivityById(portRef.NodeId)
	if targetActivity != nil {
		return targetActivity.Jointers.GetInputById(portRef.PortId)
	}
	return nil
}

func (l *LogicFlow) getSourceJointerByPortRef(portRef dsl.PortRef) *Jointer {
	sourceActivity := l.getActivityById(portRef.NodeId)
	if sourceActivity != nil {
		return sourceActivity.Jointers.GetOutputById(portRef.PortId)
	}
	return nil
}

func (l *LogicFlow) getActivityById(id string) *BaseActivity[any] {
	for i := range l.baseActivites {
		activity := l.baseActivites[i]
		if activity.Id == id {
			return activity
		}
	}

	return nil
}
