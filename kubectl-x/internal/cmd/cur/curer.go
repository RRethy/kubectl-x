package cur

import (
	"context"
	"fmt"

	"github.com/RRethy/kubectl-x/internal/kubeconfig"
	"k8s.io/cli-runtime/pkg/genericiooptions"
)

type Curer struct {
	KubeConfig kubeconfig.Interface
	IoStreams  genericiooptions.IOStreams
}

func (c Curer) Cur(ctx context.Context) error {
	currentContext, err := c.KubeConfig.GetCurrentContext()
	if err != nil {
		return fmt.Errorf("getting current context: %w", err)
	}

	currentNamespace, err := c.KubeConfig.GetCurrentNamespace()
	if err != nil {
		return fmt.Errorf("getting current namespace: %w", err)
	}
	if currentNamespace == "" {
		currentNamespace = "default"
	}

	fmt.Fprintf(c.IoStreams.Out, "--context %s --namespace %s\n", currentContext, currentNamespace)

	return nil
}
