package ctx

import (
	"context"
	"os"

	"github.com/RRethy/kubectl-x/internal/kubeconfig"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
)

func Ctx(ctx context.Context, configFlags *genericclioptions.ConfigFlags, contextSubstring, namespaceSubstring string, exactMatch bool) error {
	kubeConfig, err := kubeconfig.NewKubeConfig()
	if err != nil {
		return err
	}
	ioStreams := genericiooptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}
	return Ctxer{kubeConfig, ioStreams, configFlags}.Ctx(ctx, contextSubstring, namespaceSubstring, exactMatch)
}
