package primaza

import (
	"context"

	primazaiov1alpha1 "github.com/primaza/primaza/api/v1alpha1"
	"github.com/primaza/primaza/pkg/primaza/clustercontext"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ApplicationNamespaceClient struct {
	cli    client.Client
	config *rest.Config
}

func NewApplicationClientForClusterEnvironment(
	ctx context.Context,
	cli client.Client,
	ce primazaiov1alpha1.ClusterEnvironment,
) (
	*ApplicationNamespaceClient,
	error,
) {
	cfg, err := clustercontext.GetClusterRESTConfig(ctx, cli, ce.Namespace, ce.Spec.ClusterContextSecret)
	if err != nil {
		return nil, err
	}

	// TODO: add Primaza's CRDs here
	oc := client.Options{
		// Scheme: scheme,
		// Mapper: mapper,
	}
	c, err := client.New(cfg, oc)
	if err != nil {
		return nil, err
	}

	return &ApplicationNamespaceClient{cli: c, config: cfg}, nil
}

func (c *ApplicationNamespaceClient) GetServiceBindings(
	ctx context.Context,
	namespace string,
) (
	[]primazaiov1alpha1.ServiceBinding,
	error,
) {
	sbb := primazaiov1alpha1.ServiceBindingList{}
	o := client.ListOptions{Namespace: namespace}
	if err := c.cli.List(ctx, &sbb, &o); err != nil {
		return nil, err
	}

	return sbb.Items, nil
}
