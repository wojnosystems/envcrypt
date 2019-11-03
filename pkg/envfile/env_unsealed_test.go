package envfile

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFromUnsealedFile(t *testing.T) {
	cases := map[string]struct {
		file        string
		expected    EnvUnsealed
		expectedErr error
	}{
		"simple": {
			file: "a=a\nb=b\nx=1\n",
			expected: EnvUnsealed{
				"a": "a",
				"b": "b",
				"x": "1",
			},
		},
		"no-end": {
			file: "a=a\nb=b\nx=1",
			expected: EnvUnsealed{
				"a": "a",
				"b": "b",
				"x": "1",
			},
		},
		"malformed line 1": {
			file:        "a\nb=b\nx=1",
			expected:    EnvUnsealed{},
			expectedErr: newParseError(1, "no equal (=) sign"),
		},
		"malformed line 3": {
			file: "a=\nb=b\nx1",
			expected: EnvUnsealed{
				"a": "",
				"b": "b",
			},
			expectedErr: newParseError(3, "no equal (=) sign"),
		},
	}

	for caseName, c := range cases {
		actual, err := FromUnsealedFile(bytes.NewBufferString(c.file))
		assert.Equal(t, c.expectedErr, err, caseName)
		assert.Equal(t, c.expected, actual, caseName)
	}
}

func TestEnvUnsealed_ToUnsealedFile(t *testing.T) {
	cases := map[string]struct {
		input    EnvUnsealed
		expected []byte
	}{
		"blank single": {
			input: EnvUnsealed{
				"a": "",
			},
			expected: []byte("a=\n"),
		},
		"non-blank single": {
			input: EnvUnsealed{
				"a": "b",
			},
			expected: []byte("a=b\n"),
		},
		"non-blank multiline": {
			input: EnvUnsealed{
				"a": "b",
				"c": "d",
				"e": "f",
				"g": "h",
			},
			expected: []byte("a=b\nc=d\ne=f\ng=h\n"),
		},
	}

	for caseName, c := range cases {
		actual := bytes.Buffer{}
		_, _ = c.input.ToUnsealedFile(&actual)
		assert.Equal(t, c.expected, actual.Bytes(), caseName)
	}
}
