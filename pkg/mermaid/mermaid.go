package mermaid

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

var (
	ErrInvalidType error = fmt.Errorf("given object can not be serialized as mermaid")
)

type Graph struct {
	Name  string
	Links []Link
	Nodes []Node
}

type Node struct {
	Name        string
	Description string
}

type Link []string

func (a Link) String() string {
	nn := make([]string, len(a))

	for i, n := range a {
		sf, ef := a.nodeShapes(i)
		nn[i] = fmt.Sprintf("%s%s%s%s", n, sf, n, ef)
	}

	return strings.Join(nn, " --> ")
}

func (a Link) nodeShapes(depth int) (string, string) {
	switch depth {
	case 0:
		return "[/", "\\]"
	case 1:
		return "{{", "}}"
	case 2:
		return "[\\", "/]"
	default:
		return "[", "]"
	}
}

func (m Graph) String() string {
	return m.StringFormat(html.EscapeString)
}

func (m Graph) StringUnescaped() string {
	return m.StringFormat(func(s string) string { return s })
}

func (m Graph) StringFormat(f func(string) string) string {
	b := strings.Builder{}

	b.WriteString("graph TD;\n")
	b.WriteString(fmt.Sprintf("\taccTitle: %s;\n", f(m.Name)))

	for _, a := range m.Links {
		b.WriteString(fmt.Sprintf("\t%s;\n", a))
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
		return nil, ErrInvalidType
	}

	gh, err := g.ToGraph()
	if err != nil {
		return nil, err
	}
	return &gh, nil
}
