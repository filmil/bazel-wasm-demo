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
	prefix := app.Getenv("GOAPP_PROXY_PATH")
	h.client = api.NewGreeterHTTPClient(prefix)
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
	return app.Div().Class("container mt-5").Body(
		app.Div().Class("row justify-content-center").Body(
			app.Div().Class("col-md-6").Body(
				app.Div().Class("card shadow").Body(
					app.Div().Class("card-header bg-primary text-white").Body(
						app.H3().Class("card-title mb-0").Text("WASM gRPC Demo"),
					),
					app.Div().Class("card-body").Body(
						app.Div().Class("mb-3").Body(
							app.Label().Class("form-label").Text("Enter your name:"),
							app.Input().
								Class("form-control").
								Type("text").
								Value(h.name).
								Placeholder("Name").
								OnChange(h.ValueTo(&h.name)),
						),
						app.Button().
							Class("btn btn-primary w-100").
							Text("Say Hello").
							OnClick(h.getGreeting),
						app.If(h.message != "", func() app.UI {
							return app.Div().Class("alert alert-success mt-3").Body(
								app.Text(h.message),
							)
						}),
					),
				),
			),
		),
	)
}
