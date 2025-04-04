package pola_test

import (
	"testing"

	"github.com/ipsusila/pola"
	"github.com/stretchr/testify/assert"
)

func TestDecoder(t *testing.T) {
	names := []string{
		"sample.jsonnet",
		"sample.yml",
		"sample.toml",
		"sample.json",
	}

	f := fsSub()
	for _, name := range names {
		var dest map[string]any
		dec := pola.NewFsDecoder(name, f)
		err := dec.Decode(&dest)
		assert.NoError(t, err)
		printJson(dest)
	}
}
