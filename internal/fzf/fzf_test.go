package fzf

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	"k8s.io/utils/exec"
	fakeexec "k8s.io/utils/exec/testing"
)

func TestFzf_Run(t *testing.T) {
	fcmd := &fakeexec.FakeCmd{
		RunScript: []fakeexec.FakeAction{
			func() ([]byte, []byte, error) { return []byte("bar\n"), nil, nil },
		},
	}
	fexec := &fakeexec.FakeExec{
		CommandScript: []fakeexec.FakeCommandAction{
			func(cmd string, args ...string) exec.Cmd { return fakeexec.InitFakeCmd(fcmd, cmd, args...) },
		},
	}
	ioStreams := genericiooptions.IOStreams{In: bytes.NewReader([]byte("")), Out: bytes.NewBuffer([]byte("")), ErrOut: bytes.NewBuffer([]byte(""))}
	pipeReader, pipeWriter := io.Pipe()
	fzf := NewFzf(WithExec(fexec), WithIOStreams(ioStreams), WithPipeReader(pipeReader), WithPipeWriter(pipeWriter))
	selected, err := fzf.Run([]string{"foo", "bar", "baz"})
	require.NoError(t, err)
	assert.Equal(t, "bar", selected)
}
