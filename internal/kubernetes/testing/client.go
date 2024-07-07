package testing

import (
	"context"
	"errors"

	"github.com/RRethy/kubectl-x/internal/kubernetes"
)

var _ kubernetes.Interface = &FakeClient{}

type FakeClient struct {
	resources map[string][]any
}

func NewFakeClient(resources map[string][]any) *FakeClient {
	return &FakeClient{resources}
}

func (fake *FakeClient) List(ctx context.Context, resourceType string) ([]any, error) {
	if resources, ok := fake.resources[resourceType]; ok {
		return resources, nil
	}
	return nil, errors.New("resource type not found")
}
