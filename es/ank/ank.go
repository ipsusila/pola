package ank

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/ipsusila/pola/es"
	"github.com/mattn/anko/ast"
	"github.com/mattn/anko/ast/astutil"
	"github.com/mattn/anko/core"
	"github.com/mattn/anko/env"
	"github.com/mattn/anko/parser"
	"github.com/mattn/anko/vm"
)

/*
	"func":     FUNC,
	"return":   RETURN,
	"var":      VAR,
	"throw":    THROW,
	"if":       IF,
	"for":      FOR,
	"break":    BREAK,
	"continue": CONTINUE,
	"in":       IN,
	"else":     ELSE,
	"new":      NEW,
	"true":     TRUE,
	"false":    FALSE,
	"nil":      NIL,
	"module":   MODULE,
	"try":      TRY,
	"catch":    CATCH,
	"finally":  FINALLY,
	"switch":   SWITCH,
	"case":     CASE,
	"default":  DEFAULT,
	"go":       GO,
	"chan":     CHAN,
	"struct":   STRUCT,
	"make":     MAKE,
	"type":     TYPE,
	"len":      LEN,
	"delete":   DELETE,
	"close":    CLOSE,
	"map":      MAP,
	"import":   IMPORT,
*/

type ankoExecutor struct {
	wd      string
	src     string
	srcName string
	debug   bool
	imp     vm.Importer

	// baseEnvironment
	env  *env.Env
	stmt ast.Stmt
}

// NewAnkoExecutor create Anko embedded script executor
func NewAnkoExecutor(opts ...es.ExecutorOption) es.Executor {
	ae := ankoExecutor{
		env:   env.NewEnv(),
		imp:   nil,
		debug: false,
	}
	for _, op := range opts {
		op(&ae)
	}
	// Setup Host information
	ae.env.Define(es.SpvHost, es.MakeHost())

	return &ae
}

// Debug sets debug option.
// Debug() -> set debug to true
// Debug(true) -> set debug to true
// Debug(false) -> set debug to false
// Debug(false, true) -> set debug to false
func Debug(dbg ...bool) es.ExecutorOption {
	debug := true
	if len(dbg) > 0 {
		debug = dbg[0]
	}
	return func(e es.Executor) {
		if ae, ok := e.(*ankoExecutor); ok {
			ae.debug = debug
		}
	}
}

// WithCorePackages import core packages into vm
func WithCorePackages() es.ExecutorOption {
	return func(e es.Executor) {
		if ae, ok := e.(*ankoExecutor); ok {
			// Import core symbols
			core.Import(ae.env)
		}
	}
}

// WithWorkingDirectory set working directory
func WithWorkingDirectory(dir string) es.ExecutorOption {
	return func(e es.Executor) {
		if ae, ok := e.(*ankoExecutor); ok {
			wd, err := filepath.Abs(dir)
			if err != nil {
				panic("Error setting working directory, error: " + err.Error())
			}
			ae.wd = wd
		}
	}
}

// WithPackagesImporter create executor with given importer
func WithPackagesImporter(imp vm.Importer) es.ExecutorOption {
	return func(e es.Executor) {
		if ae, ok := e.(*ankoExecutor); ok {
			if ae.imp != nil {
				ae.imp.Append(imp)
			} else {
				ae.imp = imp
			}
		}
	}
}

// WithStdPackages create Executor with given package name
func WithStdPackages(pkgs ...string) es.ExecutorOption {
	return func(e es.Executor) {
		if ae, ok := e.(*ankoExecutor); ok {
			imp := vm.NewStdPackagesImporter(pkgs...)
			if ae.imp != nil {
				ae.imp.Append(imp)
			} else {
				ae.imp = imp
			}
		}
	}
}

// WithScript specify script source to be executed
func WithScript(str string) es.ExecutorOption {
	return func(e es.Executor) {
		if err := e.SetScript(str); err != nil {
			panic("SetScript error: " + err.Error())
		}
	}
}

// WithScriptFile specify script file to be executed.
// Optionally, it can specify the fs.FS in where to look the script
func WithScriptFile(name string, fd ...fs.FS) es.ExecutorOption {
	return func(e es.Executor) {
		if err := e.SetScriptFile(name, fd...); err != nil {
			panic("SetScriptFile error: " + err.Error() + ", file: " + name)
		}
	}
}

// WithSymbols define symbels in the executing environment
func WithSymbols(syms map[string]any) es.ExecutorOption {
	return func(e es.Executor) {
		if ae, ok := e.(*ankoExecutor); ok {
			for sym, val := range syms {
				ae.env.Define(sym, val)
			}
		}
	}
}

func (a *ankoExecutor) workingDirectory() string {
	if a.wd == "" {
		wd, err := filepath.Abs(".")
		if err != nil {
			panic("Working directory error: " + err.Error())
		}
		return wd
	}
	return a.wd
}

func (a *ankoExecutor) Execute(argv ...any) (any, error) {
	return a.ExecuteContext(context.Background(), argv...)
}
func (a *ankoExecutor) ExecuteContext(ctx context.Context, argv ...any) (any, error) {
	if a.stmt == nil {
		return nil, es.ErrScriptNotSpecified
	}

	// define host and runtime symbols
	rt := es.MakeRuntime()
	rt.Ctx = ctx
	rt.Argv = argv
	rt.Cwd = a.workingDirectory()
	defer func() {
		rt.Cleanup()
		a.stmt = nil
	}()

	// not sure wether compiled statement is not changed after execution
	// therefore, we need to compile it for every execution.
	stmt := a.stmt
	if stmt == nil {
		s, err := parser.ParseSrc(a.src)
		if err != nil {
			return nil, err
		}
		stmt = s
	}

	// Copy the environment and define runtime variables
	enviro := a.env.Copy()
	enviro.Define(es.SpvRuntime, &rt)

	// setup options
	options := vm.Options{
		Debug:       a.debug,
		PkgImporter: a.imp,
	}
	if _, err := vm.RunContext(ctx, enviro, &options, stmt); err != nil {
		return nil, fmt.Errorf("error while executing script '%s', %w", a.srcName, err)
	}

	// test
	a.callFunc(enviro, "onResponse", "Arg1", 2, "Arg3")
	a.callFunc(enviro, "toInt", "123.45")
	a.callFuncChild(enviro, "onResponse", "MyArg1", 123, time.Now())

	// Get result from script
	return rt.Ret, rt.Err
}
func (a *ankoExecutor) Define(sym string, val any) error {
	return a.env.Define(sym, val)
}
func (a *ankoExecutor) SetScript(s string, name ...string) error {
	stmt, err := parser.ParseSrc(s)
	if err != nil {
		return err
	}
	a.src = s
	a.stmt = stmt

	if len(name) > 0 {
		a.srcName = name[0]
	} else {
		a.srcName = "_noname_"
	}

	return nil
}
func (a *ankoExecutor) SetScriptFile(name string, fd ...fs.FS) error {
	var fl fs.FS
	if len(fd) > 0 && fd[0] != nil {
		fl = fd[0]
	} else {
		fl = os.DirFS(a.workingDirectory())
	}

	f, err := fl.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	return a.SetScript(string(data), name)
}

func (a *ankoExecutor) inspect(stmt ast.Stmt) error {
	cp := 0
	cnp := 0
	err := astutil.Walk(stmt, func(i interface{}) error {
		if p, ok := i.(ast.Pos); ok {
			pp := p.Position()
			fmt.Printf("@%d,%d:", pp.Line, pp.Column)
			cp++
		} else {
			cnp++
		}
		fmt.Println(i, ">", reflect.TypeOf(i))
		return nil
	})
	fmt.Println("#Pos:", cp, ",non-Pos:", cnp, ",error:", err)
	return nil
}
func (a *ankoExecutor) callFunc(e *env.Env, name string, args ...any) (any, error) {
	v, err := e.GetValue(name)
	if err != nil {
		return nil, err
	}
	if v.Kind() != reflect.Func {
		return nil, fmt.Errorf("`%s` is not a function", name)
	}

	fmt.Println("Func:", name, ">", v.Type())
	t := v.Type()

	// input types:
	fmt.Println("Function info, name:", t.Name(),
		",Variadic:", t.IsVariadic(),
		",In:", t.NumIn(),
		",Out:", t.NumOut())
	for i := range t.NumIn() {
		ti := t.In(i)

		fmt.Printf("  -[%d] (%v) > `%v`\n", i, ti, ti.String())
	}
	fmt.Println()

	fmt.Println("Output args")
	for i := range t.NumOut() {
		fmt.Printf("  +[%d] (%v)\n", i, t.Out(i))
	}
	fmt.Println("Is VM Function?", a.isVmFunc(t))
	fmt.Println()

	// call it
	if len(args) != t.NumIn()-1 {
		return nil, nil
	}

	if a.isVmFunc(t) {
		inps := make([]reflect.Value, 0, t.NumIn())
		inps = append(inps, reflect.ValueOf(context.Background()))
		for _, arg := range args {
			rv := reflect.ValueOf(arg)
			inps = append(inps, reflect.ValueOf(rv))
		}
		outs := v.Call(inps)
		fmt.Println("Result: ", outs)
		for _, rs := range outs {
			fmt.Println(">", rs.Interface().(reflect.Value))
		}
		fmt.Println()
	}

	return nil, nil
}

func (a *ankoExecutor) isVmFunc(t reflect.Type) bool {
	ist := reflect.TypeOf([]interface{}{})
	rvt := reflect.TypeOf(reflect.Value{})
	tCtx := reflect.TypeOf((*context.Context)(nil)).Elem()
	if t.NumOut() != 2 || t.NumIn() < 1 || t.In(0) != tCtx || t.Out(0) != rvt || t.Out(1) != rvt {
		return false
	}

	if t.NumIn() > 1 {
		if t.IsVariadic() {
			if t.In(t.NumIn()-1) != ist {
				return false
			}
		} else {
			if t.In(t.NumIn()-1) != rvt {
				return false
			}
		}
		for i := 1; i < t.NumIn()-1; i++ {
			if t.In(i) != rvt {
				return false
			}
		}
	}
	return true
}

func (a *ankoExecutor) callFuncChild(e *env.Env, name string, args ...any) (any, error) {
	// create temporary script and call it
	f, err := e.GetValue(name)
	if err != nil {
		return nil, err
	}
	if f.Kind() != reflect.Func {
		return nil, fmt.Errorf("`%s` is not a function", name)
	}
	narg := len(args)
	ft := f.Type()
	if ft.IsVariadic() {
		return nil, nil
	}

	// is vm function
	if a.isVmFunc(ft) {
		numIn := ft.NumIn() - 1
		if narg < numIn {
			return nil, fmt.Errorf("function required #%d arg, called with #%d arg", numIn, narg)
		}

		selfEnv := e.NewEnv()
		argv := make([]string, 0, numIn)
		for i := range numIn {
			sym := fmt.Sprintf("__a%d__", i)
			selfEnv.Define(sym, args[i])
			argv = append(argv, sym)
		}
		script := fmt.Sprintf("%s(%s)", name, strings.Join(argv, ","))
		ret, err := vm.ExecuteContext(context.Background(), selfEnv, nil, script)
		fmt.Printf("Calling `%s`, result is %v (%T), err: %v\n", name, ret, ret, err)

		if va, ok := ret.([]any); ok {
			for i, r := range va {
				fmt.Printf("  r[%d](%T)>%v\n", i, r, r)
			}
		}
	}

	return nil, nil
}
