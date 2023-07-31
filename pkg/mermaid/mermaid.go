package mermaid

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
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
	f := func(s string) string {
		return html.EscapeString(s)
	}
	return m.StringFormat(f)
}

func (m Graph) StringUnescaped() string {
	return m.StringFormat(func(s string) string { return s })
}

func (m Graph) StringFormat(f func(string) string) string {
	b := strings.Builder{}

	b.WriteString(fmt.Sprintf("graph TD;\n"))
	b.WriteString(fmt.Sprintf("\taccTitle: %s;\n", f(m.Name)))

	for _, a := range m.Adjacencies {
		b.WriteString(fmt.Sprintf("\t%s --> %s;\n", f(a.Start), f(a.Text)))
		b.WriteString(fmt.Sprintf("\t%s --> %s;\n", f(a.Text), f(a.End)))
	}

	if len(m.Nodes) != 0 {
		b.WriteString("\n")
	}

	for _, n := range m.Nodes {
		b.WriteString(fmt.Sprintf("\tclick %s call callback()\n", f(n.Name)))
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
