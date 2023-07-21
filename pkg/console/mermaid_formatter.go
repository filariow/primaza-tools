package console

import (
	"github.com/primaza/primaza-tools/pkg/mermaid"
)

type MermaidFormatter struct{}

func (f *MermaidFormatter) Format(v any) ([]byte, error) {
	gh, err := mermaid.NewGraph(v)
	if err != nil {
		return nil, err
	}

	return []byte(gh.String()), nil
}
