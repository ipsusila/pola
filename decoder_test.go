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

	// can use full-path
	var dst map[string]any
	path := "/Users/ipsusila/Workspaces/golang/crawler/pola/_data/sample.json"
	err := pola.NewFsDecoder(path).Decode(&dst)
	assert.NoError(t, err)

	// different approach
	clear(dst)
	path = "_data/sample.yml"
	err = pola.NewFsDecoder(path).Decode(&dst)
	assert.NoError(t, err)

	// cast type
	clear(dst)
	err = pola.FormattedTextFile(path).Decode(&dst)
	assert.NoError(t, err)

	clear(dst)
	js := `{"debug": true, "packages": ["fmt", "io", "os"]}`
	dec := pola.JsonText(js)
	err = dec.Decode(&dst)
	assert.NoError(t, err)
}
