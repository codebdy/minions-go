package runtime

import (
	"context"
	"reflect"
	"strings"

	"github.com/codebdy/minions-go/dsl"
)

type LogicFlow struct {
	//Id       string
	Jointers *ActivityJointers
	flowMeta *dsl.LogicFlowMeta
	//存Activity 指针
	baseActivites []*BaseActivity

	ctx context.Context
}

func AttachSubFlowsToContext(flowMetas *[]dsl.SubLogicFlowMeta, ctx context.Context) context.Context {
	return context.WithValue(ctx, CONTEXT_KEY_SUBMETAS, flowMetas)
}

// ctx用于传递value，minions.CONTEXT_KEY_SUBMETAS 对应*[]dsl.LogicFlowDefine， 子编排metas
func NewLogicflow(flowMeta dsl.LogicFlowMeta, ctx context.Context) *LogicFlow {
	var logicFlow LogicFlow
	//logicFlow.Id = flowMeta.Id
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
		case dsl.ACTIVITY_TYPE_ACTIVITY, dsl.ACTIVITY_TYPE_LOGICFLOWACTIVITY, dsl.ACTIVITY_TYPE_EMBEDDEDFLOW:
			l.newActivity(activityMeta)
		}
	}
}

func (l *LogicFlow) newActivity(activityMeta dsl.ActivityDefine) {
	if activityMeta.ActivityName != "" {
		activityType := activitiesMap[activityMeta.ActivityName]
		if activityType == nil {
			panic("Can not find activity by name:" + activityMeta.ActivityName)
		}
		if activityType.Kind() != reflect.Struct {
			panic("expect struct")
		}

		activity := reflect.New(activityType).Interface()

		if activityMeta.Type == dsl.ACTIVITY_TYPE_EMBEDDEDFLOW {
			//重新构造子节点，主要目的：把父节点端口转换成子流程的开始节点跟结束节点
			activityMeta = refactorChildren(activityMeta)
		}

		activityValue := reflect.ValueOf(activity)
		parentActivityValue := reflect.Indirect(activityValue).FieldByName("Activity")
		if !parentActivityValue.IsValid() {
			panic("Activity has not parent Activity")
		}
		baseActivityValue := reflect.Indirect(parentActivityValue).FieldByName("BaseActivity")
		if baseActivityValue.IsValid() {
			v := baseActivityValue.Addr().Interface().(*BaseActivity)
			v.Init(&activityMeta, l.ctx)
			//构造Jointers
			for _, out := range activityMeta.OutPorts {
				v.Jointers.outputs = append(v.Jointers.outputs, &Jointer{Id: out.Id, Name: out.Name})
			}

			for _, input := range activityMeta.InPorts {
				v.Jointers.inputs = append(v.Jointers.inputs, &Jointer{Id: input.Id, Name: input.Name})
			}

			//调用具体活动的初始化
			m := activityValue.MethodByName("Init")
			if m.IsValid() {
				inputs := make([]reflect.Value, 0)
				m.Call(inputs)
			}

			for i := range v.Jointers.inputs {
				input := v.Jointers.inputs[i]

				m := activityValue.MethodByName(FirstUpper(input.Name))
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
							if inputValue == nil {
								//防止nil参数异常：Call using zero Value argument
								inputs[0] = reflect.ValueOf((*string)(nil))
							} else {
								inputs[0] = reflect.ValueOf(inputValue)
							}
						}
						if mt.NumIn() > 1 {
							inputs[1] = reflect.ValueOf(ctx)
						}
						m.Call(inputs)
					})
				} else {
					//处理动态端口
					m = activityValue.MethodByName("DynamicInput")
					if m.IsValid() {
						mt := m.Type()
						input.Connect(func(inputValue any, ctx context.Context) {
							inputs := make([]reflect.Value, mt.NumIn())
							if mt.NumIn() == 0 {
								panic("DynamicInput must have one args")
							}
							inputs[0] = reflect.ValueOf(input.Name)
							if mt.NumIn() > 1 {
								if inputValue == nil {
									//防止nil参数异常：Call using zero Value argument
									inputs[1] = reflect.ValueOf((*string)(nil))
								} else {
									inputs[1] = reflect.ValueOf(inputValue)
								}
							}
							if mt.NumIn() > 2 {
								inputs[2] = reflect.ValueOf(ctx)
							}
							m.Call(inputs)
						})
						//如果不是子流程，必须要定义端口处理函数
					} else if activityMeta.Type != dsl.ACTIVITY_TYPE_LOGICFLOWACTIVITY {
						panic("Can not find input handler:" + input.Name)
					}
				}
			}

			l.baseActivites = append(l.baseActivites, v)
		} else {
			panic("Activity has not BaseActivity")
		}
	}
}
func caseInsenstiveMethodByName(v reflect.Value, name string) reflect.Value {
	name = strings.ToLower(name)
	return v.FieldByNameFunc(func(n string) bool { return strings.ToLower(n) == name })
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

		//如果不是，连接Target activity的Input
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

func (l *LogicFlow) getActivityById(id string) *BaseActivity {
	for i := range l.baseActivites {
		activity := l.baseActivites[i]
		if activity.Id == id {
			return activity
		}
	}

	return nil
}

//重新构造children，添加边界节点，修改连线
func refactorChildren(parentMeta dsl.ActivityDefine) dsl.ActivityDefine {
	//对象被复制
	newMeta := parentMeta

	//清空连线，连线需要重新被构建，节点只是新增，可以保留旧的
	newMeta.Children.Lines = []dsl.LineDefine{}

	//父节点的input创建为start, portId=>start 节点 id
	for _, input := range parentMeta.InPorts {
		newMeta.Children.Nodes = append(newMeta.Children.Nodes, dsl.ActivityDefine{
			Id:           input.Id,
			Type:         dsl.ACTIVITY_TYPE_START,
			ActivityName: "start",
			Name:         input.Name,
		})
	}

	//父节点的output创建为end, portId=>end 节点 id
	for _, output := range parentMeta.OutPorts {
		newMeta.Children.Nodes = append(newMeta.Children.Nodes, dsl.ActivityDefine{
			Id:           output.Id,
			Type:         dsl.ACTIVITY_TYPE_END,
			ActivityName: "end",
			Name:         output.Name,
		})
	}

	for _, line := range parentMeta.Children.Lines {
		//复制连线
		newLine := line
		//起点是父节点输入端口， 连接到新创建的开始节点
		if line.Source.NodeId == parentMeta.Id && line.Source.PortId != "" {
			newLine.Source.NodeId = line.Source.PortId
			newLine.Source.PortId = ""
		}
		//终点是父节点输入端口, 连接到新创建的结束节点
		if line.Target.NodeId == parentMeta.Id && line.Target.PortId != "" {
			newLine.Target.NodeId = line.Target.PortId
			newLine.Target.PortId = ""
		}
		newMeta.Children.Lines = append(newMeta.Children.Lines, newLine)
	}

	return newMeta
}
