package mermaid

import (
	"fmt"
	"strings"
)

var (
	InvalidTypeError error = fmt.Errorf("given object can not be serialized as mermaid")
)

type Graph struct {
	Name        string
	Adjacencies []Adjancency
}

type Adjancency struct {
	Start string
	End   string
	Text  string
}

func (m Graph) String() string {
	b := strings.Builder{}

	b.WriteString(fmt.Sprintf("graph TD;\n"))
	b.WriteString(fmt.Sprintf("\taccTitle: %s;\n", m.Name))

	for i, a := range m.Adjacencies {
		b.WriteString(fmt.Sprintf("\t%s --%s--> %s;", a.Start, a.Text, a.End))
		if i < len(m.Adjacencies) {
			b.WriteRune('\n')
		}
	}

	return b.String()
}

type Graphable interface {
	ToGraph() (Graph, error)
}

func NewGraph(v any) (*Graph, error) {
	g, ok := v.(Graphable)
	if !ok {
		return nil, InvalidTypeError
	}

	gh, err := g.ToGraph()
	if err != nil {
		return nil, err
	}
	return &gh, nil
}
