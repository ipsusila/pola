//go:build ignore

package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/ipsusila/pola/es/ank"
)

func main() {
	/*
		prefix := ""
		imps := []string{"time"}
		dir := `/opt/homebrew/Cellar/go/1.24.1/libexec/src/time/`
		targetPkg := "pkg"
		targetFile := "pkg/time.go"
	*/
	var (
		fPrefix = flag.String("prefix", "", "package prefix, e.g. database")
		fImp    = flag.String("imports", "", "import file, e.g. database/sql")
		fDir    = flag.String("path", "", "directory or file")
		fPkg    = flag.String("pkg", "pkg", "target package")
		fOut    = flag.String("out", "pkg/_out.go", "output file name")
	)

	imps := strings.FieldsFunc(*fImp, func(r rune) bool {
		return r == ',' || r == ';'
	})
	st, err := os.Lstat(*fDir)
	if err != nil {
		log.Fatalln(err)
	}

	res := ank.NewInspectionResult()
	if st.IsDir() {
		if err := ank.InspectDir(res, *fDir, *fPrefix); err != nil {
			log.Fatalln(err)
		}
	} else {
		if err := ank.InspectGoFile(res, *fDir, *fPrefix); err != nil {
			log.Fatalln(err)
		}
	}
	res.Sort()

	fd, err := os.Create(*fOut)
	if err != nil {
		log.Fatalln(err)
	}
	defer fd.Close()

	res.FPrint(fd, *fPkg, imps...)
}
