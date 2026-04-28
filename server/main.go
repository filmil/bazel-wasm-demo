package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/filmil/bazel-wasm-demo/protos/api"
	"github.com/filmil/bazel-wasm-demo/ui"
	"github.com/golang/glog"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

//go:embed web/*
var webAssets embed.FS

var (
	port = flag.Int("port", 8080, "default port to use")
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

	app.Route("/", func() app.Composer { return &ui.Hello{} })

	appHandler := &app.Handler{
		Name:        "Hello WASM",
		Description: "A simple Hello World WASM app",
		Icon: app.Icon{
			Default: "/web/icon.png",
		},
	}

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.FS(webAssets))
	mux.HandleFunc("/web/", func(w http.ResponseWriter, r *http.Request) {
		glog.Infof("Serving static web asset: %s", r.URL.Path)
		fileServer.ServeHTTP(w, r)
	})

	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		glog.Infof("Serving file: /favicon.ico")
		file, err := webAssets.ReadFile("web/favicon.ico")
		if err != nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "image/x-icon")
		w.Write(file)
	})

	api.RegisterGreeterHTTPMux(mux, &server{})
	mux.Handle("/", appHandler)

	glog.Infof("Listening on http://localhost:%v", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%v", *port), mux); err != nil {
		glog.Errorf("failed to serve: %v", err)
		os.Exit(1)
	}
}
