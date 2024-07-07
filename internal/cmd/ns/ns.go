package ns

import (
	"context"
	"os"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"

	"github.com/RRethy/kubectl-x/internal/kubeconfig"
	"github.com/RRethy/kubectl-x/internal/kubernetes"
)

func Ns(ctx context.Context, configFlags *genericclioptions.ConfigFlags, namespace string) error {
	kubeConfig, err := kubeconfig.NewKubeConfig()
	if err != nil {
		return err
	}
	ioStreams := genericiooptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}
	k8sClient := kubernetes.NewClient(configFlags)
	return Nser{
		kubeConfig,
		ioStreams,
		k8sClient,
	}.Ns(ctx, namespace)
}