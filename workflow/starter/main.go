package main

import (
	"context"
	"fmt"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"log"
	sample "temporal-samples/workflow"
)

func main() {
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		panic(err)
	}

	w := worker.New(c, "taskQueue", worker.Options{})

	w.RegisterWorkflow(sample.ExampleWorkflow)

	we, err := c.ExecuteWorkflow(context.Background(), client.StartWorkflowOptions{
		TaskQueue: "taskQueue",
	}, sample.ExampleWorkflow)
	if err != nil {
		panic(err)
	}
	fmt.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start Worker", err)
	}
}
