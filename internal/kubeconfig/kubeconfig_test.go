package kubeconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/tools/clientcmd/api"
)

func TestKubeConfig_Contexts(t *testing.T) {
	kubeConfig := KubeConfig{
		apiConfig: &api.Config{
			Contexts: map[string]*api.Context{
				"context1": {},
				"context2": {},
			},
		},
	}

	contexts := kubeConfig.Contexts()
	assert.ElementsMatch(t, []string{"context1", "context2"}, contexts)
}

func TestKubeConfig_UseContextWithNamespace(t *testing.T) {
	kubeConfig := KubeConfig{
		apiConfig: &api.Config{
			Contexts: map[string]*api.Context{
				"context1": {
					Namespace: "namespace1",
				},
				"context2": {
					Namespace: "namespace2",
				},
			},
		},
	}

	err := kubeConfig.UseContextWithNamespace("context1", "namespace2")
	require.Nil(t, err)
	assert.Equal(t, "context1", kubeConfig.apiConfig.CurrentContext)
	assert.Equal(t, "namespace2", kubeConfig.apiConfig.Contexts["context1"].Namespace)
}

func TestKubeConfig_UseContext(t *testing.T) {
	kubeConfig := KubeConfig{
		apiConfig: &api.Config{
			Contexts: map[string]*api.Context{
				"context1": {},
				"context2": {
					Namespace: "namespace2",
				},
			},
		},
	}

	err := kubeConfig.UseContext("context1")
	require.Nil(t, err)
	assert.Equal(t, "context1", kubeConfig.apiConfig.CurrentContext)
	assert.Equal(t, "default", kubeConfig.apiConfig.Contexts["context1"].Namespace)

	err = kubeConfig.UseContext("context2")
	require.Nil(t, err)
	assert.Equal(t, "context2", kubeConfig.apiConfig.CurrentContext)
	assert.Equal(t, "namespace2", kubeConfig.apiConfig.Contexts["context2"].Namespace)
}

func TestKubeConfig_UseNamespace(t *testing.T) {
	kubeConfig := KubeConfig{
		apiConfig: &api.Config{
			Contexts: map[string]*api.Context{
				"context1": {},
			},
		},
	}

	err := kubeConfig.UseContext("context1")
	require.Nil(t, err)
	err = kubeConfig.UseNamespace("namespace1")
	require.Nil(t, err)
	assert.Equal(t, "namespace1", kubeConfig.apiConfig.Contexts["context1"].Namespace)
}
