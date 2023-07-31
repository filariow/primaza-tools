package console

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"strings"

	"html/template"

	"github.com/primaza/primaza-tools/pkg/mermaid"
)

var (
	//go:embed resources/template.html
	HTMLTemplate string
)

type HTMLFormatter struct{}

type HTMLData struct {
	Graph     string
	NodesJSON string
}

func (f *HTMLFormatter) Format(v any) ([]byte, error) {
	gh, err := mermaid.NewGraph(v)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New("graph").Parse(HTMLTemplate)
	if err != nil {
		return nil, err
	}

	nj, err := json.Marshal(gh.Nodes)
	if err != nil {
		return nil, err
	}

	var buf []byte
	w := bytes.NewBuffer(buf)
	if err := tmpl.Execute(w, HTMLData{
		Graph: gh.StringFormat(func(s string) string {
			return strings.ReplaceAll(s, "\"", "#quot;")
		}),
		NodesJSON: string(nj),
	}); err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}
