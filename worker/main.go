package main

import (
	"io"
	"log"
	"net/http"
	"os"

	shared "github.com/temporal-sa/worker-api-keys-demo"
	"go.temporal.io/sdk/worker"
)

func main() {
	params, err := shared.ParseParams(os.Args[1:])
	if err != nil {
		log.Fatalln("Failed to parse input parameters", err)
	}

	// "kms" service
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		params.ApiKey = string(body)
		log.Default().Println("API key updated")

		w.WriteHeader(http.StatusAccepted)
	})

	go func() {
		err := http.ListenAndServe(":3333", nil)
		if err != nil {
			log.Fatalln("Unable to start webserver", err)
		}
	}()

	c, err := shared.Connect(&params)
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
