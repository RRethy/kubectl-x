package shell

import (
	"context"
	"os"

	"github.com/RRethy/kubectl-x/internal/fzf"
	"github.com/RRethy/kubectl-x/internal/history"
	"github.com/RRethy/kubectl-x/internal/kubeconfig"
	"github.com/RRethy/kubectl-x/internal/kubernetes"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
)

func Ctx(ctx context.Context, configFlags *genericclioptions.ConfigFlags, resourceBuilderFlags *genericclioptions.ResourceBuilderFlags, contextSubstring, namespaceSubstring string, exactMatch bool) error {
	kubeConfig, err := kubeconfig.NewKubeConfig()
	if err != nil {
		return err
	}
	ioStreams := genericiooptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}
	k8sClient := kubernetes.NewClient(configFlags, resourceBuilderFlags)
	fzf := fzf.NewFzf(fzf.WithIOStreams(ioStreams), fzf.WithExactMatch(exactMatch))
	history, err := history.NewHistory(history.NewConfig())
	if err != nil {
		return err
	}
	return Ctxer{kubeConfig, ioStreams, k8sClient, fzf, history}.Ctx(ctx, contextSubstring, namespaceSubstring)
}
