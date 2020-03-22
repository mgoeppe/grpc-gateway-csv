package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/mathias-zeller/grpc-gateway-csv"
	"google.golang.org/grpc"
	"net"
	"net/http"
)

//go:generate protoc -I/usr/local/include -I. -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:. example-service.proto
//go:generate protoc -I/usr/local/include -I. -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:. example-service.proto

type Server struct{}

func (s Server) Example(context.Context, *ExampleRequest) (*ExampleResponse, error) {
	return &ExampleResponse{
		Slice: []*Outer{
			{
				Col1: "dreggn",
				Col2: 42,
				Inner: &Inner{
					Col3: true,
					Col4: []string{"one", "two"},
					Col5: map[string]string{"k1": "v1"},
				},
			},
			{
				Col1: "dreggn",
				Col2: 42,
				Inner: &Inner{
					Col3: true,
					Col4: []string{"eins", "zwo"},
					Col5: map[string]string{"k1": "v1"},
				},
			},
		},
	}, nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux(runtime.WithMarshalerOption("text/csv", &csv.Marshaler{}))
	srv := grpc.NewServer()
	RegisterExampleServiceServer(srv, Server{})
	err := RegisterExampleServiceHandlerFromEndpoint(ctx, mux, ":8080", []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		panic(err)
	}

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	fmt.Println("grpc on :8080 ..")
	go func() {
		if err := srv.Serve(l); err != nil {
			panic(err)
		}
	}()

	fmt.Println("http on :8081 ..")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		panic(err)
	}
}
