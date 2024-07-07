package ctx

import (
	"context"
	"fmt"

	"k8s.io/cli-runtime/pkg/genericiooptions"

	"github.com/RRethy/kubectl-x/internal/cmd/ns"
	"github.com/RRethy/kubectl-x/internal/fzf"
	"github.com/RRethy/kubectl-x/internal/kubeconfig"
	"github.com/RRethy/kubectl-x/internal/kubernetes"
)

type Ctxer struct {
	KubeConfig kubeconfig.Interface
	IoStreams  genericiooptions.IOStreams
	K8sClient  kubernetes.Interface
	Fzf        fzf.Interface
}

func (c Ctxer) Ctx(ctx context.Context, contextSubstring, namespaceSubstring string) error {
	selectedContext, err := c.Fzf.Run(contextSubstring, c.KubeConfig.Contexts())
	if err != nil {
		return fmt.Errorf("selecting context: %s", err)
	}

	err = c.KubeConfig.SetContext(selectedContext)
	if err != nil {
		return fmt.Errorf("setting context: %w", err)
	}

	err = c.KubeConfig.Write()
	if err != nil {
		return fmt.Errorf("writing kubeconfig: %w", err)
	}

	fmt.Fprintf(c.IoStreams.Out, "Switched to context \"%s\".\n", selectedContext)

	return ns.Nser{
		KubeConfig: c.KubeConfig,
		IoStreams:  c.IoStreams,
		K8sClient:  c.K8sClient,
		Fzf:        c.Fzf,
	}.Ns(ctx, namespaceSubstring)
}
