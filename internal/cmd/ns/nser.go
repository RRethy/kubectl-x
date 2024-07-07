package ns

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericiooptions"

	"github.com/RRethy/kubectl-x/internal/fzf"
	"github.com/RRethy/kubectl-x/internal/kubeconfig"
	"github.com/RRethy/kubectl-x/internal/kubernetes"
)

type Nser struct {
	KubeConfig kubeconfig.Interface
	IoStreams  genericiooptions.IOStreams
	K8sClient  kubernetes.Interface
	Fzf        fzf.Interface
}

func (n Nser) Ns(ctx context.Context, namespace string) error {
	namespaces, err := kubernetes.List[*corev1.Namespace](ctx, n.K8sClient)
	if err != nil {
		return fmt.Errorf("listing namespaces: %s", err)
	}

	namespaceNames := make([]string, len(namespaces))
	for i, ns := range namespaces {
		namespaceNames[i] = ns.Name
	}

	selectedNamespace, err := n.Fzf.Run(namespace, namespaceNames)
	if err != nil {
		return fmt.Errorf("selecting namespace: %s", err)
	}

	err = n.KubeConfig.SetNamespace(selectedNamespace)
	if err != nil {
		return fmt.Errorf("setting namespace: %w", err)
	}

	err = n.KubeConfig.Write()
	if err != nil {
		return fmt.Errorf("writing kubeconfig: %w", err)
	}

	fmt.Fprintf(n.IoStreams.Out, "Switched to namespace \"%s\".\n", selectedNamespace)

	return nil
}
