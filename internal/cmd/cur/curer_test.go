package cur

import (
	"bytes"
	"context"
	"testing"

	kubeconfig "github.com/RRethy/kubectl-x/internal/kubeconfig/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/cli-runtime/pkg/genericiooptions"
)

func TestCurer_Cur(t *testing.T) {
	tests := []struct {
		name             string
		currentContext   string
		currentNamespace string
		expectedOut      string
		err              bool
	}{
		{
			name:             "prints context and namespace successfully",
			currentContext:   "foobar",
			currentNamespace: "baz",
			expectedOut:      "Current context: \"foobar\"\nCurrent namespace: \"baz\"\n",
			err:              false,
		},
		{
			name:             "returns error when getting current context fails",
			currentContext:   "",
			currentNamespace: "baz",
			err:              true,
		},
		{
			name:             "returns error when getting current namespace fails",
			currentContext:   "foobar",
			currentNamespace: "",
			err:              true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			err := Curer{
				KubeConfig: kubeconfig.NewFakeKubeConfig(nil, test.currentContext, test.currentNamespace),
				IoStreams:  genericiooptions.IOStreams{Out: out},
			}.Cur(context.Background())
			if test.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expectedOut, out.String())
			}
		})
	}
}
