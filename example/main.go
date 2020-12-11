package main

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	csv "github.com/matoubidou/grpc-gateway-csv"
	"google.golang.org/grpc"
)

//go:generate protoc -I . -I $HOME/go/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out . --go_opt paths=source_relative --go-grpc_out . --grpc-gateway_out . --go-grpc_opt paths=source_relative example-service.proto

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

func (s Server) mustEmbedUnimplementedExampleServiceServer() {

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
