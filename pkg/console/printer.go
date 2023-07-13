package console

import (
	"fmt"
)

type Printer interface {
	Println(any) error
}

func NewPrinter(format Format) (Printer, error) {
	f, err := NewFormatter(format)
	if err != nil {
		return nil, err
	}

	return &printer{Formatter: f}, nil
}

func NewPrinterOrDie(format Format) Printer {
	f, err := NewFormatter(format)
	if err != nil {
		panic(err)
	}

	return &printer{Formatter: f}
}

type printer struct {
	Formatter Formatter
}

func (p *printer) Println(v any) error {
	b, err := p.Formatter.Format(v)
	if err != nil {
		return err
	}

	fmt.Println(string(b))
	return nil
}
