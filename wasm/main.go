package main

import (
	"fmt"
	"github.com/filmil/bazel-wasm-demo/ui"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

func main() {
	prefix := app.Getenv("GOAPP_PROXY_PATH")
	fmt.Printf("WASM starting with prefix: %s\n", prefix)
	
	// Register the root route.
	app.Route("/", func() app.Composer { return &ui.Hello{} })
	
	// If a prefix is set, also register it so that reloading works.
	if prefix != "" && prefix != "/" {
		app.Route(prefix, func() app.Composer { return &ui.Hello{} })
		app.Route(prefix+"/", func() app.Composer { return &ui.Hello{} })
	}
	
	app.RunWhenOnBrowser()
}
