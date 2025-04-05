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

	/*

		targetPkg = "usr"
		outDir = filepath.Join("pkg", targetPkg)
		hints2, err := pola.GetGoPackageHints("github.com")
		assert.NoError(t, err)
		for _, h := range hints2.Hints {
			res := ank.NewInspectionResult()
			if err := ank.InspectDir(res, h); err == nil {
				tp := filepath.Join(outDir, h.OutputFilename())
				res.WriteFile(tp, targetPkg, h.ImportPath)
			}
		}
	*/

}
