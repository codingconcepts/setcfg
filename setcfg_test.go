package main

import (
	"errors"
	"strings"
	"testing"

	"github.com/codingconcepts/setcfg/internal/pkg/test"
	"gopkg.in/yaml.v2"
)

func TestSet(t *testing.T) {
	cases := []struct {
		name   string
		input  string
		env    string
		exp    string
		expErr error
	}{
		{
			name:  "single level with no placeholders",
			input: "a: b",
			env:   "",
			exp:   "a: b\n",
		},
		{
			name:  "single level with no matching placeholders",
			input: "a: b",
			env:   "b: c",
			exp:   "a: b\n",
		},
		{
			name:  "single level with nullifying placeholders",
			input: "a: ~hello~",
			env:   "hello:",
			exp:   "a: null\n",
		},
		{
			name:  "single level with a matching placeholder - string",
			input: "a: ~hello~",
			env:   "hello: hi",
			exp:   "a: hi\n",
		},
		{
			name:  "single level with a matching placeholder - list",
			input: "a: ~hello~",
			env:   "hello:\n- 1\n- 2",
			exp:   "a:\n- 1\n- 2\n",
		},
		{
			name:  "single level with a matching placeholder - map",
			input: "a: ~hello~",
			env:   "hello:\n  one: 1\n  two: 2",
			exp:   "a:\n  one: 1\n  two: 2\n",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			input, err := parse(strings.NewReader(c.input))
			if err != nil {
				t.Fatalf("error parsing input: %v", err)
			}

			env, err := parse(strings.NewReader(c.env))
			if err != nil {
				t.Fatalf("error parsing env: %v", err)
			}

			err = setParsed(input, env)
			test.Equals(t, c.expErr, errors.Unwrap(err))
			if err != nil {
				return
			}

			act, err := yaml.Marshal(input)
			if err != nil {
				t.Fatalf("error marshalling output: %v", err)
			}

			test.Equals(t, c.exp, string(act))
		})
	}
}
