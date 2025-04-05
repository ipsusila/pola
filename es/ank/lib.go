package ank

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ipsusila/pola"
)

func createTopPkgDef(outDir, tgtPkgName string) error {
	fname := tgtPkgName + "_lib.go"
	fname = filepath.Join(outDir, fname)
	fd, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer fd.Close()

	fmt.Fprintf(fd, `// Auto generated file
package %s

import (
	"reflect"

	"github.com/mattn/anko/vm"
)

var (
	Pkgs     = make(map[string]map[string]reflect.Value)
	PkgTypes = make(map[string]map[string]reflect.Type)
)

func NewImporter(pkgs ...string) vm.Importer {
	imp := vm.NewPackagesImporter(nil, nil)
	return imp.AppendMap(Pkgs, PkgTypes, pkgs...)
}
`, tgtPkgName)

	return nil
}

func generatePackagesFromHints(hints *pola.GoPackageHints, outDir, tgtPkgName string) error {
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return err
	}

	// cteate top package definition
	if err := createTopPkgDef(outDir, tgtPkgName); err != nil {
		return err
	}

	for _, h := range hints.Hints {
		res := NewInspectionResult()
		if err := InspectDir(res, h); err == nil {
			tp := filepath.Join(outDir, h.OutputFilename())
			if werr := res.WriteFile(tp, tgtPkgName, h.ImportPath); werr != nil {
				return werr
			}
		}
	}

	return nil
}

func GenerateStdPackages(goVer, rootPkgDir, tgtPkgName string, rootDir ...string) error {
	hints, err := pola.GetStdGoPackageHints(goVer, rootDir...)
	if err != nil {
		return err
	}
	if goVer == "" {
		goVer = hints.GoVersion
	}

	vDir := "go" + goVer
	outDir := filepath.Join(rootPkgDir, vDir, tgtPkgName)
	return generatePackagesFromHints(hints, outDir, tgtPkgName)
}

func GenerateCustomPackages(version, importPath, pkgSrcDir, rootPkgDir, tgtPkgName string, rootDir ...string) error {
	hints, err := pola.GetGoPackageHints(version, importPath, pkgSrcDir, rootDir...)
	if err != nil {
		return err
	}

	outDir := filepath.Join(rootPkgDir, tgtPkgName)
	return generatePackagesFromHints(hints, outDir, tgtPkgName)

}
