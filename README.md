# gRPC Server Example

## Quick Start

Start server using `go run main.go`. To print available arguments, run `go run main.go -h`. If using gateway mode or gateway-hybrid mode, you can play with APIs at http://localhost:8080/swagger.

## Development

### Prerequisites

1. [Go](https://go.dev/dl)
2. [Protocol buffer compiler](https://grpc.io/docs/protoc-installation), version 3
3. Go plugins for the protocol buffer compiler
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```
4. gRPC gateway and OpenAPI plugin for the protocol buffer compiler
```
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
```

### Generate Server Library

Run the example command in the "grpc_example_proto" repository to generate server library.
```
OUTPUT_DIR=./go_gens
TARGET_DIR=../grpc_example_server
make OUTPUT_DIR=${OUTPUT_DIR} \
  EXTRA_FLAGS="--grpc-gateway_opt=logtostderr=true,generate_unbound_methods=true --grpc-gateway_out=${OUTPUT_DIR} --openapiv2_opt=logtostderr=true,generate_unbound_methods=true,allow_merge=true --openapiv2_out=${OUTPUT_DIR}" \
  && cp -r go_gens/grpc_example/proto ${TARGET_DIR} \
  && cp go_gens/apidocs.swagger.json ${TARGET_DIR}/third_party/swagger_ui
```

### Build

1. Build server binary for the Linux AMD64 platform
```
env GOOS=linux GOARCH=amd64 go build -o build/server main.go
```
2. Build Docker image
```
docker build -t grpc_example_server .
```

### Testing

- Run all tests
```
go test github.com/zmzhang8/grpc_example/...
```
- Run tests in a package
```
go test github.com/zmzhang8/grpc_example/path/to/package
```

## References

- [Go Testing](https://pkg.go.dev/testing)
- [Protocol Buffers in Go](https://developers.google.com/protocol-buffers/docs/reference/go-generated)
- [gRPC in Go](https://grpc.io/docs/languages/go)
- [gRPC Metadata in Go](https://github.com/grpc/grpc-go/blob/master/Documentation/grpc-metadata.md)
- [gRPC Error Handling](https://www.grpc.io/docs/guides/error)
- [gRPC Performance Best Practices](https://www.grpc.io/docs/guides/performance)
- [gRPC Health Checking Protocol](https://github.com/grpc/grpc/blob/master/doc/health-checking.md)
- [gRPC Server Reflection](https://github.com/grpc/grpc/blob/master/doc/server-reflection.md)
- [gRPC Web](https://grpc.io/docs/platforms/web/basics)
- [The state of gRPC in the browser](https://grpc.io/blog/state-of-grpc-web)
- [gRPC Gateway](https://github.com/grpc-ecosystem/grpc-gateway)
- [Swagger UI](https://github.com/swagger-api/swagger-ui)