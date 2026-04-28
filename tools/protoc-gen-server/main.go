package main

import (
	"bytes"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	protogen.Options{}.Run(func(gen *protogen.Plugin) error {
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			generateServer(gen, f)
		}
		return nil
	})
}

const serviceTpl = `package {{ .GoPackageName }}

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

{{ range .Services }}
type {{ .GoName }}HTTPServerInterface interface {
	{{ range .Methods }}
	{{ .GoName }}(ctx context.Context, in *{{ .Input.GoIdent.GoName }}) (*{{ .Output.GoIdent.GoName }}, error)
	{{ end }}
}

func Register{{ .GoName }}HTTPMux(mux *http.ServeMux, srv {{ .GoName }}HTTPServerInterface) {
	{{ $serviceName := .GoName }}
	{{ $protoPackageName := .Desc.ParentFile.Package }}

	{{ range .Methods }}
	{{ $methodPath := printf "/%s.%s/%s" $protoPackageName $serviceName .GoName }}
	mux.HandleFunc("{{ $methodPath }}", func(w http.ResponseWriter, r *http.Request) {
		in := new({{ .Input.GoIdent.GoName }})
		inJSON, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer r.Body.Close()

		if len(inJSON) > 0 {
			err = json.Unmarshal(inJSON, in)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
		}

		ret, err := srv.{{ .GoName }}(context.Background(), in)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		retJSON, err := json.Marshal(ret)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(retJSON)
	})
	{{ end }}
}
{{ end }}
`

func generateServer(gen *protogen.Plugin, file *protogen.File) {
	if len(file.Services) == 0 {
		return
	}
	filename := file.GeneratedFilenamePrefix + ".server.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)

	tmpl, err := template.New("server").Parse(serviceTpl)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, file)
	if err != nil {
		panic(err)
	}
	g.P(buf.String())
}
