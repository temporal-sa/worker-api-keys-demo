# NOTE: this isn't a script to run; it's a cheatsheet if you want to setup/run
#       the demo from the CLI.

# SETUP

# Create a namespace with API key auth enabled
NAMESPACE_SHORTNAME=apikey-test-2
REGION=us-east-1

NAMESPACE=`tcld namespace create --namespace $NAMESPACE_SHORTNAME --region $REGION --auth-method 'api_key' | jq -r '.requestStatus.resourceId'`
# note: wait for the namespace creation request to be fulfilled
GRPC_ENDPOINT=`tcld namespace get --namespace $NAMESPACE | jq -r '.uri.regionalGrpc'`

# Create a service account with admin permission for the namespace
SA_ID=`tcld service-account create --name from-cli --account-role Read --namespace-permission $NAMESPACE=Admin | jq -r '.serviceAccountId'`

# Create 2 API keys
KEY_NAME=demo-key

# You'll probably have to echo these out to transfer them between terminals for the worker and client
KEY1=`tcld apikey create --name $KEY_NAME-1 --duration 30d --service-account-id $SA_ID | jq -r '.secretKey'`
KEY2=`tcld apikey create --name $KEY_NAME-2 --duration 30d --service-account-id $SA_ID | jq -r '.secretKey'`
KEY1_ID=`tcld apikey l | jq -r '.apiKeys[] | select(.spec.displayName == "'$KEY_NAME'-1").id'`


# SIMPLE DEMO

# Run a basic 'hello, world!' demo
go run ./worker -namespace $NAMESPACE -grpcEndpoint $GRPC_ENDPOINT -apikey $KEY1
go run ./starter -namespace $NAMESPACE -grpcEndpoint $GRPC_ENDPOINT -apikey $KEY2


# KEY REVOCATION DEMO

# Disable key 1 (the worker key)
tcld apikey disable --id $KEY1_ID

# This workflow won't progress (it may take a moment for the disable request to complete)
go run ./starter -namespace $NAMESPACE -grpcEndpoint $GRPC_ENDPOINT -apikey $KEY2

# Enable key 1 (the worker key) - workflow will complete
tcld apikey enable --id $KEY1_ID


# KEY HOT RELOADING DEMO

# Disable key 1 (the worker key)
tcld apikey disable --id $KEY1_ID

# This workflow won't progress (it may take a moment for the disable request to complete)
go run ./starter -namespace $NAMESPACE -grpcEndpoint $GRPC_ENDPOINT -apikey $KEY2

# Update the key used by the worker - workflow will complete
curl -X PUT http://localhost:3333 -d $KEY2