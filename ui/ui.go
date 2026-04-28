package ui

import (
	"log"

	"github.com/filmil/bazel-wasm-demo/protos/api"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type Hello struct {
	app.Compo
	name    string
	message string
	client  api.GreeterHTTPClientInterface
}

func (h *Hello) OnMount(ctx app.Context) {
	h.client = api.NewGreeterHTTPClient("")
}

func (h *Hello) getGreeting(ctx app.Context, e app.Event) {
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

func (h *Hello) Render() app.UI {
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
