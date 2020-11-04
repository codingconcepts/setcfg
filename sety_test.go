package main

import (
	"errors"
	"strings"
	"testing"

	"github.com/codingconcepts/sety/internal/pkg/test"
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
			inputReader := strings.NewReader(c.input)
			partsReader := strings.NewReader(c.parts)

			output, err := set(inputReader, partsReader)
			test.Equals(t, c.expErr, errors.Unwrap(err))
			if err != nil {
				return
			}

			test.Equals(t, c.exp, output)
		})
	}
}
