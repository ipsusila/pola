package pola_test

import (
	"embed"
	"fmt"
	"io/fs"

	"github.com/TylerBrock/colorjson"
)

//go:embed _data
var fsData embed.FS

func fsSub(name ...string) fs.FS {
	nm := "_data"
	if len(name) > 0 && name[0] != "" {
		nm = name[0]
	}
	f, _ := fs.Sub(fsData, nm)
	return f
}

func printJson(data any) {
	f := colorjson.NewFormatter()
	f.Indent = 2

	s, _ := f.Marshal(data)
	fmt.Println(string(s))
}
