package console

import (
	"fmt"
)

type Formatter interface {
	Format(any) ([]byte, error)
}

func NewFormatter(format Format) (Formatter, error) {
	switch format {
	case FormatHTML:
		return &HTMLFormatter{}, nil
	case FormatJson:
		return &JsonFormatter{}, nil
	case FormatMermaid:
		return &MermaidFormatter{}, nil
	}

	return nil, fmt.Errorf("Invalid format: %s", format)
}
