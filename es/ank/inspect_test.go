package ank_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/ipsusila/pola/es/ank"
	"github.com/stretchr/testify/assert"
)

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

	res.FPrint(os.Stdout, "pkg", "github.com/jmoiron/sqlx")
}

func TestDoInspect(t *testing.T) {
	// dir /opt/homebrew/Cellar/go/1.24.1/libexec/src/time
	/*
		prefix := "jmoiron"
		imps := []string{"github.com/jmoiron/sqlx"}
		dir := `~/go/pkg/mod/github.com/jmoiron/sqlx@v1.4.0/`
		targetPkg := "pkg"
		targetFile := "pkg/sqlx.go"
	*/
	prefix := ""
	imps := []string{"database/sql"}
	dir := `/opt/homebrew/Cellar/go/1.24.1/libexec/src/database/sql/`
	targetPkg := "pkg"
	targetFile := "pkg/_sql.go"

	res := ank.NewInspectionResult()
	err := ank.InspectDir(res, dir, prefix)
	assert.NoError(t, err)
	res.Sort()

	fd, err := os.Create(targetFile)
	assert.NoError(t, err)
	defer fd.Close()

	res.FPrint(fd, targetPkg, imps...)
}
