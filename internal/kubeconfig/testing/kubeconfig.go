package testing

import (
	"errors"

	"github.com/RRethy/kubectl-x/internal/kubeconfig"
)

var _ kubeconfig.Interface = &FakeKubeConfig{}

type FakeKubeConfig struct {
	contexts         []string
	currentContext   string
	currentNamespace string
}

func NewFakeKubeConfig(contexts []string, currentContext, currentNamespace string) *FakeKubeConfig {
	return &FakeKubeConfig{
		contexts,
		currentContext,
		currentNamespace,
	}
}

func (fake *FakeKubeConfig) Contexts() []string {
	return fake.contexts
}

func (fake *FakeKubeConfig) SetContext(context string) error {
	if context == "" {
		return errors.New("context cannot be empty")
	}
	fake.currentContext = context
	return nil
}

func (fake *FakeKubeConfig) SetNamespace(namespace string) error {
	if namespace == "" {
		return errors.New("namespace cannot be empty")
	}
	fake.currentNamespace = namespace
	return nil
}

func (fake *FakeKubeConfig) GetCurrentContext() (string, error) {
	if fake.currentContext == "" {
		return "", errors.New("current context not set")
	}
	return fake.currentContext, nil
}

func (fake *FakeKubeConfig) GetCurrentNamespace() (string, error) {
	if fake.currentNamespace == "" {
		return "", errors.New("current namespace not set")
	}
	return fake.currentNamespace, nil
}

func (fake *FakeKubeConfig) Write() error {
	return nil
}
