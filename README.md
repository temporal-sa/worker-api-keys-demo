# Temporal API Keys (Data Plane) Demo

## Prerequisites

* A Temporal Cloud account with [API keys enabled](https://docs.temporal.io/cloud/api-keys#manage-api-keys).
* A namespace with API key authentication enabled (*Allow API Key authentication*). \* Note the gRPC endpoint for the namespace.
* [An API key](https://docs.temporal.io/cloud/api-keys) for Namespace authentication.

## Run this demo

1. Create [an API key](https://docs.temporal.io/cloud/api-keys) for Namespace authentication (see [prequisites](#prerequisites)).
2. Create a name
3. Start the worker
```
go run ./worker -namespace <namespace> -grpcEndpoint <grpcEndpoint> -apikey <apikey>
```
4. Run a workflow
```
go run ./starter -namespace <namespace> -grpcEndpoint <grpcEndpoint> -apikey <apikey>
```