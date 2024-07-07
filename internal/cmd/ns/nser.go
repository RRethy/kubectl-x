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
	kubeConfig kubeconfig.KubeConfig
	ioStreams  genericiooptions.IOStreams
	k8sClient  kubernetes.Interface
}

func (n Nser) Ns(ctx context.Context, namespace string) error {
	namespaces, err := kubernetes.List[*corev1.Namespace](ctx, n.k8sClient)
	if err != nil {
		return fmt.Errorf("listing namespaces: %s", err)
	}

	namespaceNames := make([]string, len(namespaces))
	for i, ns := range namespaces {
		namespaceNames[i] = ns.Name
	}

	selectedNamespace, err := fzf.NewFzf(fzf.WithIOStreams(n.ioStreams)).Run(namespaceNames)
	if err != nil {
		return fmt.Errorf("selecting namespace: %s", err)
	}

	err = n.kubeConfig.SetNamespace(selectedNamespace)
	if err != nil {
		return fmt.Errorf("setting namespace: %w", err)
	}

	err = n.kubeConfig.Write()
	if err != nil {
		return fmt.Errorf("writing kubeconfig: %w", err)
	}

	fmt.Fprintf(n.ioStreams.Out, "Switched to namespace \"%s\".\n", selectedNamespace)

	return nil
}
