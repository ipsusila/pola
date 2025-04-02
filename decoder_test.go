package pola_test

import (
	"testing"

	"github.com/ipsusila/pola"
	"github.com/stretchr/testify/assert"
)

func TestDecoder(t *testing.T) {
	names := []string{
		"_data/sample.jsonnet",
		"_data/sample.yml",
		"_data/sample.toml",
		"_data/sample.json",
	}
	for _, name := range names {
		var dest map[string]any
		dec := pola.NewFsDecoder(name, dataFs)
		err := dec.Decode(&dest)
		assert.NoError(t, err)
		printJson(dest)
	}
}
