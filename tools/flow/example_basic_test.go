package flow_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"cuelang.org/go/cue"
	"cuelang.org/go/tools/flow"
)

func TestSimpleFlow(t *testing.T) {
	var r cue.Runtime
	inst, err := r.Compile("example.cue", `
	a: {
		input: "world"
		output: string
	}
	b: {
		input: a.output
		output: string
	}
	`)
	if err != nil {
		log.Fatal(err)
	}
	controller := flow.New(nil, inst, ioTaskFunc)
	rv, err := controller.Run(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	aout, _ := rv.LookupPath(cue.ParsePath("a.output")).String()
	assert.Equal(t, "hello world", aout)
	bout, _ := rv.LookupPath(cue.ParsePath("b.output")).String()
	assert.Equal(t, "hello hello world", bout)
}

func ioTaskFunc(v cue.Value) (flow.Runner, error) {
	inputPath := cue.ParsePath("input")

	input := v.LookupPath(inputPath)
	if !input.Exists() {
		return nil, nil
	}

	return flow.RunnerFunc(func(t *flow.Task) error {
		inputVal, err := t.Value().LookupPath(inputPath).String()
		if err != nil {
			return fmt.Errorf("input not of type string")
		}

		outputVal := fmt.Sprintf("hello %s", inputVal)
		// Output:
		// setting a.output to "hello world"
		// setting b.output to "hello hello world"
		// fmt.Printf("setting %s.output to %q\n", t.Path(), outputVal)

		return t.Fill(map[string]string{
			"output": outputVal,
		})
	}), nil
}
