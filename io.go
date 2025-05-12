package pola

import (
	"errors"
	"io"
	"io/fs"
	"net"
	"net/url"
	"os"
	"strings"
	"sync"
)

const (
	IoStdout  = "<stdout>"
	IoStderr  = "<stderr>"
	IoStdin   = "<stdin>"
	IoNull    = "<null>"
	IoDevNull = "/dev/null"
	IoEmpty   = ""
)

// CurrentDirFS return fs.FS for current working directory
func CurrentDirFS() (fs.FS, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return os.DirFS(pwd), nil
}

// WriteCloser
type nopWCloser struct {
	io.Writer
}

// NopWriteCloser create writer with Close method
func NopWriteCloser(w io.Writer) io.WriteCloser {
	if _, ok := w.(io.ReaderFrom); ok {
		return nopWCloserReadFrom{w}
	}
	return nopWCloser{w}
}

func (nopWCloser) Close() error {
	return nil
}

// WriteCloser + ReaderFrom
type nopWCloserReadFrom struct {
	io.Writer
}

func (nopWCloserReadFrom) Close() error {
	return nil
}
func (n nopWCloserReadFrom) ReadFrom(r io.Reader) (int64, error) {
	return n.Writer.(io.ReaderFrom).ReadFrom(r)
}

// ReadWriteCloser
type nopRWCloser struct {
	io.ReadWriter
}

func (nopRWCloser) Close() error {
	return nil
}

// NopReadWriteCloser wrap io.ReadWriter with nop Close method.
func NopReadWriteCloser(rw io.ReadWriter) io.ReadWriteCloser {
	_, rf := rw.(io.ReaderFrom)
	_, wt := rw.(io.WriterTo)
	if rf && wt {
		return nopRWCloserReadFromWriteTo{rw}
	} else if rf {
		return nopRWCloserReadFrom{rw}
	} else if wt {
		return nopRWCloserWriteTo{rw}
	}
	return nopRWCloser{rw}
}

type nopRWCloserReadFrom struct {
	io.ReadWriter
}

func (nopRWCloserReadFrom) Close() error {
	return nil
}
func (n nopRWCloserReadFrom) ReadFrom(r io.Reader) (int64, error) {
	return n.ReadWriter.(io.ReaderFrom).ReadFrom(r)
}

type nopRWCloserWriteTo struct {
	io.ReadWriter
}

func (nopRWCloserWriteTo) Close() error {
	return nil
}
func (n nopRWCloserWriteTo) WriteTo(w io.Writer) (int64, error) {
	return n.ReadWriter.(io.WriterTo).WriteTo(w)
}

type nopRWCloserReadFromWriteTo struct {
	io.ReadWriter
}

func (nopRWCloserReadFromWriteTo) Close() error {
	return nil
}
func (n nopRWCloserReadFromWriteTo) WriteTo(w io.Writer) (int64, error) {
	return n.ReadWriter.(io.WriterTo).WriteTo(w)
}
func (n nopRWCloserReadFromWriteTo) ReadFrom(r io.Reader) (int64, error) {
	return n.ReadWriter.(io.ReaderFrom).ReadFrom(r)
}

// DevNull mimics /dev/null behaviour
// It discard on write, and return EOF on read.
var (
	DevNull     = devNull{}
	poolDevNull = NewBytesPool()
)

type devNull struct{}

func (devNull) ReadByte() (byte, error) {
	return 0, io.EOF
}
func (devNull) WriteByte(c byte) error {
	return nil
}
func (devNull) ReadRune() (rune, int, error) {
	return 0, 0, io.EOF
}
func (devNull) WriteString(s string) (int, error) {
	return len(s), nil
}
func (devNull) Write(p []byte) (int, error) {
	return len(p), nil
}
func (devNull) Read(p []byte) (int, error) {
	return 0, io.EOF
}
func (devNull) WriteTo(w io.Writer) (int64, error) {
	return 0, io.EOF
}
func (devNull) ReadFrom(r io.Reader) (n int64, err error) {
	bp := poolDevNull.Get().(*[]byte)
	szRd := 0
	for {
		szRd, err = r.Read(*bp)
		n += int64(szRd)
		if err != nil {
			poolDevNull.Put(bp)
			if errors.Is(err, io.EOF) {
				return n, nil
			}
			return
		}
	}
}
func (devNull) Close() error {
	return nil
}

// NewBytesPool return pool with pre-allocated []byte
// Arg `cap` is opitonal, and if not specified,
// the capacity is 8192 bytes.
func NewBytesPool(cap ...int) *sync.Pool {
	nc := 8192
	if len(cap) > 0 && cap[0] > 0 {
		nc = cap[0]
	}
	return &sync.Pool{
		New: func() any {
			b := make([]byte, nc)
			return &b
		},
	}
}

// ReadCloserFromDescriptor create reader with closer from given descriptor
func ReadCloserFromDescriptor(desc string) (io.ReadCloser, error) {
	ldesc := strings.ToLower(desc)
	switch ldesc {
	case IoEmpty, IoNull, IoDevNull:
		return devNull{}, nil
	case IoStdin:
		return io.NopCloser(os.Stdin), nil
	default:
		return rwclFromDescriptor(desc, false)
	}
}

// WriteCloserFromDescriptor return io.WriteCloser from given descriptr.
// Valid descriptor are: <null>, <stdout>, <stderr> and "desc" as filename
func WriteCloserFromDescriptor(desc string) (io.WriteCloser, error) {
	ldesc := strings.ToLower(desc)
	switch ldesc {
	case IoEmpty, IoNull, IoDevNull:
		return devNull{}, nil
	case IoStdout:
		return NopWriteCloser(os.Stdout), nil
	case IoStderr:
		return NopWriteCloser(os.Stderr), nil
	default:
		return rwclFromDescriptor(desc, true)
	}
}

func rwclFromDescriptor(desc string, wr bool) (io.ReadWriteCloser, error) {
	// try parse
	u, err := url.Parse(desc)
	if err == nil {
		// create based on scheme
		scheme := strings.ToLower(u.Scheme)
		switch scheme {
		case "tcp", "tcp4", "tcp6", "udp", "udp4", "udp6":
			return net.Dial(scheme, u.Host)
		case "unix", "unixgram":
			return net.Dial(scheme, u.Path)
		case "file":
			desc = u.Path
		}
	}

	// open/create file
	if wr {
		return os.Create(desc)
	}
	return os.Open(desc)
}

// Closers hold list of io.Closer
type Closers interface {
	io.Closer
	Append(io.Closer) Closers
	Remove(io.Closer) Closers
	TakeFirst() (io.Closer, bool)
	TakeLast() (io.Closer, bool)
	Empty() bool
	Len() int
	Clear()
}

type safeCloser struct {
	c io.Closer
}

func SafeCloser(c io.Closer) io.Closer {
	return &safeCloser{c}
}

func (sc *safeCloser) Close() error {
	if c := sc.c; c != nil {
		sc.c = nil
		return c.Close()
	}
	return nil
}

type safeSyncCloser struct {
	sync.Mutex
	c io.Closer
}

func SafeSyncCloser(c io.Closer) io.Closer {
	return &safeCloser{c: c}
}

func (sc *safeSyncCloser) Close() error {
	sc.Lock()
	defer sc.Unlock()

	if c := sc.c; c != nil {
		sc.c = nil
		return c.Close()
	}
	return nil
}

type closers struct {
	sync.RWMutex
	items []io.Closer
}

func NewClosers(c ...io.Closer) Closers {
	return &closers{items: c}
}

func (cs *closers) Empty() bool {
	cs.RLock()
	defer cs.RUnlock()

	return len(cs.items) == 0
}
func (cs *closers) Len() int {
	cs.RLock()
	defer cs.RUnlock()

	return len(cs.items)
}
func (cs *closers) Close() error {
	cs.Lock()
	defer cs.Unlock()

	var errs error
	for _, c := range cs.items {
		if c != nil {
			errs = errors.Join(errs, c.Close())
		}
	}
	cs.items = nil

	return errs
}
func (cs *closers) Append(c io.Closer) Closers {
	cs.Lock()
	defer cs.Unlock()

	cs.items = append(cs.items, c)
	return cs
}
func (cs *closers) TakeFirst() (io.Closer, bool) {
	cs.Lock()
	defer cs.Unlock()

	if len(cs.items) == 0 {
		return nil, false
	}
	c := cs.items[0]
	cs.items = cs.items[1:]
	return c, true
}
func (cs *closers) TakeLast() (io.Closer, bool) {
	cs.Lock()
	defer cs.Unlock()

	if len(cs.items) == 0 {
		return nil, false
	}
	ne := len(cs.items) - 1
	c := cs.items[ne]
	cs.items = cs.items[:ne]

	return c, true
}
func (cs *closers) Remove(c io.Closer) Closers {
	cs.Lock()
	defer cs.Unlock()

	pos := -1
	for i, ci := range cs.items {
		if ci == c {
			pos = i
			break
		}
	}
	if pos == -1 {
		return cs
	}

	n := len(cs.items)
	for i := pos + 1; i < n; i++ {
		cs.items[i-1] = cs.items[i]
	}
	cs.items = cs.items[:n-1]

	return cs
}
func (cs *closers) Clear() {
	cs.Lock()
	cs.items = nil
	cs.Unlock()
}

// check file exists
func PathExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}
