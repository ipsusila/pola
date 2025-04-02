package pola_test

import (
	"embed"
	"fmt"

	"github.com/TylerBrock/colorjson"
)

//go:embed _data
var dataFs embed.FS

func printJson(data any) {
	f := colorjson.NewFormatter()
	f.Indent = 2

	s, _ := f.Marshal(data)
	fmt.Println(string(s))
}
