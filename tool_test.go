package pola_test

import (
	"fmt"
	"go/build"
	"testing"

	"github.com/ipsusila/pola"
	"github.com/k0kubun/pp/v3"
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

func TestBuild(t *testing.T) {
	ctx := build.Default
	dirs := ctx.SrcDirs()
	for _, dir := range dirs {
		fmt.Println("Source: ", dir)
	}
	fmt.Println("GOPATH:", ctx.GOPATH)
	fmt.Println("GOROOT:", ctx.GOROOT)
	fmt.Println("GOARCH:", ctx.GOARCH)
	fmt.Println("GOOS  :", ctx.GOOS)

	pkg, err := build.Import("io/fs", ".", build.FindOnly)
	assert.NoError(t, err)
	//fmt.Printf("%#v\n", pkg)
	pp.Println(pkg)
}
