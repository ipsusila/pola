package ank_test

import (
	"testing"

	"github.com/ipsusila/pola/es/ank"
	"github.com/stretchr/testify/assert"
)

func TestGenStdPackages(t *testing.T) {
	goVer := "1.23.0"
	rootPkgDir := "pkg"
	tgtPkgName := "std"
	err := ank.GenerateStdPackages(goVer, rootPkgDir, tgtPkgName)
	assert.NoError(t, err)
}

func TestGenCurrent(t *testing.T) {
	pkgVer := ""
	pkgImport := "github.com/ipsusila/pola/es/ank"
	pkgSrcDir := "."
	tgtOutPkgDir := "tst"
	tgtOutPkg := "usr"
	rootDirs := []string{}
	err := ank.GenerateCustomPackages(pkgVer, pkgImport, pkgSrcDir, tgtOutPkgDir, tgtOutPkg, rootDirs...)
	assert.NotNil(t, err)
}
