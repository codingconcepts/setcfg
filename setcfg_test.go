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
		parts  string
		exp    string
		expErr error
	}{
		{
			name:  "single level with no placeholders",
			input: "a: b",
			parts: "",
			exp:   "a: b\n",
		},
		{
			name:  "single level with no matching placeholders",
			input: "a: b",
			parts: "b: c",
			exp:   "a: b\n",
		},
		{
			name:  "single level with nullifying placeholders",
			input: "a: ~hello~",
			parts: "hello:",
			exp:   "a: null\n",
		},
		{
			name:  "single level with a matching placeholder - string",
			input: "a: ~hello~",
			parts: "hello: hi",
			exp:   "a: hi\n",
		},
		{
			name:  "single level with a matching placeholder - list",
			input: "a: ~hello~",
			parts: "hello:\n- 1\n- 2",
			exp:   "a:\n- 1\n- 2\n",
		},
		{
			name:  "single level with a matching placeholder - map",
			input: "a: ~hello~",
			parts: "hello:\n  one: 1\n  two: 2",
			exp:   "a:\n  one: 1\n  two: 2\n",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			input, err := parse(strings.NewReader(c.input))
			if err != nil {
				t.Fatalf("error parsing input: %v", err)
			}

			parts, err := parse(strings.NewReader(c.parts))
			if err != nil {
				t.Fatalf("error parsing parts: %v", err)
			}

			err = setParsed(input, parts)
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
