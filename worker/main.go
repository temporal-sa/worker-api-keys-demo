package main

import (
	"log"
	"os"

	shared "github.com/temporal-sa/worker-api-keys-demo"
	"go.temporal.io/sdk/worker"
)

func main() {
	c, err := shared.Connect(os.Args[1:])
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	w := worker.New(c, shared.TASK_QUEUE, worker.Options{})

	w.RegisterWorkflow(shared.HelloWorkflow)
	w.RegisterActivity(shared.HelloActivity)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start Worker", err)
	}
}
