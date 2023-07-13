package runtime

import (
	"context"
	"testing"

	"github.com/codebdy/minions-go/dsl"
)

var testInputValue any

var inited bool

type TestActivity struct {
	BaseActivity BaseActivity
}

func init() {
	RegisterActivity(
		"test1",
		TestActivity{},
	)
}

func (t TestActivity) Input(inputValue any) {
	testInputValue = inputValue
}

func (t TestActivity) Init() {
	inited = true
}

func TestNewLogicflow(t *testing.T) {
	logicFlowMetas := dsl.LogicFlowMeta{
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

	if testInputValue != "From Test" {
		t.Error("Can not pass input value")
	}

	if !inited {
		t.Error("Init is not invoked")
	}

}
