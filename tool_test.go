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
	pkgImport := "github.com/BurntSushi/toml"
	pkgSrcDir := ""
	hints2, err := pola.GetGoPackageHints(tgtVer, pkgImport, pkgSrcDir)
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

func TestGetDir(t *testing.T) {
	gopath := pola.GoPath()
	goroot := pola.GoRoot()
	fmt.Println("Path:", gopath, "> Root:", goroot)
}
