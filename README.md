### Steps to run:
1. Setup environment by installing Go, protoc and gRPC tools.
2. Clone this repo
3. `make generate` to generate code
4. `make bump` to get the dependencies in generated code
5. `go mod tidy`
6. `make build` to build the server and client binaries
7. `./bin/server` to run the server binary
8. `./bin/client` to run the client binary

**__Note:__**`localhost:50051` is used by the server and client binaries for gRPC communication. Ensure the port is free before running the binaries