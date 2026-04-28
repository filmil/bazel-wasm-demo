package main

import (
	"github.com/filmil/bazel-wasm-demo/protos/api"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"log"
)

type hello struct {
	app.Compo
	name    string
	message string
	client  api.GreeterHTTPClientInterface
}

func (h *hello) OnMount(ctx app.Context) {
	h.client = api.NewGreeterHTTPClient("")
}

func (h *hello) getGreeting(ctx app.Context, e app.Event) {
	ctx.Async(func() {
		reply, err := h.client.SayHello(ctx, &api.HelloRequest{Name: h.name})
		if err != nil {
			log.Println("Error:", err)
			return
		}
		ctx.Dispatch(func(ctx app.Context) {
			h.message = reply.Message
		})
	})
}

func (h *hello) Render() app.UI {
	return app.Div().Body(
		app.H1().Text("WASM gRPC Demo"),
		app.P().Body(
			app.Input().
				Type("text").
				Value(h.name).
				Placeholder("Enter your name").
				OnChange(h.ValueTo(&h.name)),
			app.Button().
				Text("Say Hello").
				OnClick(h.getGreeting),
		),
		app.P().Text(h.message),
	)
}

func main() {
	app.Route("/", &hello{})
	app.RunWhenOnBrowser()
}
