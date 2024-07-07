package kubernetes

import (
	"context"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/kubectl/pkg/scheme"
)

type Interface interface {
	List(ctx context.Context, resourceType string) ([]any, error)
}

type Client struct {
	ConfigFlags *genericclioptions.ConfigFlags
}

func NewClient(configFlags *genericclioptions.ConfigFlags) Interface {
	return &Client{ConfigFlags: configFlags}
}

func (c *Client) List(ctx context.Context, resourceType string) ([]any, error) {
	infos, err := resource.NewBuilder(c.ConfigFlags).
		WithScheme(scheme.Scheme, scheme.Scheme.PrioritizedVersionsAllGroups()...).
		// NamespaceParam("default").
		// FieldSelectorParam(c.FieldSelector).
		// LabelSelectorParam(c.LabelSelector).
		ContinueOnError().
		ResourceTypeOrNameArgs(true, string(resourceType)).
		Flatten().
		Do().
		Infos()
	if err != nil {
		return nil, err
	}

	var res []any
	for _, info := range infos {
		res = append(res, info.Object)
	}
	return res, nil
}
