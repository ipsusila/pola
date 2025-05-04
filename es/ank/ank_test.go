package ank_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/ipsusila/pola/es"
	"github.com/ipsusila/pola/es/ank"
	"github.com/ipsusila/pola/es/ank/pkg/go1.23.0/std"
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

	imp := std.NewImporter("fmt", "time", "sync", "database/sql", "context").
		AppendMap(pkgs, types)
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

	imp := std.NewImporter("fmt", "time", "sync").
		AppendMap(pkgs, types)
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

func TestCallback(t *testing.T) {
	dir := "../../_data"
	file := "callback.ank"

	imp := std.NewImporter("fmt", "time", "sync", "io").
		AppendMap(pkgs, types)
	ex := ank.NewAnkoExecutor(
		ank.Debug(),
		ank.WithCorePackages(),
		ank.WithWorkingDirectory(dir),
		ank.WithScriptFile(file),
		ank.WithPackagesImporter(imp),
	)
	res, err := ex.Execute()
	assert.NotNil(t, res)
	assert.NoError(t, err)
	//pp.Println(res)

	vv := reflect.ValueOf(res)
	//pp.Println(vv)
	//tf := reflect.TypeOf(res)
	//pp.Println(tf)

	args := make([]reflect.Value, 0, 3)
	args = append(args, reflect.ValueOf(context.Background()))

	ra := reflect.ValueOf(int64(11))
	rb := reflect.ValueOf(int64(12))
	args = append(args, reflect.ValueOf(ra))
	args = append(args, reflect.ValueOf(rb))
	rv := vv.Call(args)
	for i, rr := range rv {
		fmt.Printf("%d (%d) > %T > %#v, %s, %T\n", i, len(rv), rr, rr, rr.Kind(), reflect.TypeOf(rr))
	}
	if rv[0].CanInterface() {
		x := rv[0].Interface()
		xv := rv[0].Interface().(reflect.Value)
		fmt.Printf("%T (%#v) > %v [%d]\n", x, x, xv.Kind(), xv.Int())
	}

	// VM Function signature func
	// func(context.Context, reflect.Value, reflect.Value) (reflect.Value, reflect.Value)
	// First argument always a context

	// Test reflection
	n := 100
	vN := reflect.ValueOf(n)
	fmt.Printf("Value: %T> %v (%#v), can set: %v, can iface: %v\n", vN, vN, vN, vN.CanSet(), vN.CanInterface())
	vvN := reflect.ValueOf(vN)
	fmt.Printf("Value: %T> %v (%#v), can set: %v, can iface: %v\n", vvN, vvN, vvN, vvN.CanSet(), vvN.CanInterface())

	fmt.Println(vN, ">>", vvN)

	// VM funtions receive:
	// except context, it'is in the form of reflect.ValueOf(reflect.ValueOf(v)),
	// If we assign function to Host Variable which has type of function signature,
	// VM function will be converted into the function signature.

}

func TestModule(t *testing.T) {
	dir := "../../_data"
	file := "module.ank"

	imp := std.NewImporter("fmt", "time", "sync").
		AppendMap(pkgs, types)
	ex := ank.NewAnkoExecutor(
		ank.Debug(),
		ank.WithCorePackages(),
		ank.WithWorkingDirectory(dir),
		ank.WithScriptFile(file),
		ank.WithPackagesImporter(imp),
	)
	res, err := ex.Execute()
	fmt.Println(res)
	assert.NotNil(t, res)
	assert.NoError(t, err)
}

func TestCallFuncs(t *testing.T) {
	dir := "../../_data"
	file := "callback.ank"

	imp := std.NewImporter("fmt", "time", "sync", "io").
		AppendMap(pkgs, types)
	ex := ank.NewAnkoExecutor(
		ank.Debug(),
		ank.WithCorePackages(),
		ank.WithWorkingDirectory(dir),
		ank.WithScriptFile(file),
		ank.WithPackagesImporter(imp),
	)
	res, err := ex.Execute()
	fmt.Println("Result:", res, " error>", err)
}
