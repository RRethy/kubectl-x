package ctx

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	"k8s.io/client-go/tools/clientcmd/api"

	fzf "github.com/RRethy/kubectl-x/internal/fzf/testing"
	history "github.com/RRethy/kubectl-x/internal/history/testing"
	kubeconfig "github.com/RRethy/kubectl-x/internal/kubeconfig/testing"
	kubernetes "github.com/RRethy/kubectl-x/internal/kubernetes/testing"
)

func TestCtxer_Ctx(t *testing.T) {
	tests := []struct {
		name              string
		initialContext    string
		initialNamespace  string
		selectedContext   string
		selectedNamespace string
		expectedOut       string
		err               bool
	}{
		{
			name:              "switches context and namespace successfully",
			initialContext:    "fo",
			initialNamespace:  "ba",
			selectedContext:   "foobar",
			selectedNamespace: "baz",
			expectedOut:       "Switched to context \"foobar\".\nSwitched to namespace \"baz\".\n",
		},
		{
			name:              "returns error when selecting context fails",
			initialContext:    "fo",
			initialNamespace:  "ba",
			selectedContext:   "",
			selectedNamespace: "baz",
			err:               true,
		},
		{
			name:              "returns error when selecting namespace fails",
			initialContext:    "fo",
			initialNamespace:  "ba",
			selectedContext:   "foobar",
			selectedNamespace: "",
			err:               true,
		},
		{
			name:              "switches to context and namespace from history",
			initialContext:    "-",
			initialNamespace:  "",
			selectedContext:   "",
			selectedNamespace: "",
			expectedOut:       "Switched to context \"old-bar\".\nSwitched to namespace \"old-ns-bar\".\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			history := &history.FakeHistory{Data: map[string][]string{"context": {"old-foo", "old-bar", "old-baz"}}}
			err := Ctxer{
				kubeconfig.NewFakeKubeConfig(
					map[string]*api.Context{
						"old-foo": {Cluster: "old-foo", Namespace: "old-ns-foo"},
						"old-bar": {Cluster: "old-bar", Namespace: "old-ns-bar"},
					},
					test.selectedContext,
					test.selectedNamespace,
				),
				genericiooptions.IOStreams{Out: out},
				kubernetes.NewFakeClient(map[string][]any{
					"namespace": {
						&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "foo"}},
						&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "bar"}},
						&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "baz"}},
					},
				}),
				fzf.NewFakeFzf([]fzf.InputOutput{
					{Input: test.initialContext, Output: test.selectedContext},
					{Input: test.initialNamespace, Output: test.selectedNamespace},
				}),
				history,
			}.Ctx(context.Background(), test.initialContext, test.initialNamespace)

			if test.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expectedOut, out.String())
			}
		})
	}
}
