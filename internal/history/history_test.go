package history

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHistory(t *testing.T) {
	tests := []struct {
		name     string
		contents string
		data     map[string][]string
		err      bool
	}{
		{
			name: "parses history file successfully",
			contents: `
data:
  context:
  - context1
  - context2
  namespace:
  - namespace1
  - namespace2
`,
			data: map[string][]string{
				"context":   {"context1", "context2"},
				"namespace": {"namespace1", "namespace2"},
			},
			err: false,
		},
		{
			name:     "returns error when parsing history file fails",
			contents: `\`,
			err:      true,
		},
		{
			name:     "returns no error when history file does not exist",
			contents: "",
			data:     map[string][]string(nil),
			err:      false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tempDir, err := os.MkdirTemp("", "kubectl-x-testing")
			require.NoError(t, err)
			historyPath := filepath.Join(tempDir, "history.yaml")
			os.WriteFile(historyPath, []byte(test.contents), 0o644)
			defer os.RemoveAll(tempDir)

			h, err := NewHistory(NewConfig(WithHistoryPath(historyPath)))
			if test.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.data, h.Data)
			}
		})
	}
}

func TestHistory_Get(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string][]string
		group    string
		distance int
		expected string
		err      bool
	}{
		{
			name:     "returns item from history",
			group:    "context",
			distance: 1,
			data: map[string][]string{
				"context":   {"context1", "context2"},
				"namespace": {"namespace1", "namespace2"},
			},
			expected: "context2",
			err:      false,
		},
		{
			name:  "returns error when group does not exist",
			group: "context",
			data:  map[string][]string{},
			err:   true,
		},
		{
			name:     "returns error when distance is out of bounds",
			group:    "context",
			distance: 2,
			data: map[string][]string{
				"context": {"context1", "context2"},
			},
			err: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &History{Data: test.data}
			item, err := h.Get(test.group, test.distance)
			if test.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expected, item)
			}
		})
	}
}

func TestHistory_Add(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string][]string
		group    string
		item     string
		expected map[string][]string
	}{
		{
			name: "adds item to history",
			data: map[string][]string{
				"context":   {"context1", "context2"},
				"namespace": {"namespace1", "namespace2"},
			},
			group: "context",
			item:  "context3",
			expected: map[string][]string{
				"context":   {"context3", "context1"},
				"namespace": {"namespace1", "namespace2"},
			},
		},
		{
			name:  "adds to empty group",
			data:  map[string][]string{},
			group: "context",
			item:  "context1",
			expected: map[string][]string{
				"context": {"context1"},
			},
		},
		{
			name: "does not go over max history size",
			data: map[string][]string{
				"context": {"context1", "context2"},
			},
			group: "context",
			item:  "context101",
			expected: map[string][]string{
				"context": {"context101", "context1"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := &History{Data: test.data}
			h.Add(test.group, test.item)
			assert.Equal(t, test.expected, h.Data)
		})
	}
}

func TestHistory_Write(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string][]string
		expected string
	}{
		{
			name: "writes history to file",
			data: map[string][]string{
				"context":   {"context1", "context2"},
				"namespace": {"namespace1", "namespace2"},
			},
			expected: `data:
  context:
  - context1
  - context2
  namespace:
  - namespace1
  - namespace2
`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tempDir, err := os.MkdirTemp("", "kubectl-x-testing")
			require.NoError(t, err)
			historyPath := filepath.Join(tempDir, "history.yaml")
			defer os.RemoveAll(tempDir)

			h := &History{Data: test.data, path: historyPath}
			require.NoError(t, err)
			err = h.Write()
			require.NoError(t, err)

			contents, err := os.ReadFile(historyPath)
			require.NoError(t, err)
			assert.Equal(t, test.expected, string(contents))
		})
	}
}
