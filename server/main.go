package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/filmil/bazel-wasm-demo/protos/api"
	"github.com/filmil/bazel-wasm-demo/ui"
	"github.com/golang/glog"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

//go:embed web/*
var webAssets embed.FS

var (
	port = flag.Int("port", 8080, "default port to use")
	host = flag.String("host", "", "default host to use")
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *api.HelloRequest) (*api.HelloReply, error) {
	name := in.Name
	if name == "" {
		name = "World"
	}
	glog.Infof("SayHello called with name: %s", name)
	return &api.HelloReply{Message: "Hello, " + name}, nil
}

func main() {
	flag.Parse()

	// Register the root route.
	app.Route("/", func() app.Composer { return &ui.Hello{} })

	appHandler := &app.Handler{
		Name:        "Hello WASM",
		Description: "A simple Hello World WASM app",
		Styles:      []string{"/web/bootstrap.min.css"},
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

	// Use a dynamic app handler that responds to the prefix.
	dynamicAppHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prefix := r.Header.Get("X-Forwarded-Prefix")
		if prefix == "" {
			appHandler.ServeHTTP(w, r)
			return
		}
		prefix = "/" + strings.Trim(prefix, "/")
		hCopy := *appHandler
		hCopy.Resources = app.PrefixedLocation(prefix)
		if hCopy.Env == nil {
			hCopy.Env = make(map[string]string)
		}
		hCopy.Env["GOAPP_PROXY_PATH"] = prefix
		hCopy.ServeHTTP(w, r)
	})
	mux.Handle("/", dynamicAppHandler)

	glog.Infof("Listening on http://%v:%v", *host, *port)

	// The global proxy handler handles path stripping.
	globalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prefix := r.Header.Get("X-Forwarded-Prefix")
		if prefix != "" {
			prefix = "/" + strings.Trim(prefix, "/")
			if strings.HasPrefix(r.URL.Path, prefix) {
				http.StripPrefix(prefix, mux).ServeHTTP(w, r)
				return
			}
		}
		mux.ServeHTTP(w, r)
	})

	if err := http.ListenAndServe(fmt.Sprintf("%v:%v", *host, *port), globalHandler); err != nil {
		glog.Errorf("failed to serve: %v", err)
		os.Exit(1)
	}
}
