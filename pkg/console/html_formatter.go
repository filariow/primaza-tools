package console

import (
	"bytes"
	_ "embed"

	"html/template"

	"github.com/primaza/primaza-tools/pkg/mermaid"
)

var (
	//go:embed resources/template.html
	HTMLTemplate string
)

type HTMLFormatter struct{}

type HTMLData struct {
	Graph string
}

func (f *HTMLFormatter) Format(v any) ([]byte, error) {
	tmpl, err := template.New("graph").Parse(HTMLTemplate)
	if err != nil {
		return nil, err
	}

	gh, err := mermaid.NewGraph(v)
	if err != nil {
		return nil, err
	}

	var buf []byte
	w := bytes.NewBuffer(buf)
	if err := tmpl.Execute(w, HTMLData{Graph: gh.String()}); err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}
