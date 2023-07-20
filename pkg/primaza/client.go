package primaza

import (
	primazaiov1alpha1 "github.com/primaza/primaza/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewClient(cfg *rest.Config) (client.Client, error) {
	scheme := runtime.NewScheme()
	utilruntime.Must(primazaiov1alpha1.AddToScheme(scheme))
	utilruntime.Must(corev1.AddToScheme(scheme))

	oc := client.Options{
		Scheme: scheme,
		// Mapper: mapper,
	}
	return client.New(cfg, oc)
}
