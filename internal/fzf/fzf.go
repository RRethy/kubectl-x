package fzf

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"

	"k8s.io/cli-runtime/pkg/genericiooptions"
	"k8s.io/utils/exec"
)

const binaryName = "fzf"

type Interface interface {
	Run(initialSearch string, items []string) (string, error)
}

type FzfOption func(*Fzf)

func WithExec(exec exec.Interface) FzfOption {
	return func(f *Fzf) {
		f.exec = exec
	}
}

func WithIOStreams(ioStreams genericiooptions.IOStreams) FzfOption {
	return func(f *Fzf) {
		f.ioStreams = ioStreams
	}
}

func WithExactMatch(exactMatch bool) FzfOption {
	return func(f *Fzf) {
		f.exactMatch = exactMatch
	}
}

func WithSorted(sorted bool) FzfOption {
	return func(f *Fzf) {
		f.sorted = sorted
	}
}

type Fzf struct {
	exec       exec.Interface
	ioStreams  genericiooptions.IOStreams
	exactMatch bool
	sorted     bool
}

func NewFzf(opts ...FzfOption) *Fzf {
	fzf := &Fzf{
		exec:       exec.New(),
		ioStreams:  genericiooptions.IOStreams{},
		exactMatch: false,
		sorted:     true,
	}
	for _, opt := range opts {
		opt(fzf)
	}
	return fzf
}

func (f *Fzf) Run(initialSearch string, items []string) (string, error) {
	cmd := f.exec.Command(binaryName, f.buildArgs()...)
	pipeReader, pipeWriter := io.Pipe()

	go func() {
		defer pipeWriter.Close()
		var filteredItems []string
		for _, item := range items {
			if strings.Contains(item, initialSearch) {
				filteredItems = append(filteredItems, item)
			}
		}
		if f.sorted {
			sort.Strings(filteredItems)
		}
		if _, err := fmt.Fprint(pipeWriter, strings.Join(filteredItems, "\n")); err != nil {
			panic(err)
		}
	}()

	cmd.SetStdin(pipeReader)

	var out bytes.Buffer
	cmd.SetStdout(&out)
	cmd.SetStderr(f.ioStreams.ErrOut)

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("running fzf: %s", err)
	}

	output := strings.TrimSpace(out.String())
	if output == "" {
		return "", fmt.Errorf("no item selected")
	}
	return output, nil
}

func (f Fzf) buildArgs() []string {
	args := []string{
		"--height",
		"30%",
		"--ansi",
		"--select-1",
		"--exit-0",
		"--color=dark",
		"--layout=reverse",
	}
	if f.exactMatch {
		args = append(args, "--exact")
	}
	return args
}
