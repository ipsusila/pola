package ank_test

import (
	"path/filepath"
	"testing"

	"github.com/ipsusila/pola"
	"github.com/ipsusila/pola/es/ank"
	"github.com/stretchr/testify/assert"
)

/*
func TestInspect(t *testing.T) {
	prefix := "jmoiron"
	dir := `/Users/ipsusila/go/pkg/mod/github.com/jmoiron/sqlx@v1.4.0/`
	res := ank.NewInspectionResult()
	err := ank.InspectDir(res, dir, prefix)
	assert.NoError(t, err)
	// Prints
	fmt.Printf("%s/%s\n", res.Prefix, res.PkgName)

	fnPrint := func(label string, ids []ank.Ident) {
		fmt.Printf("  %s\n", label)
		for _, id := range ids {
			if !id.Exported {
				continue
			}
			fmt.Printf("    %s[%v]\n", id.Name, id.HasParam)
		}
	}
	res.Sort()

	fnPrint("FUNCS", res.Funcs)
	fnPrint("STRUCTS", res.Structs)
	fnPrint("VARS", res.Vars)
	fnPrint("CONST", res.Consts)

	res.Fprint(os.Stdout, "pkg", "github.com/jmoiron/sqlx")
}

func TestDoInspect(t *testing.T) {
	// dir /opt/homebrew/Cellar/go/1.24.1/libexec/src/time

	homeDir, err := os.UserHomeDir()
	assert.NoError(t, err)
	goVersion := "1.23.0"
	outDir := "pkg"
	tgtPkgName := "pkg"
	stdSrcDir := fmt.Sprintf("%s/sdk/go%s/src/", homeDir, goVersion)
	fnInsArg := func(srcDir, pkgPath string) ank.InspectionArg {
		inpPkgs := strings.Split(pkgPath, "/")
		fileName := strings.Join(inpPkgs, ".") + ".go"
		if srcDir == "" {
			items := append([]string{stdSrcDir}, inpPkgs...)
			srcDir = filepath.Join(items...)
		}

		prefix := ""
		if n := len(inpPkgs) - 1; n > 0 {
			cText := cases.Title(language.Und, cases.NoLower)
			for i := 0; i < n; i++ {
				prefix += cText.String(inpPkgs[i])
			}
		}
		return ank.InspectionArg{
			Prefix:     prefix,
			Imports:    []string{pkgPath},
			SrcDir:     srcDir,
			TargetPkg:  tgtPkgName,
			TargetFile: filepath.Join(outDir, fileName),
		}
	}

	stdPkgs := []string{
		"io",
		"io/fs",
		"net",
		"net/http",
		"net/url",
		"net/rpc/jsonrpc",
		"os",
		"os/exec",
		"os/signal",
		"os/user",
		"encoding/json",
		"encoding/xml",
		"context",
	}
	for _, sp := range stdPkgs {
		ig := fnInsArg("", sp)
		fmt.Println(ig.String())

		res := ank.NewInspectionResult()
		err := ank.InspectDir(res, ig.SrcDir, ig.Prefix)
		assert.NoError(t, err)
		res.Sort()

		err = res.WriteFile(ig.TargetFile, ig.TargetPkg, ig.Imports...)
		assert.NoError(t, err)
	}
}
*/

func TestInspect(t *testing.T) {
	hints, err := pola.GetStdGoPackageHints("1.23.0")
	assert.NoError(t, err)
	assert.NotEmpty(t, hints)

	targetPkg := "std"
	outDir := filepath.Join("pkg", targetPkg)
	for _, h := range hints.Hints {
		res := ank.NewInspectionResult()
		if err := ank.InspectDir(res, h); err == nil {
			tp := filepath.Join(outDir, h.OutputFilename())
			res.WriteFile(tp, targetPkg, h.ImportPath)
		}
	}

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

	// get inspecton results

	/*
		fmt.Println("=======================")
		for _, h := range hints.Hints {
			fmt.Printf("[%s](%s) > %s\n", h.ImportPath, h.ID(), h.OutputFilename())
		}
		fmt.Println("-----------------------")
		for _, h := range hints2.Hints {
			fmt.Printf("[%s](%s) > %s\n", h.ImportPath, h.ID(), h.OutputFilename())
		}
	*/

}
