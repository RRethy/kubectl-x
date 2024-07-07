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

func WithPipeReader(r *io.PipeReader) FzfOption {
	return func(f *Fzf) {
		f.pipeReader = r
	}
}

func WithPipeWriter(w *io.PipeWriter) FzfOption {
	return func(f *Fzf) {
		f.pipeWriter = w
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
	pipeReader *io.PipeReader
	pipeWriter *io.PipeWriter
	exactMatch bool
	sorted     bool
}

func NewFzf(opts ...FzfOption) *Fzf {
	pipeReader, pipeWriter := io.Pipe()
	fzf := &Fzf{
		exec:       exec.New(),
		ioStreams:  genericiooptions.IOStreams{},
		pipeReader: pipeReader,
		pipeWriter: pipeWriter,
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

	go func() {
		defer f.pipeWriter.Close()
		var filteredItems []string
		for _, item := range items {
			if strings.Contains(item, initialSearch) {
				filteredItems = append(filteredItems, item)
			}
		}
		if f.sorted {
			sort.Strings(filteredItems)
		}
		if _, err := fmt.Fprint(f.pipeWriter, strings.Join(filteredItems, "\n")); err != nil {
			panic(err)
		}
	}()

	cmd.SetStdin(f.pipeReader)

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
