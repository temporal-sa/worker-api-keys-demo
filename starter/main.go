package main

import (
	"context"
	"log"
	"os"

	shared "github.com/temporal-sa/worker-api-keys-demo"
	"go.temporal.io/sdk/client"
)

func main() {
	c, err := shared.Connect(os.Args[1:])
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: shared.TASK_QUEUE,
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, shared.HelloWorkflow, "Temporal")
	if err != nil {
		log.Fatalln("Failed to execute workflow", err)
	}

	log.Println("Workflow started", "id", we.GetID(), "runId", we.GetRunID())

	var result string
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Failed to get workflow result", err)
	}

	log.Println("Workflow result: ", result)
}
