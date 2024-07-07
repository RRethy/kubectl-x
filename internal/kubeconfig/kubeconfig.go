package kubeconfig

import (
	"errors"
	"fmt"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

type KubeConfig struct {
	configAccess clientcmd.ConfigAccess
	apiConfig    *api.Config
}

func NewKubeConfig() (KubeConfig, error) {
	configAccess := clientcmd.NewDefaultPathOptions()
	config, err := configAccess.GetStartingConfig()
	if err != nil {
		return KubeConfig{}, err
	}

	return KubeConfig{configAccess: configAccess, apiConfig: config}, nil
}

func (kubeConfig KubeConfig) Contexts() []string {
	contexts := make([]string, 0, len(kubeConfig.apiConfig.Contexts))
	for context := range kubeConfig.apiConfig.Contexts {
		contexts = append(contexts, context)
	}
	return contexts
}

func (kubeConfig KubeConfig) SetContext(context string) error {
	if len(context) == 0 {
		return errors.New("context cannot be empty")
	}

	ctx, ok := kubeConfig.apiConfig.Contexts[context]
	if !ok {
		return fmt.Errorf("context '%s' not found", context)
	}

	kubeConfig.apiConfig.CurrentContext = context
	if len(ctx.Namespace) > 0 {
		return kubeConfig.SetNamespace(ctx.Namespace)
	}
	return kubeConfig.SetNamespace("default")
}

func (kubeConfig KubeConfig) SetNamespace(namespace string) error {
	if len(namespace) == 0 {
		return errors.New("namespace cannot be empty")
	}

	ctx, ok := kubeConfig.apiConfig.Contexts[kubeConfig.apiConfig.CurrentContext]
	if !ok {
		return errors.New("current context not found")
	}

	ctx.Namespace = namespace
	return nil
}

func (kubeConfig KubeConfig) CurrentContext() (string, error) {
	if len(kubeConfig.apiConfig.CurrentContext) == 0 {
		return "", errors.New("current context not set")
	}
	return kubeConfig.apiConfig.CurrentContext, nil
}

func (kubeConfig KubeConfig) CurrentNamespace() (string, error) {
	ctx, ok := kubeConfig.apiConfig.Contexts[kubeConfig.apiConfig.CurrentContext]
	if !ok {
		return "", errors.New("current context not found")
	}
	return ctx.Namespace, nil
}

func (kubeConfig KubeConfig) Write() error {
	return clientcmd.ModifyConfig(kubeConfig.configAccess, *kubeConfig.apiConfig, true)
}
