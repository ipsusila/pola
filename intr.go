package pola

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// RunnerFunc is adapter to allow function to be used as Runner
type RunnerFunc func(context.Context) error

// Runner is the interface that wrap a runnable task/function/action
type Runner interface {
	RunContext(ctx context.Context) error
}

func (f RunnerFunc) RunContext(ctx context.Context) error {
	return f(ctx)
}

// InterruptibleFunc execute runner with given function
func InterruptibleFunc(fn func(ctx context.Context) error, sigs ...os.Signal) error {
	var rf RunnerFunc = fn
	return InterruptibleContext(context.Background(), rf, sigs...)
}

// Interruptible execute Runner and cancel it when SIGINT is captured.
// Here, context.Background is used as the parent context.
// By default, the function listen to os.Interrupt and syscall.SIGINT,
// if additional signals are needed, pass them to optional `sigs` argument.
func Interruptible(r Runner, sigs ...os.Signal) error {
	return InterruptibleContext(context.Background(), r, sigs...)
}

// InterruptibleContext execute Runner and cancel it when SIGINT is captured.
// The argument `ctx` is used as the parent context for constructing cancelable context.
// By default, the function listen to os.Interrupt and syscall.SIGINT,
// if additional signals are needed, pass them to optional `sigs` argument.
func InterruptibleContext(ctx context.Context, r Runner, sigs ...os.Signal) error {
	cctx, cancel := context.WithCancel(ctx)
	defer cancel()

	chSigs := make(chan os.Signal, 1)
	chDone := make(chan bool, 1)
	chRunDone := make(chan bool, 1)

	// Handle several cases:
	// 1. Signal retrieved
	// 2. App/task done
	// 3. Canceled by other through parent context
	notif := []os.Signal{os.Interrupt, syscall.SIGINT}
	notif = append(notif, sigs...)
	signal.Notify(chSigs, notif...)
	go func() {
		defer close(chDone)
		select {
		case <-chSigs:
			cancel()
			return
		case <-chRunDone:
			return
		case <-cctx.Done():
			return
		}
	}()

	err := r.RunContext(cctx)
	close(chRunDone)

	// wait until go routine done
	<-chDone

	return err
}
