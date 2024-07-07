package cur

import (
	"context"
	"os"

	"github.com/RRethy/kubectl-x/internal/kubeconfig"
	"k8s.io/cli-runtime/pkg/genericiooptions"
)

func Cur(ctx context.Context) error {
	kubeConfig, err := kubeconfig.NewKubeConfig()
	if err != nil {
		return err
	}
	ioStreams := genericiooptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}
	return Curer{kubeConfig, ioStreams}.Cur(ctx)
}
