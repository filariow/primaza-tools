package dependencies

import (
	"context"
	"errors"
	"fmt"

	"github.com/primaza/primaza-tools/pkg/primaza"
	primazaiov1alpha1 "github.com/primaza/primaza/api/v1alpha1"
	"github.com/primaza/primaza/pkg/primaza/constants"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ServiceDependenciesCrawler struct {
	cli client.Client
}

func NewServiceDependeciesCrawler(cli client.Client) *ServiceDependenciesCrawler {
	return &ServiceDependenciesCrawler{cli: cli}
}

func (c *ServiceDependenciesCrawler) fetchRegisteredService(
	ctx context.Context,
	pcli *primaza.ControlPlaneClient,
	tenant string,
	sb primazaiov1alpha1.ServiceBinding,
	fetched map[string]*primazaiov1alpha1.RegisteredService,
) (*primazaiov1alpha1.RegisteredService, error) {

	rsn, ok := sb.Annotations[constants.BoundRegisteredServiceNameAnnotation]
	if !ok {
		return nil, fmt.Errorf("missing RegisteredService annotation in ServiceBinding %s", sb.Name)
	}

	if f, ok := fetched[rsn]; ok {
		return f, nil
	}

	return pcli.GetRegisteredService(ctx, tenant, rsn)
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
			Workloads:          []unstructured.Unstructured{},
		}

		for _, ns := range ce.Spec.ApplicationNamespaces {
			ww := map[primaza.WorkloadKey]*unstructured.Unstructured{}

			// fetch service bindings
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
				// fetch registered services
				rs, err := c.fetchRegisteredService(ctx, pcli, tenant, sb, rss)
				if err != nil {
					errs = append(errs, err)
				} else {
					rss[rs.Name] = rs
				}

				// fetch workload
				a := sb.Spec.Application
				for _, c := range sb.Status.Connections {
					k := primaza.WorkloadKey{
						ApiVersion: a.APIVersion,
						Kind:       a.Kind,
						Name:       c.Name,
					}
					if _, ok := ww[k]; ok {
						continue
					}

					w, err := pcli.GetWorkload(ctx, sb.Namespace, k)
					if err != nil {
						errs = append(errs, err)
					} else {
						ww[k] = w
					}
				}
			}

			// store registered services
			for _, rs := range rss {
				sd.RegisteredServices = append(sd.RegisteredServices, *rs)
			}

			// store workloads
			for _, w := range ww {
				sd.Workloads = append(sd.Workloads, *w)
			}
		}
		sdd = append(sdd, sd)
	}

	return sdd, errors.Join(errs...)
}
