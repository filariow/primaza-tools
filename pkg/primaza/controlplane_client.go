package primaza

import (
	"context"

	primazaiov1alpha1 "github.com/primaza/primaza/api/v1alpha1"
	"github.com/primaza/primaza/pkg/primaza/clustercontext"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ControlPlaneClient struct {
	cli client.Client
}

func NewClient(cli client.Client) *ControlPlaneClient {
	return &ControlPlaneClient{cli: cli}
}

func (c *ControlPlaneClient) ListClusterEnvironments(
	ctx context.Context,
	namespace string,
) (
	[]primazaiov1alpha1.ClusterEnvironment,
	error,
) {
	cee := primazaiov1alpha1.ClusterEnvironmentList{}
	o := client.ListOptions{Namespace: namespace}
	if err := c.cli.List(ctx, &cee, &o); err != nil {
		return nil, err
	}

	return cee.Items, nil
}

func (c *ControlPlaneClient) GetClusterContextSecret(
	ctx context.Context,
	ce *primazaiov1alpha1.ClusterEnvironment,
) (
	*corev1.Secret,
	error,
) {
	return clustercontext.GetClusterContextSecret(ctx, c.cli, ce)
}

func (c *ControlPlaneClient) NewApplicationClientForClusterEnvironment(
	ctx context.Context,
	ce primazaiov1alpha1.ClusterEnvironment,
) (
	*ApplicationNamespaceClient,
	error,
) {
	return NewApplicationClientForClusterEnvironment(ctx, c.cli, ce)
}
