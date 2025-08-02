package testing

import "fmt"

type InputOutput struct {
	Input  string
	Output string
}

type FakeFzf struct {
	runScript      []InputOutput
	runScriptIndex int
}

func NewFakeFzf(runScript []InputOutput) *FakeFzf {
	return &FakeFzf{runScript, 0}
}

func (f *FakeFzf) Run(initialSearch string, items []string) (string, error) {
	if f.runScriptIndex < len(f.runScript) {
		inputOutput := f.runScript[f.runScriptIndex]
		f.runScriptIndex++
		if inputOutput.Input == initialSearch {
			return inputOutput.Output, nil
		}
		panic(fmt.Sprintf("expected input %s but got %s", inputOutput.Input, initialSearch))
	}
	panic("not enough items in run script")
}
