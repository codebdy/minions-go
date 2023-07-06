package runtime

import "github.com/codebdy/minions-go/dsl"

type LogicFlow struct {
	Id       string
	Jointers *ActivityJointers
	flowMeta *dsl.LogicFlowDefine
	activities [] *any
}

func NewLogicflow(flowMeta dsl.LogicFlowDefine, logicflowContext any) *LogicFlow {
	var logicFlow LogicFlow
	logicFlow.Id = flowMeta.Id
	logicFlow.Jointers = &ActivityJointers{}
	logicFlow.flowMeta = &flowMeta
	//第一步，解析节点
	logicFlow.constructActivities()
	//第二步， 构建连接关系
	logicFlow.contructLines()
	return &logicFlow
}

//构建一个图的所有节点
func (l *LogicFlow) constructActivities() {
	for _, activityMeta := range l.flowMeta.Nodes {
		switch (activityMeta.Type) {
		case dsl.ACTIVITY_TYPE_START:
			name := "input"
			if activityMeta.activityName != "" {
				name = activityMeta.activityName
			}
			//start只有一个端口，所以name可以跟meta name一样
			l.Jointers.inputs = append(l.Jointers.inputs, Jointer{Id:activityMeta.id, name:name})
		case dsl.ACTIVITY_TYPE_END:
			name := "output"
			if activityMeta.activityName != "" {
				name = activityMeta.activityName
			}
			//end 只有一个端口，所以name可以跟meta name一样
			l.Jointers.outputs = append(l.Jointers.outputs, Jointer{Id:activityMeta.id, name:name})
		case dsl.ACTIVITY_TYPE_ACTIVITY, dsl.ACTIVITY_TYPE_LOGICFLOWACTIVITY:
			if activityMeta.activityName != "" {
				activityType := activitiesMap[activityMeta.activityName]
				if (activityType == nil) {
					panic("Can not find activity by name:" + activityMeta.activityName)
				}
				if activityType.Kind() != reflect.Struct {
					panic("expect struct")
				}

				activity := reflect.New(activityType) 
				//new activityClass(activityMeta, this.context);

				//构造Jointers
				for (const out of activityMeta.outPorts || []) {
					activity.jointers.outputs.push(new Jointer(out.id, out.name))
				}
				for (const input of activityMeta.inPorts || []) {
					activity.jointers.inputs.push(new Jointer(input.id, input.name))
				}

				//把input端口跟处理函数相连
				for (const inputName of Object.keys(activityInfo.methodMap)) {
					const handleName = activityInfo.methodMap[inputName]
					// eslint-disable-next-line @typescript-eslint/no-explicit-any
					const handle = (activity as any)?.[handleName];
					const handleWithThis = handle?.bind(activity);
					handleName && activity.jointers.getInput(inputName)?.connect(handleWithThis)
				}

				//处理动态端口
				if (activityInfo.dynamicMethod) {
					// eslint-disable-next-line @typescript-eslint/no-explicit-any
					const handle = (activity as any)?.[activityInfo.dynamicMethod];
					const handleWithThis = handle?.bind(activity);

					for (const input of activity.jointers.inputs) {
						const handeWraper = (inputValue: unknown) => {
							return handleWithThis?.(input.name, inputValue)
						}
						input.connect(handeWraper)
					}
				}

				this.activities.push(activity)
			}
	}
	}
}

//连接一个图的所有节点，把所有的jointer连起来
func (l *LogicFlow) contructLines() {

}
