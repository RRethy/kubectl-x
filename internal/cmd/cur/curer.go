package cur

import (
	"context"
	"fmt"

	"github.com/RRethy/kubectl-x/internal/kubeconfig"
	"k8s.io/cli-runtime/pkg/genericiooptions"
)

type Curer struct {
	kubeConfig kubeconfig.KubeConfig
	ioStreams  genericiooptions.IOStreams
}

func (c Curer) Cur(ctx context.Context) error {
	currentContext, err := c.kubeConfig.CurrentContext()
	if err != nil {
		return fmt.Errorf("getting current context: %w", err)
	}

	currentNamespace, err := c.kubeConfig.CurrentNamespace()
	if err != nil {
		return fmt.Errorf("getting current namespace: %w", err)
	}

	fmt.Fprintf(c.ioStreams.Out, "Current context: \"%s\"\n", currentContext)
	fmt.Fprintf(c.ioStreams.Out, "Current namespace: \"%s\"\n", currentNamespace)

	return nil
}
