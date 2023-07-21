package dependencies

import (
	"context"
	"errors"
	"fmt"

	"github.com/primaza/primaza-tools/pkg/mermaid"
	"github.com/primaza/primaza-tools/pkg/primaza"
	primazaiov1alpha1 "github.com/primaza/primaza/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
				Start: sb.Name,
				End:   c.Name,
			}
			g.Adjacencies = append(g.Adjacencies, a)
		}
	}

	return g, nil
}

type ServiceDependenciesCrawler struct {
	cli client.Client
}

func NewServiceDependeciesCrawler(cli client.Client) *ServiceDependenciesCrawler {
	return &ServiceDependenciesCrawler{cli: cli}
}

func (c *ServiceDependenciesCrawler) CrawlServiceDependencies(ctx context.Context, tenant string) ([]ServiceDependencies, error) {
	pcli := primaza.NewControlPlaneClient(c.cli)

	cee, err := pcli.ListClusterEnvironments(ctx, tenant)
	if err != nil {
		return nil, err
	}

	sdd := []ServiceDependencies{}
	errs := []error{}

	for _, ce := range cee {
		acli, err := pcli.NewApplicationClientForClusterEnvironment(ctx, ce)
		if err != nil {
			werr := fmt.Errorf("error building client for Cluster Environment '%s': %w", ce.Name, err)
			errs = append(errs, werr)
			continue
		}
		sd := ServiceDependencies{
			ClusterEnvironment: ce,
			ServiceBindings:    []primazaiov1alpha1.ServiceBinding{},
		}

		for _, ns := range ce.Spec.ApplicationNamespaces {
			sbb, err := acli.GetServiceBindings(ctx, ns)
			if err != nil {
				werr := fmt.Errorf(
					"error retrieving service bindings from namespace '%s' of Cluster Environment '%s'",
					ns, ce.Name)
				errs = append(errs, werr)
				continue
			}

			sd.ServiceBindings = append(sd.ServiceBindings, sbb...)
		}
		sdd = append(sdd, sd)
	}

	return sdd, errors.Join(errs...)
}
