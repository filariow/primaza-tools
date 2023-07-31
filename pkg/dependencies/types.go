package dependencies

import (
	"encoding/json"
	"strings"

	"github.com/primaza/primaza-tools/pkg/mermaid"
	primazaiov1alpha1 "github.com/primaza/primaza/api/v1alpha1"
	"github.com/primaza/primaza/pkg/primaza/constants"
)

type Dependencies []ServiceDependencies

type ServiceDependencies struct {
	ClusterEnvironment primazaiov1alpha1.ClusterEnvironment
	ServiceBindings    []primazaiov1alpha1.ServiceBinding
	RegisteredServices []primazaiov1alpha1.RegisteredService
}

func (d *ServiceDependencies) ToGraph() (mermaid.Graph, error) {
	g := mermaid.Graph{Name: d.ClusterEnvironment.Name, Adjacencies: []mermaid.Adjancency{}, Nodes: []mermaid.Node{}}

	for _, sb := range d.ServiceBindings {
		for _, c := range sb.Status.Connections {
			a := mermaid.Adjancency{
				Start: c.Name,
				End:   sb.Annotations[constants.BoundRegisteredServiceNameAnnotation],
				Text:  sb.Name,
			}
			g.Adjacencies = append(g.Adjacencies, a)

			sbj, err := json.MarshalIndent(&sb, "", "  ")
			if err == nil {
				n := mermaid.Node{Name: a.Text, Description: strings.ReplaceAll(string(sbj), "\"", "'")}
				g.Nodes = append(g.Nodes, n)
			}
		}
	}

	for _, rs := range d.RegisteredServices {
		rsj, err := json.MarshalIndent(&rs, "", "  ")
		if err == nil {
			n := mermaid.Node{Name: rs.Name, Description: strings.ReplaceAll(string(rsj), "\"", "'")}
			g.Nodes = append(g.Nodes, n)
		}
	}

	return g, nil
}
