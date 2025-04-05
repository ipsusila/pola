package pola_test

import (
	"fmt"
	"testing"

	"github.com/ipsusila/pola"
	"github.com/stretchr/testify/assert"
)

func TestTool(t *testing.T) {
	hints, err := pola.GetStdGoPackageHints("1.23.0")
	assert.NoError(t, err)
	assert.NotEmpty(t, hints)

	tgtVer := "" // any version is fine
	targetPath := "github.com/BurntSushi/toml"
	hints2, err := pola.GetGoPackageHints(tgtVer, targetPath)
	assert.NoError(t, err)

	fmt.Println("=======================")
	for _, h := range hints.Hints {
		fmt.Printf("[%s](%s) > %s\n", h.ImportPath, h.ID(), h.OutputFilename())
	}
	fmt.Println("-----------------------")
	for _, h := range hints2.Hints {
		fmt.Printf("[%s](%s) > %s\n", h.ImportPath, h.ID(), h.OutputFilename())
	}

}
