package pola_test

import (
	"fmt"
	"io"
	"net/url"
	"testing"

	"github.com/ipsusila/pola"
	"github.com/stretchr/testify/assert"
)

func TestIo(t *testing.T) {
	f := fsSub()
	fd, err := f.Open("sample.json")
	assert.NoError(t, err)
	defer fd.Close()

	n, err := io.Copy(pola.DevNull, fd)
	assert.NoError(t, err)

	fmt.Println("/dev/null Copy:", n)

	// Test closers (with 5 items)
	cs := pola.NewClosers(nil, nil, pola.DevNull, nil, nil)
	assert.Equal(t, cs.Len(), 5)
	assert.NoError(t, cs.Close())
	assert.Equal(t, cs.Remove(nil).Len(), 4)
	assert.Equal(t, cs.Remove(pola.DevNull).Len(), 3)
	assert.Equal(t, cs.Append(io.NopCloser(pola.DevNull)).Len(), 4)

	for cs.Len() > 2 {
		cs.TakeFirst()
	}

	_, ok := cs.TakeFirst()
	assert.True(t, ok)

	_, ok = cs.TakeFirst()
	assert.True(t, ok)
	assert.Equal(t, cs.Len(), 0)

	cs.Append(nil)
	_, ok = cs.TakeLast()
	assert.True(t, ok)
}

func TestParseUrl(t *testing.T) {
	patterns := []string{
		"",
		"abc.txt",
		"abc",
		"/",
		".",
		"./",
		"..",
		"../",
		`C:`,
		`C:\windows`,
		`D:\data\abc`,
		"file:///tmp",
		"tcp://localhost",
		"tcp://127.0.0.1:1000",
		"unix:///temp/try.sock",
		"file://localhost/tmp/test",
	}

	for _, p := range patterns {
		u, err := url.Parse(p)
		if err != nil {
			fmt.Printf("Error: [%s], error: %s\n", p, err.Error())
		} else {
			fmt.Printf("Input: [%s], scheme: `%s`, host: `%s`, port: `%s`, path: %s\n", p, u.Scheme, u.Host, u.Port(), u.Path)
		}
	}
}
