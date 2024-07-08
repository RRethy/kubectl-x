package ns

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericiooptions"

	"github.com/RRethy/kubectl-x/internal/fzf"
	"github.com/RRethy/kubectl-x/internal/history"
	"github.com/RRethy/kubectl-x/internal/kubeconfig"
	"github.com/RRethy/kubectl-x/internal/kubernetes"
)

type Nser struct {
	KubeConfig kubeconfig.Interface
	IoStreams  genericiooptions.IOStreams
	K8sClient  kubernetes.Interface
	Fzf        fzf.Interface
	History    history.Interface
}

func (n Nser) Ns(ctx context.Context, namespace string) error {
	var selectedNamespace string
	if strings.HasPrefix(namespace, "-") {
		if namespace == "-" {
			namespace = "-1"
		}
		num, err := strconv.ParseInt(strings.TrimPrefix(namespace, "-"), 10, 8)
		if err != nil {
			return fmt.Errorf("parsing namespace argument: %s", err)
		}

		selectedNamespace, err = n.History.Get("namespace", int(num))
		if err != nil {
			return fmt.Errorf("getting namespace from history: %s", err)
		}
	} else {
		namespaces, err := kubernetes.List[*corev1.Namespace](ctx, n.K8sClient)
		if err != nil {
			return fmt.Errorf("listing namespaces: %s", err)
		}

		namespaceNames := make([]string, len(namespaces))
		for i, ns := range namespaces {
			namespaceNames[i] = ns.Name
		}

		selectedNamespace, err = n.Fzf.Run(namespace, namespaceNames)
		if err != nil {
			return fmt.Errorf("selecting namespace: %s", err)
		}
	}

	err := n.KubeConfig.SetNamespace(selectedNamespace)
	if err != nil {
		return fmt.Errorf("setting namespace: %w", err)
	}

	n.History.Add("namespace", selectedNamespace)

	err = n.KubeConfig.Write()
	if err != nil {
		return fmt.Errorf("writing kubeconfig: %w", err)
	}

	err = n.History.Write()
	if err != nil {
		fmt.Fprintf(n.IoStreams.ErrOut, "writing history: %s\n", err)
	}

	fmt.Fprintf(n.IoStreams.Out, "Switched to namespace \"%s\".\n", selectedNamespace)

	return nil
}
