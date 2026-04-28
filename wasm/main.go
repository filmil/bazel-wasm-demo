package main

import (
	"github.com/filmil/bazel-wasm-demo/ui"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

func main() {
	app.Route("/", &ui.Hello{})
	app.RunWhenOnBrowser()
}
