package mermaid

import (
	"fmt"
	"strings"
)

type Graph struct {
	Name        string
	Adjacencies []Adjancency
}

type Adjancency struct {
	Start string
	End   string
}

func (m Graph) String() string {
	b := strings.Builder{}

	b.WriteString(fmt.Sprintf("graph TD;\n"))
	b.WriteString(fmt.Sprintf("\taccTitle: %s;\n", m.Name))

	for i, a := range m.Adjacencies {
		b.WriteString(fmt.Sprintf("\t%s --> %s;", a.Start, a.End))
		if i < len(m.Adjacencies) {
			b.WriteRune('\n')
		}
	}

	return b.String()
}
