package shared

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"strings"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var TASK_QUEUE string = "api-keys-demo"

func HelloWorkflow(ctx workflow.Context, name string) (string, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("HelloWorkflow started", "name", name)

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	})

	var result string

	err := workflow.ExecuteActivity(ctx, HelloActivity, name).Get(ctx, &result)
	if err != nil {
		logger.Error("ExecuteActivity failed", "result", result)
		return "", err
	}

	logger.Info("HelloWorkflow finished", "result", result)

	return result, nil
}

func HelloActivity(ctx context.Context, name string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Activity", "name", name)
	return "Hello, " + name, nil
}

func Connect(args []string) (client.Client, error) {
	set := flag.NewFlagSet("worker-api-keys-demo", flag.ExitOnError)

	namespace := set.String("namespace", "Default", "Namespace for the server")
	grpcEndpoint := set.String("grpcEndpoint", "us-east-1.aws.api.temporal.io:7233", "Namespace gRPC endpoint")
	apiKey := set.String("apikey", "", "Data plane API key")
	serverName := strings.Split(*grpcEndpoint, ":")[0]

	if err := set.Parse(args); err != nil {
		return nil, fmt.Errorf("failed parsing args: %w", err)
	} else if *apiKey == "" {
		return nil, fmt.Errorf("-namespace is required")
	}

	return client.Dial(client.Options{
		HostPort:    *grpcEndpoint,
		Namespace:   *namespace,
		Credentials: client.NewAPIKeyStaticCredentials(*apiKey),
		ConnectionOptions: client.ConnectionOptions{
			TLS: &tls.Config{
				InsecureSkipVerify: true,
				ServerName:         serverName,
			},
			DialOptions: []grpc.DialOption{
				grpc.WithUnaryInterceptor(
					func(ctx context.Context, method string, req any, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
						return invoker(
							metadata.AppendToOutgoingContext(ctx, "temporal-namespace", *namespace),
							method,
							req,
							reply,
							cc,
							opts...,
						)
					},
				),
			},
		},
	})
}
