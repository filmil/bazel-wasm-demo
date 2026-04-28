package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/filmil/bazel-wasm-demo/protos/api"
	"github.com/filmil/bazel-wasm-demo/ui"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

var (
	wasmLoc    = flag.String("wasm-path", "", "path to the web app wasm file")
	iconLoc    = flag.String("icon-path", "", "path to the icon")
	faviconLoc = flag.String("favicon-path", "", "path to the favicon")
	port       = flag.Int("port", 8080, "default port to use")
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
	flag.Parse()
	if *wasmLoc == "" {
		log.Fatalf("The flag --wasm-path is required.")
	}

	app.Route("/", func() app.Composer { return &ui.Hello{} })

	appHandler := &app.Handler{
		Name:        "Hello WASM",
		Description: "A simple Hello World WASM app",
		Icon: app.Icon{
			Default: "/web/icon.png",
		},
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/web/app.wasm", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, *wasmLoc)
	})
	if *iconLoc != "" {
		mux.HandleFunc("/web/icon.png", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, *iconLoc)
		})
	}
	if *faviconLoc != "" {
		mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, *faviconLoc)
		})
	}

	api.RegisterGreeterHTTPMux(mux, &server{})
	mux.Handle("/", appHandler)

	log.Printf("Listening on http://localhost:%v\n", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%v", *port), mux); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

