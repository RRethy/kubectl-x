package ctx

import (
	"context"
	"fmt"

	"k8s.io/cli-runtime/pkg/genericiooptions"

	"github.com/RRethy/kubectl-x/pkg/cli/ns"
	"github.com/RRethy/kubectl-x/pkg/fzf"
	"github.com/RRethy/kubectl-x/pkg/history"
	"github.com/RRethy/kubectl-x/pkg/kubeconfig"
	"github.com/RRethy/kubectl-x/pkg/kubernetes"
)

type Ctxer struct {
	KubeConfig kubeconfig.Interface
	IoStreams  genericiooptions.IOStreams
	K8sClient  kubernetes.Interface
	Fzf        fzf.Interface
	History    history.Interface
}

func NewCtxer(kubeConfig kubeconfig.Interface, ioStreams genericiooptions.IOStreams, k8sClient kubernetes.Interface, fzf fzf.Interface, history history.Interface) Ctxer {
	return Ctxer{
		KubeConfig: kubeConfig,
		IoStreams:  ioStreams,
		K8sClient:  k8sClient,
		Fzf:        fzf,
		History:    history,
	}
}

func (c Ctxer) Ctx(ctx context.Context, contextSubstring, namespaceSubstring string) error {
	var selectedContext string
	var selectedNamespace string
	var err error
	if contextSubstring == "-" {
		selectedContext, err = c.History.Get("context", 1)
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
		nser := ns.NewNser(c.KubeConfig, c.IoStreams, c.K8sClient, c.Fzf, c.History)
		return nser.Ns(ctx, namespaceSubstring)
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
