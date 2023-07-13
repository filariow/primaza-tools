package dependencies

import (
	"context"
	"errors"
	"fmt"

	"github.com/primaza/primaza-tools/pkg/primaza"
	primazaiov1alpha1 "github.com/primaza/primaza/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ServiceDependency struct {
	ClusterEnvironment   primazaiov1alpha1.ClusterEnvironment
	ApplicationNamespace string
	ServiceBinding       primazaiov1alpha1.ServiceBinding
}

type ServiceDependenciesCrawler struct {
	cli client.Client
}

func NewServiceDependeciesCrawler(cli client.Client) *ServiceDependenciesCrawler {
	return &ServiceDependenciesCrawler{cli: cli}
}

func (c *ServiceDependenciesCrawler) CrawlServiceDependencies(ctx context.Context, tenant string) ([]ServiceDependency, error) {
	pcli := primaza.NewClient(c.cli)

	cee, err := pcli.ListClusterEnvironments(ctx, tenant)
	if err != nil {
		return nil, err
	}

	sdd := []ServiceDependency{}
	errs := []error{}

	for _, ce := range cee {
		acli, err := pcli.NewApplicationClientForClusterEnvironment(ctx, ce)
		if err != nil {
			werr := fmt.Errorf("error building client for Cluster Environment '%s': %w", ce.Name, err)
			errs = append(errs, werr)
			continue
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

			for _, sb := range sbb {
				sd := ServiceDependency{
					ClusterEnvironment:   ce,
					ApplicationNamespace: ns,
					ServiceBinding:       sb,
				}
				sdd = append(sdd, sd)
			}
		}
	}

	return sdd, errors.Join(errs...)
}
