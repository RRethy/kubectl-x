package kustomize

import (
	"github.com/RRethy/utils/kustomizelite/pkg/exec"
	"github.com/RRethy/utils/kustomizelite/pkg/helm"
)

// Option is a functional option for configuring a Kustomizer.
type Option func(*kustomization)

// WithHelmTemplater sets a custom helm templater (useful for testing).
func WithHelmTemplater(templater helm.Templater) Option {
	return func(k *kustomization) {
		k.helmTemplater = templater
	}
}

// WithExecWrapper sets a custom exec wrapper for running commands.
func WithExecWrapper(wrapper exec.Wrapper) Option {
	return func(k *kustomization) {
		k.execWrapper = wrapper
	}
}
