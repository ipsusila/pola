package es

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"os"
	"runtime"

	"github.com/ipsusila/pola"
)

// Special variables name
const (
	SpvRuntime = "__runtime__"
	SpvHost    = "__host__"
)

var (
	ErrScriptNotSpecified = errors.New("script not specified")
)

type Executor interface {
	ExecuteContext(ctx context.Context, argv ...any) (any, error)
	Execute(argv ...any) (any, error)
	Define(symbol string, val any) error
	SetScript(s string, name ...string) error
	SetScriptFile(name string, fd ...fs.FS) error
}

// ExecutorOption define function for setting up executor
type ExecutorOption func(e Executor)

// AnyArray convert array of T to array of any
func AnyArray[T any](args []T) []any {
	argv := []any{}
	for _, arg := range args {
		argv = append(argv, arg)
	}
	return argv
}

type Host struct {
	OS      string
	Arch    string
	NumCpu  int
	Name    string
	Version string
}

func MakeHost() Host {
	name, _ := os.Hostname()
	return Host{
		OS:      runtime.GOOS,
		Arch:    runtime.GOARCH,
		NumCpu:  runtime.NumCPU(),
		Name:    name,
		Version: runtime.Version(),
	}
}

type Runtime struct {
	// inputs to script
	Ctx  context.Context
	Argv []any
	Cwd  string

	// return value
	Err error
	Ret any

	cls pola.Closers
}

func MakeRuntime() Runtime {
	return Runtime{
		cls: pola.NewClosers(),
	}
}
func (r *Runtime) AddCloser(c io.Closer) io.Closer {
	if c == nil {
		return pola.DevNull
	}
	sc := pola.SafeSyncCloser(c)
	r.cls.Append(sc)

	return sc
}
func (r *Runtime) Cleanup() error {
	return r.cls.Close()
}
func (r *Runtime) ReadMemStats() *runtime.MemStats {
	ms := runtime.MemStats{}
	runtime.ReadMemStats(&ms)

	return &ms
}
