package history

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

const (
	maxHistoryItems = 2
)

var (
	_ Interface = &History{}

	defaultHistoryPath = filepath.Join(os.ExpandEnv("$HOME"), ".local", "share", "kubectl-x", "history.yaml")
)

type Interface interface {
	Get(group string, distance int) (string, error)
	Add(group, item string)
	Write() error
}

type ConfigOption func(*Config)

func WithHistoryPath(path string) ConfigOption {
	return func(config *Config) {
		config.historyPath = path
	}
}

type Config struct {
	historyPath string
}

func NewConfig(options ...ConfigOption) *Config {
	config := &Config{historyPath: defaultHistoryPath}
	for _, option := range options {
		option(config)
	}
	return config
}

type History struct {
	Data map[string][]string `json:"data"`

	path string
}

func NewHistory(config *Config) (*History, error) {
	contents, err := os.ReadFile(config.historyPath)
	history := History{path: config.historyPath}
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("reading file: %s", err)
	} else if err == nil {
		err = yaml.Unmarshal(contents, &history)
		if err != nil {
			return nil, fmt.Errorf("unmarshalling history: %s", err)
		}
	}
	return &history, nil
}

func (h *History) Get(group string, distance int) (string, error) {
	if h.Data == nil {
		return "", fmt.Errorf("no history found")
	}

	groupHistory, ok := h.Data[group]
	if !ok {
		return "", fmt.Errorf("group '%s' not found in history", group)
	}

	if distance >= len(groupHistory) {
		return "", fmt.Errorf("unable to go back %d items in history", distance)
	}

	return groupHistory[distance], nil
}

func (h *History) Add(group, item string) {
	if h.Data == nil {
		h.Data = make(map[string][]string)
	}

	h.Data[group] = append([]string{item}, h.Data[group]...)
	if len(h.Data[group]) > maxHistoryItems {
		h.Data[group] = h.Data[group][:maxHistoryItems]
	}
}

func (h *History) Write() error {
	contents, err := yaml.Marshal(h)
	if err != nil {
		return fmt.Errorf("marshalling history: %s", err)
	}

	err = os.MkdirAll(filepath.Dir(h.path), 0o755)
	if err != nil {
		return fmt.Errorf("creating directory: %s", err)
	}

	err = os.WriteFile(h.path, contents, 0o644)
	if err != nil {
		return fmt.Errorf("writing file: %s", err)
	}

	return nil
}
