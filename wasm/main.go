package main

import (
	"github.com/filmil/bazel-wasm-demo/ui"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

func main() {
	app.Route("/", func() app.Composer { return &ui.Hello{} })
	app.RunWhenOnBrowser()
}
