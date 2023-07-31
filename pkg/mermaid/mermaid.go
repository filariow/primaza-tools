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
	Nodes       []Node
}

type Node struct {
	Name        string
	Description string
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

	for _, a := range m.Adjacencies {
		b.WriteString(fmt.Sprintf("\t%s --> %s;\n", a.Start, a.Text))
		b.WriteString(fmt.Sprintf("\t%s --> %s;\n", a.Text, a.End))
	}

	if len(m.Nodes) != 0 {
		b.WriteString("\n")
	}

	for _, n := range m.Nodes {
		b.WriteString(fmt.Sprintf("\tclick %s callback \"%s\"\n", n.Name, n.Description))
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
