package primaza

import (
	"context"

	primazaiov1alpha1 "github.com/primaza/primaza/api/v1alpha1"
	"github.com/primaza/primaza/pkg/primaza/clustercontext"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ControlPlaneClient struct {
	cli client.Client
}

func NewControlPlaneClient(cli client.Client) *ControlPlaneClient {
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

func (c *ControlPlaneClient) GetRegisteredService(ctx context.Context, namespace, name string) (*primazaiov1alpha1.RegisteredService, error) {
	var rs primazaiov1alpha1.RegisteredService
	o := types.NamespacedName{Name: name, Namespace: namespace}
	if err := c.cli.Get(ctx, o, &rs); err != nil {
		return nil, err
	}
	return &rs, nil
}

type WorkloadKey struct {
	ApiVersion string
	Kind       string
	Name       string
}

func (c *ControlPlaneClient) GetWorkload(ctx context.Context, namespace string, k WorkloadKey) (*unstructured.Unstructured, error) {
	n := types.NamespacedName{Namespace: namespace, Name: k.Name}
	a := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": k.ApiVersion,
			"kind":       k.Kind,
			"metadata": map[string]interface{}{
				"name":      k.Name,
				"namespace": namespace,
			},
		},
	}

	if err := c.cli.Get(ctx, n, &a); err != nil {
		return nil, err
	}
	return &a, nil
}
