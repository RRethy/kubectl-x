package ctx

import (
	"context"
	"fmt"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"

	"github.com/RRethy/kubectl-x/internal/cmd/ns"
	"github.com/RRethy/kubectl-x/internal/fzf"
	"github.com/RRethy/kubectl-x/internal/kubeconfig"
)

type Ctxer struct {
	kubeConfig           kubeconfig.KubeConfig
	ioStreams            genericiooptions.IOStreams
	configFlags          *genericclioptions.ConfigFlags
	resourceBuilderFlags *genericclioptions.ResourceBuilderFlags
}

func (c Ctxer) Ctx(ctx context.Context, contextSubstring, namespaceSubstring string, exactMatch bool) error {
	selectedContext, err := fzf.NewFzf(fzf.WithIOStreams(c.ioStreams), fzf.WithExactMatch(exactMatch)).Run(contextSubstring, c.kubeConfig.Contexts())
	if err != nil {
		return fmt.Errorf("selecting context: %s", err)
	}

	err = c.kubeConfig.SetContext(selectedContext)
	if err != nil {
		return fmt.Errorf("setting context: %w", err)
	}

	err = c.kubeConfig.Write()
	if err != nil {
		return fmt.Errorf("writing kubeconfig: %w", err)
	}

	fmt.Fprintf(c.ioStreams.Out, "Switched to context \"%s\".\n", selectedContext)

	return ns.Ns(ctx, c.configFlags, c.resourceBuilderFlags, namespaceSubstring, exactMatch)
}
