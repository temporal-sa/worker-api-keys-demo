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

type Params struct {
	Namespace    string
	GrpcEndpoint string
	ApiKey       string
	ServerName   string
}

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

func ParseParams(args []string) (Params, error) {
	set := flag.NewFlagSet("worker-api-keys-demo", flag.ExitOnError)

	namespace := set.String("namespace", "Default", "Namespace for the server")
	grpcEndpoint := set.String("grpcEndpoint", "us-east-1.aws.api.temporal.io:7233", "Namespace gRPC endpoint")
	apiKey := set.String("apikey", "", "Data plane API key")
	serverName := strings.Split(*grpcEndpoint, ":")[0]

	if err := set.Parse(args); err != nil {
		return Params{}, fmt.Errorf("failed parsing args: %w", err)
	} else if *apiKey == "" {
		return Params{}, fmt.Errorf("-namespace is required")
	}

	return Params{
		Namespace:    *namespace,
		GrpcEndpoint: *grpcEndpoint,
		ApiKey:       *apiKey,
		ServerName:   serverName,
	}, nil
}

func Connect(params *Params) (client.Client, error) {
	return client.Dial(client.Options{
		HostPort:  params.GrpcEndpoint,
		Namespace: params.Namespace,
		Credentials: client.NewAPIKeyDynamicCredentials(
			func(context.Context) (string, error) {
				return params.ApiKey, nil
			},
		),
		ConnectionOptions: client.ConnectionOptions{
			TLS: &tls.Config{
				InsecureSkipVerify: true,
				ServerName:         params.ServerName,
			},
			DialOptions: []grpc.DialOption{
				grpc.WithUnaryInterceptor(
					func(ctx context.Context, method string, req any, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
						return invoker(
							metadata.AppendToOutgoingContext(ctx, "temporal-namespace", params.Namespace),
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
