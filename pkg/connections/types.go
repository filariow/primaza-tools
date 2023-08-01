package connections

import (
	"encoding/json"

	"github.com/primaza/primaza-tools/pkg/mermaid"
	primazaiov1alpha1 "github.com/primaza/primaza/api/v1alpha1"
	"github.com/primaza/primaza/pkg/primaza/constants"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Connections []ServiceConnections

type ServiceConnections struct {
	ClusterEnvironment primazaiov1alpha1.ClusterEnvironment
	ServiceBindings    []primazaiov1alpha1.ServiceBinding
	RegisteredServices []primazaiov1alpha1.RegisteredService
	Workloads          []unstructured.Unstructured
}

func (d *ServiceConnections) ToGraph() (mermaid.Graph, error) {
	g := mermaid.Graph{Name: d.ClusterEnvironment.Name, Adjacencies: []mermaid.Adjancency{}, Nodes: []mermaid.Node{}}

	for _, sb := range d.ServiceBindings {
		for _, c := range sb.Status.Connections {
			sb := sb

			a := mermaid.Adjancency{
				Start: c.Name,
				End:   sb.Annotations[constants.BoundRegisteredServiceNameAnnotation],
				Text:  sb.Name,
			}
			g.Adjacencies = append(g.Adjacencies, a)

			sbj, err := json.MarshalIndent(&sb, "", "  ")
			if err == nil {
				n := mermaid.Node{Name: a.Text, Description: string(sbj)}
				g.Nodes = append(g.Nodes, n)
			}
		}
	}

	for _, rs := range d.RegisteredServices {
		rs := rs

		rsj, err := json.MarshalIndent(&rs, "", "  ")
		if err == nil {
			n := mermaid.Node{Name: rs.Name, Description: string(rsj)}
			g.Nodes = append(g.Nodes, n)
		}
	}

	for _, w := range d.Workloads {
		w := w
		wj, err := json.MarshalIndent(&w, "", "  ")
		if err == nil {
			n := mermaid.Node{Name: w.GetName(), Description: string(wj)}
			g.Nodes = append(g.Nodes, n)
		}
	}
	return g, nil
}
