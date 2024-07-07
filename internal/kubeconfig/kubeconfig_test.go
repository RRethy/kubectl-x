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

	err := kubeConfig.SetContext("context1")
	require.Nil(t, err)
	assert.Equal(t, "context1", kubeConfig.apiConfig.CurrentContext)
	assert.Equal(t, "default", kubeConfig.apiConfig.Contexts["context1"].Namespace)

	err = kubeConfig.SetContext("context2")
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

	err := kubeConfig.SetContext("context1")
	require.Nil(t, err)
	err = kubeConfig.SetNamespace("namespace1")
	require.Nil(t, err)
	assert.Equal(t, "namespace1", kubeConfig.apiConfig.Contexts["context1"].Namespace)
}

func TestKubeConfig_CurrentContext(t *testing.T) {
	tests := []struct {
		name      string
		apiConfig *api.Config
		expected  string
		err       bool
		errMsg    string
	}{
		{
			name: "correct context when set",
			apiConfig: &api.Config{
				CurrentContext: "context1",
				Contexts: map[string]*api.Context{
					"context1": {},
					"context2": {},
				},
			},
			expected: "context1",
			err:      false,
			errMsg:   "",
		},
		{
			name: "error when context not set",
			apiConfig: &api.Config{
				CurrentContext: "",
				Contexts: map[string]*api.Context{
					"context1": {},
					"context2": {},
				},
			},
			expected: "",
			err:      true,
			errMsg:   "current context not set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, err := KubeConfig{apiConfig: tt.apiConfig}.CurrentContext()
			if tt.err {
				require.NotNil(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				require.Nil(t, err)
				assert.Equal(t, tt.expected, context)
			}
		})
	}
}

func TestKubeConfig_CurrentNamespace(t *testing.T) {
	kubeConfig := KubeConfig{
		apiConfig: &api.Config{
			CurrentContext: "context1",
			Contexts: map[string]*api.Context{
				"context1": {
					Namespace: "namespace1",
				},
			},
		},
	}

	namespace, err := kubeConfig.CurrentNamespace()
	require.Nil(t, err)
	assert.Equal(t, "namespace1", namespace)
}
