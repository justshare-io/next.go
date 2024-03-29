# nextgo
your site, built now.

## Try it out
```shell
# hit the http api server
go run github.com/justshare-io/nextgo@latest
curl --header "Content-Type: application/json" --data '{"name": "Jane"}' http://localhost:8080/greet.GreetService/Greet 

# ui for calling the grpc server, like postman.
go run github.com/fullstorydev/grpcui/cmd/grpcui@latest -plaintext localhost:8080
```

## Learn

### Code
- server: I would start by reading the comments in [main.go](main.go). 
- proto: [greet.proto](greet/greet.proto) is the protobuf file that defines the service.

### Proto
If you change the proto file, you need to regenerate the go code. 
```shell
go install github.com/bufbuild/buf/cmd/buf@latest
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest

go generate ./...
```

### GPRC
```shell
go install github.com/fullstorydev/grpcui/cmd/grpcui@latest
grpcui -plaintext localhost:8080
```

### Dependencies
- [go](https://golang.org/doc/)
- [protobuf](https://developers.google.com/protocol-buffers)
- [grpc](https://grpc.io/docs/what-is-grpc/introduction/)

## Goals
- Works today, tomorrow, and next year.
- Scales with both your infra and team.
- Easy to learn and fun to use.

### Recommended Reading
- [Twelve-Factor App](https://12factor.net/)
- [The Go Programming Language](https://www.gopl.io/)
- [Hypermedia Systems](https://hypermedia.systems/)

## Stack
- go
  - grpc
