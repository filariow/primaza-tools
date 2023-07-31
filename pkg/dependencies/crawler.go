package dependencies

import (
	"context"
	"errors"
	"fmt"

	"github.com/primaza/primaza-tools/pkg/primaza"
	primazaiov1alpha1 "github.com/primaza/primaza/api/v1alpha1"
	"github.com/primaza/primaza/pkg/primaza/constants"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

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
	rss := map[string]*primazaiov1alpha1.RegisteredService{}

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
			RegisteredServices: []primazaiov1alpha1.RegisteredService{},
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

			for _, sb := range sbb {
				rsn, ok := sb.Annotations[constants.BoundRegisteredServiceNameAnnotation]
				if !ok {
					continue
				}

				if _, ok := rss[rsn]; ok {
					continue
				}

				rs, err := pcli.GetRegisteredService(ctx, tenant, rsn)
				if err != nil {
					continue
				}
				rss[rsn] = rs
			}
			for _, rs := range rss {
				sd.RegisteredServices = append(sd.RegisteredServices, *rs)
			}
		}
		sdd = append(sdd, sd)

	}

	return sdd, errors.Join(errs...)
}
