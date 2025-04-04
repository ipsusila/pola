package ank

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/ipsusila/pola/es"
	"github.com/mattn/anko/ast"
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
