package main

import (
	"encoding/base64"
	"github.com/stretchr/testify/assert"
	"github.com/wojnosystems/envcrypt/pkg/envfile"
	"testing"
)

func TestDecryptAES256(t *testing.T) {
	key := make([]byte, 32)
	key, _ = base64.StdEncoding.DecodeString("MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI=")
	cases := map[string]struct {
		input envfile.EnvUnsealed
	}{
		"simple": {
			input: envfile.EnvUnsealed{"a": "b"},
		},
	}

	for caseName, c := range cases {
		sealed, err := encryptAES256(key, c.input)
		assert.NoError(t, err, caseName)
		actual, err := decryptAES256(key, []envfile.EnvSealed{sealed})
		assert.NoError(t, err, caseName)
		assert.Equal(t, c.input, actual, caseName)
	}
}
