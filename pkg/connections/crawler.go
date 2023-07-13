package connections

import (
	"context"
	"fmt"

	primazaiov1alpha1 "github.com/primaza/primaza/api/v1alpha1"
	"github.com/primaza/primaza/pkg/primaza/clustercontext"
	"github.com/primaza/primaza/pkg/primaza/constants"
	"github.com/primaza/tools/pkg/primaza"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	ErrMissingAnnotations = fmt.Errorf("missing required annotations")
)

type ServiceConnectionsCrawler struct {
	cli client.Client
}

func NewServiceDependeciesCrawler(cli client.Client) *ServiceConnectionsCrawler {
	return &ServiceConnectionsCrawler{cli: cli}
}

func (c *ServiceConnectionsCrawler) fetchRegisteredService(
	ctx context.Context,
	pcli *primaza.ControlPlaneClient,
	tenant string,
	sb primazaiov1alpha1.ServiceBinding,
	cache map[string]primazaiov1alpha1.RegisteredService,
) (*primazaiov1alpha1.RegisteredService, error) {
	rsn, ok := sb.Annotations[constants.BoundRegisteredServiceNameAnnotation]
	if !ok {
		return nil, fmt.Errorf("%v: missing RegisteredService annotation in ServiceBinding %s", ErrMissingAnnotations, sb.Name)
	}

	if f, ok := cache[rsn]; ok {
		return &f, nil
	}

	return pcli.GetRegisteredService(ctx, tenant, rsn)
}

func (c *ServiceConnectionsCrawler) fetchService(
	ctx context.Context,
	cee []primazaiov1alpha1.ClusterEnvironment,
	tenant string,
	rs primazaiov1alpha1.RegisteredService,
	cache map[string]*unstructured.Unstructured,
) (*unstructured.Unstructured, error) {

	if !c.checkServiceAnnotations(&rs) {
		return nil, fmt.Errorf("%w in RegisteredService %s", ErrMissingAnnotations, rs.GetName())
	}

	s := rs.Annotations[constants.ServiceUIDAnnotation]

	if c, ok := cache[s]; ok {
		return c, nil
	}

	cen := rs.Annotations[constants.ClusterEnvironmentAnnotation]
	ce := func() *primazaiov1alpha1.ClusterEnvironment {
		for _, c := range cee {
			if c.GetName() == cen {
				return &c
			}
		}
		return nil
	}()

	cli, err := clustercontext.CreateClient(ctx, c.cli, *ce, c.cli.Scheme(), c.cli.RESTMapper())
	if err != nil {
		return nil, err
	}

	n := types.NamespacedName{
		Namespace: rs.Annotations[constants.ServiceNamespaceAnnotation],
		Name:      rs.Annotations[constants.ServiceNameAnnotation],
	}
	a := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": rs.Annotations[constants.ServiceAPIVersionAnnotation],
			"kind":       rs.Annotations[constants.ServiceKindAnnotation],
		},
	}

	if err := cli.Get(ctx, n, &a); err != nil {
		return nil, err
	}
	return &a, nil
}

func (c *ServiceConnectionsCrawler) checkServiceAnnotations(rs *primazaiov1alpha1.RegisteredService) bool {
	aa := []string{
		constants.ServiceAPIVersionAnnotation,
		constants.ServiceKindAnnotation,
		constants.ServiceNameAnnotation,
		constants.ServiceNamespaceAnnotation,
		constants.ServiceUIDAnnotation,
		constants.ClusterEnvironmentAnnotation,
	}

	for _, a := range aa {
		if _, ok := rs.ObjectMeta.Annotations[a]; !ok {
			return false
		}
	}
	return true
}

func (c *ServiceConnectionsCrawler) CrawlServiceConnections(ctx context.Context, tenant string) ([]ServiceConnections, []error, error) {
	pcli := primaza.NewControlPlaneClient(c.cli)

	cee, err := pcli.ListClusterEnvironments(ctx, tenant)
	if err != nil {
		return nil, nil, err
	}

	sdd := []ServiceConnections{}
	errs := []error{}
	rss := map[string]primazaiov1alpha1.RegisteredService{}
	ss := map[string]*unstructured.Unstructured{}

	for _, ce := range cee {
		acli, err := pcli.NewApplicationClientForClusterEnvironment(ctx, ce)
		if err != nil {
			werr := fmt.Errorf("error building client for Cluster Environment '%s': %w", ce.Name, err)
			errs = append(errs, werr)
			continue
		}
		sd := ServiceConnections{
			ClusterEnvironment: ce,
			Workloads:          []unstructured.Unstructured{},
			ServiceBindings:    []primazaiov1alpha1.ServiceBinding{},
			Services:           []unstructured.Unstructured{},
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
					rss[rs.Name] = *rs

					s, err := c.fetchService(ctx, cee, tenant, *rs, ss)
					if err != nil {
						errs = append(errs, err)
					} else {
						ss[string(s.GetUID())] = s
					}
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
			sd.RegisteredServices = rss

			// store services
			for _, s := range ss {
				sd.Services = append(sd.Services, *s)
			}

			// store workloads
			for _, w := range ww {
				sd.Workloads = append(sd.Workloads, *w)
			}
		}
		sdd = append(sdd, sd)
	}

	return sdd, errs, nil
}
