package dependencies

import (
	"github.com/primaza/primaza-tools/pkg/mermaid"
	primazaiov1alpha1 "github.com/primaza/primaza/api/v1alpha1"
	"github.com/primaza/primaza/pkg/primaza/constants"
)

type Dependencies []ServiceDependencies

type ServiceDependencies struct {
	ClusterEnvironment primazaiov1alpha1.ClusterEnvironment
	ServiceBindings    []primazaiov1alpha1.ServiceBinding
}

func (d *ServiceDependencies) ToGraph() (mermaid.Graph, error) {
	g := mermaid.Graph{Name: d.ClusterEnvironment.Name, Adjacencies: []mermaid.Adjancency{}}

	for _, sb := range d.ServiceBindings {
		for _, c := range sb.Status.Connections {
			a := mermaid.Adjancency{
				Start: c.Name,
				End:   sb.Spec.ServiceEndpointDefinitionSecret,
				Text:  sb.Annotations[constants.BoundRegisteredServiceNameAnnotation],
			}
			g.Adjacencies = append(g.Adjacencies, a)
		}
	}

	return g, nil
}
