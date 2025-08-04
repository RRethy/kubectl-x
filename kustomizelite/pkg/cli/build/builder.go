package build

import (
	"bytes"
	"context"
	"fmt"

	"gopkg.in/yaml.v3"
	"k8s.io/cli-runtime/pkg/genericiooptions"

	"github.com/RRethy/utils/kustomizelite/pkg/kustomize"
)

type Builder struct {
	IOStreams  genericiooptions.IOStreams
	kustomizer kustomize.Kustomizer
}

func (b *Builder) Build(_ context.Context, path string) error {
	resources, err := b.kustomizer.Kustomize(path, nil)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	for i, resource := range resources {
		var buf bytes.Buffer
		encoder := yaml.NewEncoder(&buf)
		encoder.SetIndent(2)
		if err := encoder.Encode(resource); err != nil {
			return fmt.Errorf("marshaling resource %d: %w", i, err)
		}

		if i > 0 {
			_, _ = fmt.Fprintln(b.IOStreams.Out, "---")
		}
		_, _ = fmt.Fprint(b.IOStreams.Out, buf.String())
	}

	return nil
}
