# gRPC-client-server-based-calculator
Calculator server with some basic functionality implemented using gRPC.


## Execution and Compilation
### To update and regenerate proto dependencies
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    messages_proto/messages.proto