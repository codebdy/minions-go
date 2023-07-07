package runtime

import (
	"context"
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
			},
		},
		Lines: []dsl.LineDefine{},
	}

	NewLogicflow(logicFlowMetas, context.Background())
}
