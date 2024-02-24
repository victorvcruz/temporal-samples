package workflow

import (
	"fmt"
	"go.temporal.io/sdk/workflow"
)

func ExampleWorkflow(ctx workflow.Context) error {
	fmt.Println("example workflow started")
	return nil
}
