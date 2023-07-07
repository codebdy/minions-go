package runtime

import (
	"context"
	"fmt"
	"testing"

	"github.com/codebdy/minions-go/dsl"
)

type TestActivity struct {
	BaseActivity BaseActivity[any]
}

func init() {
	RegisterActivity(
		"test1",
		TestActivity{},
	)
}

func (t TestActivity) Input(inputValue any) {
	fmt.Println("Test Input ", inputValue)
}

func TestNewLogicflow(t *testing.T) {
	logicFlowMetas := dsl.LogicFlowDefine{
		Id:    "test1",
		Name:  "Test 1",
		Label: "测试",
		Nodes: []dsl.ActivityDefine{
			{
				Id:           "start_id_1",
				ActivityName: "input1",
				Type:         dsl.ACTIVITY_TYPE_START,
			},
			{
				Id:           "test_id_1",
				ActivityName: "test1",
				Type:         dsl.ACTIVITY_TYPE_ACTIVITY,
				InPorts: []dsl.PortDefine{
					{
						Id:   "test_1_input1",
						Name: "input",
					},
				},
			},
		},
		Lines: []dsl.LineDefine{
			{
				Id: "line1",
				Source: dsl.PortRef{
					NodeId: "start_id_1",
				},
				Target: dsl.PortRef{
					NodeId: "test_id_1",
					PortId: "test_1_input1",
				},
			},
		},
	}

	logicFolow := NewLogicflow(logicFlowMetas, context.Background())

	logicFolow.Jointers.GetInput("input1").Push("From Test", context.Background())
}
