package workflow

import (
	"context"
	"fmt"
	"go.temporal.io/sdk/workflow"
	"time"
)

func ExampleWorkflow(ctx workflow.Context) error {
	fmt.Println("example workflow started")
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 10,
	})

	err := workflow.ExecuteActivity(ctx, ExampleActivity).Get(ctx, nil)
	if err != nil {
		return err
	}

	err = workflow.ExecuteActivity(ctx, ExampleActivity).Get(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}

func ExampleActivity(ctx context.Context) error {
	fmt.Println("example activity started")
	return nil
}
