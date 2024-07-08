package ns

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericiooptions"

	fzf "github.com/RRethy/kubectl-x/internal/fzf/testing"
	history "github.com/RRethy/kubectl-x/internal/history/testing"
	kubeconfig "github.com/RRethy/kubectl-x/internal/kubeconfig/testing"
	kubernetes "github.com/RRethy/kubectl-x/internal/kubernetes/testing"
)

func TestNser_Ns(t *testing.T) {
	tests := []struct {
		name        string
		initialNs   string
		selectedNs  string
		expectedOut string
		err         bool
	}{
		{
			name:        "switches namespace successfully",
			initialNs:   "fo",
			selectedNs:  "foobar",
			expectedOut: "Switched to namespace \"foobar\".\n",
		},
		{
			name:       "returns error when selecting namespace fails",
			initialNs:  "fo",
			selectedNs: "",
			err:        true,
		},
		{
			name:        "switches to namespace from history",
			initialNs:   "-",
			selectedNs:  "old-foo",
			expectedOut: "Switched to namespace \"old-foo\".\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			history := &history.FakeHistory{Data: map[string][]string{"namespace": {"old-foo", "old-bar", "old-baz"}}}
			err := Nser{
				KubeConfig: kubeconfig.NewFakeKubeConfig(nil, "foobar", test.selectedNs),
				IoStreams:  genericiooptions.IOStreams{Out: out},
				K8sClient: kubernetes.NewFakeClient(map[string][]any{
					"namespace": {
						&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "foo"}},
						&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "bar"}},
						&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "baz"}},
					},
				}),
				Fzf: fzf.NewFakeFzf([]fzf.InputOutput{
					{Input: test.initialNs, Output: test.selectedNs},
				}),
				History: history,
			}.Ns(context.Background(), test.initialNs)

			if test.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.selectedNs, history.Data["namespace"][0])
				assert.Equal(t, test.expectedOut, out.String())
			}
		})
	}
}
