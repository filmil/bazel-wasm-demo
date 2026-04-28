package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/filmil/bazel-wasm-demo/protos/api"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *api.HelloRequest) (*api.HelloReply, error) {
	name := in.Name
	if name == "" {
		name = "World"
	}
	return &api.HelloReply{Message: "Hello, " + name}, nil
}

func main() {
	appHandler := &app.Handler{
		Name:        "Hello WASM",
		Description: "A simple Hello World WASM app",
	}

	mux := http.NewServeMux()

	api.RegisterGreeterHTTPMux(mux, &server{})
	mux.Handle("/", appHandler)

	fmt.Println("Listening on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
