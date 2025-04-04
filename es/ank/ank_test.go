package ank_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/ipsusila/pola/es"
	"github.com/ipsusila/pola/es/ank"
	"github.com/ipsusila/pola/es/ank/pkg"
	"github.com/mattn/anko/vm"
	"github.com/stretchr/testify/assert"

	_ "github.com/mattn/anko/packages"
	_ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
)

type Descriptor struct {
	Name string
	Desc string
}

var pkgs = map[string]map[string]reflect.Value{
	"custom": {
		"RFC3339": reflect.ValueOf(time.RFC3339),
	},
}

var types = map[string]map[string]reflect.Type{
	"custom": {
		"Descriptor": reflect.TypeOf(&Descriptor{}),
	},
}

func TestScript(t *testing.T) {
	dir := "../../_data"
	files := []string{
		"new.ank",
		"vars.ank",
		"argv.ank",
		"now.ank",
		"sync.ank",
		"chan.ank",
		"anonym.ank",
		"sql.ank",
	}

	imp := vm.NewPackagesWithStdImporter(pkgs, types, "fmt", "time", "sync").
		Append(vm.NewPackagesImporter(pkg.Pkgs, pkg.PkgTypes))
	ex := ank.NewAnkoExecutor(
		ank.Debug(),
		ank.WithCorePackages(),
		ank.WithWorkingDirectory(dir),
		ank.WithPackagesImporter(imp),
	)
	for _, name := range files {
		ex.SetScriptFile(name)
		res, err := ex.Execute(es.AnyArray(files)...)
		assert.NoError(t, err)
		assert.Nil(t, res)
	}
}

func TestFuncs(t *testing.T) {
	dir := "../../_data"
	file := "funcs.ank"

	imp := vm.NewStdPackagesImporter("fmt", "time", "sync").
		Append(vm.NewPackagesImporter(pkg.Pkgs, pkg.PkgTypes))
	ex := ank.NewAnkoExecutor(
		ank.Debug(),
		ank.WithCorePackages(),
		ank.WithWorkingDirectory(dir),
		ank.WithScriptFile(file),
		ank.WithPackagesImporter(imp),
	)
	res, err := ex.Execute()
	assert.Nil(t, res)
	assert.NoError(t, err)
}
