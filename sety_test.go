package main

import (
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
			name:  "single level file without placeholders",
			input: `a: b`,
			parts: ``,
			exp:   `a: b`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			inputReader := strings.NewReader(c.input)
			partsReader := strings.NewReader(c.parts)

			output, err := set(inputReader, partsReader)
			test.Equals(t, c.expErr, err)
			if err != nil {
				return
			}

			test.Equals(t, c.exp, output)
		})
	}
}
