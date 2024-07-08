package ctx

import (
	"context"
	"fmt"

	"k8s.io/cli-runtime/pkg/genericiooptions"

	"github.com/RRethy/kubectl-x/internal/cmd/ns"
	"github.com/RRethy/kubectl-x/internal/fzf"
	"github.com/RRethy/kubectl-x/internal/history"
	"github.com/RRethy/kubectl-x/internal/kubeconfig"
	"github.com/RRethy/kubectl-x/internal/kubernetes"
)

type Ctxer struct {
	KubeConfig kubeconfig.Interface
	IoStreams  genericiooptions.IOStreams
	K8sClient  kubernetes.Interface
	Fzf        fzf.Interface
	History    history.Interface
}

func (c Ctxer) Ctx(ctx context.Context, contextSubstring, namespaceSubstring string) error {
	var selectedContext string
	var selectedNamespace string
	var err error
	if contextSubstring == "-" {
		selectedContext, err = c.History.Get("context", 0)
		if err != nil {
			return fmt.Errorf("getting context from history: %s", err)
		}

		selectedNamespace, err = c.KubeConfig.GetNamespaceForContext(selectedContext)
		if err != nil {
			return fmt.Errorf("getting namespace for context: %s", err)
		}
	} else {
		selectedContext, err = c.Fzf.Run(contextSubstring, c.KubeConfig.Contexts())
		if err != nil {
			return fmt.Errorf("selecting context: %s", err)
		}
	}

	c.History.Add("context", selectedContext)

	err = c.KubeConfig.SetContext(selectedContext)
	if err != nil {
		return fmt.Errorf("setting context: %w", err)
	}

	err = c.History.Write()
	if err != nil {
		fmt.Fprintf(c.IoStreams.ErrOut, "writing history: %s\n", err)
	}

	err = c.KubeConfig.Write()
	if err != nil {
		return fmt.Errorf("writing kubeconfig: %w", err)
	}

	fmt.Fprintf(c.IoStreams.Out, "Switched to context \"%s\".\n", selectedContext)

	if selectedNamespace == "" {
		return ns.Nser{KubeConfig: c.KubeConfig, IoStreams: c.IoStreams, K8sClient: c.K8sClient, Fzf: c.Fzf, History: c.History}.Ns(ctx, namespaceSubstring)
	}

	err = c.KubeConfig.SetNamespace(selectedNamespace)
	if err != nil {
		return fmt.Errorf("setting namespace: %w", err)
	}

	err = c.KubeConfig.Write()
	if err != nil {
		return fmt.Errorf("writing kubeconfig: %w", err)
	}

	fmt.Fprintf(c.IoStreams.Out, "Switched to namespace \"%s\".\n", selectedNamespace)
	return nil
}
