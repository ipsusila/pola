package pola_test

import (
	"path/filepath"
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
	pth, err := filepath.Abs("_data/sample.json")
	assert.NoError(t, err)

	err = pola.NewFsDecoder(pth).Decode(&dst)
	assert.NoError(t, err)

	// different approach
	clear(dst)
	pth = "_data/sample.yml"
	err = pola.NewFsDecoder(pth).Decode(&dst)
	assert.NoError(t, err)

	// cast type
	clear(dst)
	err = pola.FormattedTextFile(pth).Decode(&dst)
	assert.NoError(t, err)

	clear(dst)
	js := `{"debug": true, "packages": ["fmt", "io", "os"]}`
	dec := pola.JsonText(js)
	err = dec.Decode(&dst)
	assert.NoError(t, err)
}
