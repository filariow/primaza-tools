package connections

import (
	"encoding/json"

	"github.com/primaza/tools/pkg/mermaid"
	primazaiov1alpha1 "github.com/primaza/primaza/api/v1alpha1"
	"github.com/primaza/primaza/pkg/primaza/constants"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Connections []ServiceConnections

type ServiceConnections struct {
	ClusterEnvironment primazaiov1alpha1.ClusterEnvironment
	ServiceBindings    []primazaiov1alpha1.ServiceBinding
	RegisteredServices map[string]primazaiov1alpha1.RegisteredService
	Workloads          []unstructured.Unstructured
	Services           []unstructured.Unstructured
}

func (d *ServiceConnections) ToGraph() (mermaid.Graph, error) {
	g := mermaid.Graph{Name: d.ClusterEnvironment.Name, Links: []mermaid.Link{}, Nodes: []mermaid.Node{}}

	for _, sb := range d.ServiceBindings {
		for _, c := range sb.Status.Connections {
			sb := sb

			rsn := sb.Annotations[constants.BoundRegisteredServiceNameAnnotation]
			l := mermaid.Link{
				c.Name,
				sb.Name,
				rsn,
			}
			if s, ok := d.RegisteredServices[rsn]; ok {
				if n := s.ObjectMeta.Annotations[constants.ServiceNameAnnotation]; n != "" {
					l = append(l, n)
				}
			}
			g.Links = append(g.Links, l)

			sbj, err := json.MarshalIndent(&sb, "", "  ")
			if err == nil {
				n := mermaid.Node{Name: sb.GetName(), Description: string(sbj)}
				g.Nodes = append(g.Nodes, n)
			}
		}
	}

	for _, rs := range d.RegisteredServices {
		rs := rs

		rsj, err := json.MarshalIndent(&rs, "", "  ")
		if err == nil {
			n := mermaid.Node{Name: rs.GetName(), Description: string(rsj)}
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

	for _, s := range d.Services {
		s := s

		sj, err := json.MarshalIndent(&s, "", "  ")
		if err == nil {
			n := mermaid.Node{Name: s.GetName(), Description: string(sj)}
			g.Nodes = append(g.Nodes, n)
		}
	}
	return g, nil
}
