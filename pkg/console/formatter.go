package console

import (
	"fmt"
	"log"
)

type Formatter interface {
	Format(any) ([]byte, error)
}

func NewFormatter(format Format) (Formatter, error) {
	switch format {
	case FormatJson:
		return &JsonFormatter{}, nil
	case FormatMermaid:
		log.Println("using mermaid formatter")
		return &MermaidFormatter{}, nil
	}

	return nil, fmt.Errorf("Invalid format: %s", format)
}
