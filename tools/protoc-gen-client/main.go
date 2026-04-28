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
			generateClient(gen, f)
		}
		return nil
	})
}

const clientTpl = `package {{ .GoPackageName }}

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

{{ range .Services }}
type {{ .GoName }}HTTPClientInterface interface {
	{{ range .Methods }}
	{{ .GoName }}(ctx context.Context, in *{{ .Input.GoIdent.GoName }}) (*{{ .Output.GoIdent.GoName }}, error)
	{{ end }}
}

type {{ .GoName }}HTTPClient struct {
	url string
}

func New{{ .GoName }}HTTPClient(url string) {{ .GoName }}HTTPClientInterface {
	return &{{ .GoName }}HTTPClient{url: url}
}

{{ $serviceName := .GoName }}
{{ $protoPackageName := .Desc.ParentFile.Package }}

{{ range .Methods }}
func (c *{{ $serviceName }}HTTPClient) {{ .GoName }}(ctx context.Context, in *{{ .Input.GoIdent.GoName }}) (*{{ .Output.GoIdent.GoName }}, error) {
	inJSON, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}
	
	req, err := http.NewRequestWithContext(ctx, "POST", c.url+"/{{ $protoPackageName }}.{{ $serviceName }}/{{ .GoName }}", bytes.NewBuffer(inJSON))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	out := new({{ .Output.GoIdent.GoName }})
	err = json.Unmarshal(body, out)
	if err != nil {
		return nil, err
	}

	return out, nil
}
{{ end }}
{{ end }}
`

func generateClient(gen *protogen.Plugin, file *protogen.File) {
	if len(file.Services) == 0 {
		return
	}
	filename := file.GeneratedFilenamePrefix + ".client.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	
	tmpl, err := template.New("client").Parse(clientTpl)
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
