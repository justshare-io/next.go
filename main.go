package main

//go:generate go run github.com/bufbuild/buf/cmd/buf@latest generate

import (
	"connectrpc.com/connect"
	"connectrpc.com/grpcreflect"
	"context"
	"fmt"
	greetv1 "github.com/justshare-io/nextgo/pkg/gen/proto"
	"github.com/justshare-io/nextgo/pkg/gen/proto/protoconnect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"log"
	"log/slog"
	"net/http"
)

type GreetServer struct{}

// NOTE breadchris this enforces that the proto and go code are in sync
var _ protoconnect.GreetServiceHandler = &GreetServer{}

// NOTE breadchris this is accepts data as json over http _or_ binary protobuf requests
func (s *GreetServer) Greet(
	ctx context.Context,
	req *connect.Request[greetv1.GreetRequest],
) (*connect.Response[greetv1.GreetResponse], error) {
	log.Println("Request headers: ", req.Header())
	res := connect.NewResponse(&greetv1.GreetResponse{
		Greeting: fmt.Sprintf("Hello, %s!", req.Msg.Name),
	})
	res.Header().Set("Greet-Version", "v1")
	return res, nil
}

func main() {
	greeter := &GreetServer{}

	// NOTE breadchris we are starting with just a regular http server
	// that means you can bring all of your existing http middleware
	// here is a more involved example: https://github.com/justshare-io/justshare/blob/main/pkg/server/serve.go#L130
	mux := http.NewServeMux()

	// NOTE breadchris intercept any request and response
	interceptors := connect.WithInterceptors(NewLogInterceptor())

	path, handler := protoconnect.NewGreetServiceHandler(greeter, interceptors)

	mux.Handle(path, handler)

	// NOTE breadchris this is a reflection service that allows
	// you to remotely inspect the service, like graphql introspection.
	mux.Handle(grpcreflect.NewHandlerV1Alpha(
		grpcreflect.NewStaticReflector(protoconnect.GreetServiceName),
	))

	addr := "localhost:8080"
	slog.Info("server started", "addr", addr)
	if err := http.ListenAndServe(
		addr,
		h2c.NewHandler(mux, &http2.Server{}),
	); err != nil {
		slog.Error("failed to start server", "error", err)
	}
}

func NewLogInterceptor() connect.UnaryInterceptorFunc {
	// TODO breadchris support logging for stream calls
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			slog.Info("connect request", "request", fmt.Sprintf("%+v", req))

			resp, err := next(ctx, req)
			if err != nil {
				slog.Error("connect error", "error", fmt.Sprintf("%+v", err))
			}
			return resp, err
		}
	}
}
