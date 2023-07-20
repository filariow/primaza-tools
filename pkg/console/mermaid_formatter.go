package console

import (
	"log"

	"github.com/primaza/primaza-tools/pkg/mermaid"
)

type MermaidFormatter struct{}

func (f *MermaidFormatter) Format(v any) ([]byte, error) {
	log.Println("mermaid formatter")
	gh, err := mermaid.NewGraph(v)
	if err != nil {
		return nil, err
	}

	return []byte(gh.String()), nil
}
