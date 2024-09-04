# Temporal API Keys (Data Plane) Demo

## Prerequisites

* A Temporal Cloud account with [API keys enabled](https://docs.temporal.io/cloud/api-keys#manage-api-keys).
* A namespace with API key authentication enabled (*Allow API Key authentication*). \* Note the gRPC endpoint for the namespace.
* [An API key](https://docs.temporal.io/cloud/api-keys) for Namespace authentication.

## Run a simple demo

1. Create a service account that has writer or namespace administrator access to a namespace that has API key authentication enabled.
1. Create [an API key](https://docs.temporal.io/cloud/api-keys) for the service account.
2. Start the worker
```
go run ./worker -namespace <namespace> -grpcEndpoint <grpcEndpoint> -apikey <apikey>
```
3. Run a workflow
```
go run ./starter -namespace <namespace> -grpcEndpoint <grpcEndpoint> -apikey <apikey>
```

## Demonstrate key revocation, hot key reloading 🌶️

1. Create 2 [API keys](https://docs.temporal.io/cloud/api-keys) (`apikey1` and `apikey2`)
2. Start the worker using `apikey1`
```
go run ./worker -namespace <namespace> -grpcEndpoint <grpcEndpoint> -apikey <apikey1>
```
3. Run a workflow using `apikey2`
```
go run ./starter -namespace <namespace> -grpcEndpoint <grpcEndpoint> -apikey <apikey2>
```
4. Disable `key1`
```
tcld apikey disable --id <key1 id>
```
5. Wait for the worker to fail pollings (`WARN  Failed to poll for task.`)
6. Run a workflow using `apikey2`
```
go run ./starter -namespace <namespace> -grpcEndpoint <grpcEndpoint> -apikey <apike21>
```
Note the work flow doesn't progress.
7. Update the worker so that it uses `apikey2` (which is still enabled)
```
curl -X PUT http://localhost:3333/ -d 'apikey2'
```
8. The workflow will now complete.