package ank

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/ipsusila/pola"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	KindStruct = iota
	KindFunction
	KindVar
	KindConst
)

type Ident struct {
	Name     string
	Exported bool
	HasParam bool
	Filename string
}

type InspectionResult struct {
	PkgName string // package name
	Structs []Ident
	Funcs   []Ident
	Consts  []Ident
	Vars    []Ident
	Hint    *pola.GoPackageHint

	mVals  map[string]bool
	mTyps  map[string]bool
	cTitle cases.Caser
}

type InspectionArg struct {
	Prefix     string
	SrcDir     string
	Imports    []string
	TargetPkg  string
	TargetFile string
	Version    string
}

/*
func makeStdInspectionArg(version, pkgPath, tgtPkgName, outDir string) ([]*InspectionArg, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	stdSrcDir := fmt.Sprintf("%s/sdk/go%s/src/", homeDir, goVersion)
	inpPkgs := strings.Split(pkgPath, "/")
	fileName := strings.Join(inpPkgs, ".") + ".go"
	if srcDir == "" {
		items := append([]string{stdSrcDir}, inpPkgs...)
		srcDir = filepath.Join(items...)
	}

	prefix := ""
	if n := len(inpPkgs) - 1; n > 0 {
		cText := cases.Title(language.Und, cases.NoLower)
		for i := 0; i < n; i++ {
			prefix += cText.String(inpPkgs[i])
		}
	}
	return &InspectionArg{
		Prefix:     prefix,
		Imports:    []string{pkgPath},
		SrcDir:     srcDir,
		TargetPkg:  tgtPkgName,
		TargetFile: filepath.Join(outDir, fileName),
	}, nil
}
*/

func (i InspectionArg) String() string {
	sb := strings.Builder{}
	fmt.Fprintf(&sb, "Prefix:     %s\n", i.Prefix)
	fmt.Fprintf(&sb, "SrcDir:     %s\n", i.SrcDir)
	fmt.Fprintf(&sb, "Imports:    %v\n", i.Imports)
	fmt.Fprintf(&sb, "TargetPkg:  %s\n", i.TargetPkg)
	fmt.Fprintf(&sb, "TargetFile: %s\n", i.TargetFile)
	fmt.Fprintf(&sb, "Version:    %s\n", i.Version)

	return sb.String()
}

func NewInspectionResult() *InspectionResult {
	return &InspectionResult{
		mVals:  make(map[string]bool),
		mTyps:  make(map[string]bool),
		cTitle: cases.Title(language.Und, cases.NoLower),
	}
}

func (r *InspectionResult) Sort() {
	sort.Slice(r.Structs, func(i, j int) bool {
		return r.Structs[i].Name < r.Structs[j].Name
	})
	sort.Slice(r.Funcs, func(i, j int) bool {
		return r.Funcs[i].Name < r.Funcs[j].Name
	})
	sort.Slice(r.Consts, func(i, j int) bool {
		return r.Consts[i].Name < r.Consts[j].Name
	})
	sort.Slice(r.Vars, func(i, j int) bool {
		return r.Vars[i].Name < r.Vars[j].Name
	})
}
func (r *InspectionResult) HasExportedValue() bool {
	for _, v := range r.Funcs {
		if v.Exported && !v.HasParam {
			return true
		}
	}
	for _, v := range r.Vars {
		if v.Exported && !v.HasParam {
			return true
		}
	}
	for _, v := range r.Consts {
		if v.Exported && !v.HasParam {
			return true
		}
	}
	return false
}
func (r *InspectionResult) HasValueId(name string) bool {
	if r.mVals != nil {
		return r.mVals[name]
	}
	return false
}
func (r *InspectionResult) AddValueId(name string) {
	if r.mVals != nil {
		r.mVals[name] = true
	}
}

func (r *InspectionResult) HasExportedType() bool {
	for _, v := range r.Structs {
		if v.Exported && !v.HasParam {
			return true
		}
	}
	return false
}
func (r *InspectionResult) HasTypesId(name string) bool {
	if r.mTyps != nil {
		return r.mTyps[name]
	}
	return false
}
func (r *InspectionResult) AddTypesId(name string) {
	if r.mTyps != nil {
		r.mTyps[name] = true
	}
}
func (r *InspectionResult) ValueName() string {
	pfx := "val"
	if r.Hint != nil {
		pfx += r.Hint.GoVariableName()
	}
	return pfx
}
func (r *InspectionResult) TypeName() string {
	pfx := "typ"
	if r.Hint != nil {
		pfx += r.Hint.GoVariableName()
	}
	return pfx
}
func (r *InspectionResult) MaxNameLen(ids []Ident) int {
	n := 0
	for _, id := range ids {
		if id.Exported && !id.HasParam {
			n = max(n, len(id.Name))
		}
	}
	return n
}
func (r *InspectionResult) printIdent(w io.Writer, comment, kind, brace string, idts []Ident) {
	if len(idts) == 0 {
		return
	}

	hasItem := false
	sb := strings.Builder{}
	nl := r.MaxNameLen(idts)
	for _, id := range idts {
		if !id.Exported || id.HasParam {
			continue
		}
		sp := strings.Repeat(" ", nl-len(id.Name)+1)
		item := fmt.Sprintf(`    "%s":%sreflect.%s(%s.%s%s),`, id.Name, sp, kind, r.PkgName, id.Name, brace)
		fmt.Fprintln(&sb, item)
		hasItem = true
	}
	if hasItem {
		fmt.Fprintln(w, "    //"+comment)
		fmt.Fprint(w, sb.String())
	}
}

func (r *InspectionResult) WriteFile(name string, targetPkg string, impName ...string) error {
	if !(r.HasExportedValue() || r.HasExportedType()) {
		return nil
	}
	fd, err := os.Create(name)
	if err != nil {
		return err
	}
	defer fd.Close()

	r.doFprint(fd, targetPkg, impName...)
	return nil
}
func (r *InspectionResult) Fprint(w io.Writer, targetPkg string, impName ...string) {
	if r.HasExportedValue() || r.HasExportedType() {
		r.doFprint(w, targetPkg, impName...)
	}
}
func (r *InspectionResult) doFprint(w io.Writer, targetPkg string, impName ...string) {
	fmt.Fprintln(w, "// Auto generated code, at", time.Now().Format(time.RFC3339))
	fmt.Fprintln(w, "package", targetPkg)
	fmt.Fprintln(w)
	fmt.Fprintln(w, "import (")
	fmt.Fprintln(w, `    "reflect"`)
	for _, imp := range impName {
		fmt.Fprintf(w, `    "%s"`, imp)
		fmt.Fprintln(w)
	}
	fmt.Fprintln(w, ")")
	fmt.Fprintln(w)

	// print value
	valName := ""
	typName := ""
	hasVal := r.HasExportedValue()
	if hasVal {
		valName = r.ValueName()
		fmt.Fprintf(w, "var %s = map[string]reflect.Value{\n", valName)
		r.printIdent(w, "Function(s)", "ValueOf", "", r.Funcs)
		r.printIdent(w, "Variables(s)", "ValueOf", "", r.Vars)
		r.printIdent(w, "Constants(s)", "ValueOf", "", r.Consts)
		fmt.Fprintln(w, "}")
	}

	// print types
	hasTyp := r.HasExportedType()
	if hasTyp {
		if hasVal {
			fmt.Fprintln(w)
		}
		typName = r.TypeName()
		fmt.Fprintf(w, "var %s = map[string]reflect.Type{\n", typName)
		r.printIdent(w, "Struct(s)", "TypeOf", "{}", r.Structs)
		fmt.Fprintln(w, "}")
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, "func init() {")
	if hasVal {
		fmt.Fprintf(w, `	Pkgs["%s"] = %s`, r.Hint.ImportPath, valName)
		fmt.Fprintln(w)
	}
	if hasTyp {
		fmt.Fprintf(w, `	PkgTypes["%s"] = %s`, r.Hint.ImportPath, typName)
		fmt.Fprintln(w)
	}
	fmt.Fprintln(w, "}")
}

func InspectDir(res *InspectionResult, hint *pola.GoPackageHint) error {
	absDir, err := filepath.Abs(hint.SrcDir)
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(absDir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		loName := strings.ToLower(e.Name())
		if strings.HasSuffix(loName, "_test.go") {
			continue
		}
		ext := filepath.Ext(loName)
		if ext == ".go" {
			err := InspectGoFile(res, hint, filepath.Join(absDir, e.Name()))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func InspectGoFile(res *InspectionResult, hint *pola.GoPackageHint, srcPath string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, srcPath, nil, 0)
	if err != nil {
		return err
	}

	// package name
	if name := f.Name; name != nil {
		if name.Name == "main" || name.Name == "internal" || strings.HasSuffix(name.Name, "_test") {
			return nil
		}
		res.PkgName = name.Name
	}
	res.Hint = hint

	// search for: struct, func, var, const
	// print all function calls
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			// Ignore Method declaration
			if x.Recv == nil && x.Name != nil {
				if !res.HasValueId(x.Name.Name) {
					id := Ident{
						Name:     x.Name.Name,
						Exported: x.Name.IsExported(),
						Filename: srcPath,
					}
					if x.Type != nil {
						id.HasParam = x.Type.TypeParams != nil
					}
					res.Funcs = append(res.Funcs, id)
					res.AddValueId(id.Name)
				}
			}
		case *ast.GenDecl:
			switch x.Tok {
			case token.CONST:
				for _, sp := range x.Specs {
					if v, ok := sp.(*ast.ValueSpec); ok {
						for _, id := range v.Names {
							if res.HasValueId(id.Name) {
								continue
							}
							// add const declaration
							idt := Ident{
								Name:     id.Name,
								Exported: id.IsExported(),
								Filename: srcPath,
							}
							res.Consts = append(res.Consts, idt)
							res.AddValueId(idt.Name)
						}
					}
				}
			case token.VAR:
				for _, sp := range x.Specs {
					if v, ok := sp.(*ast.ValueSpec); ok {
						for _, id := range v.Names {
							if res.HasValueId(id.Name) {
								continue
							}
							// add const declaration
							idt := Ident{
								Name:     id.Name,
								Exported: id.IsExported(),
								Filename: srcPath,
							}
							res.Vars = append(res.Vars, idt)
							res.AddValueId(idt.Name)
						}
					}
				}
			case token.TYPE:
				for _, sp := range x.Specs {
					if ts, ok := sp.(*ast.TypeSpec); ok {
						if ts.Name != nil {
							if res.HasTypesId(ts.Name.Name) {
								continue
							}
							if _, st := ts.Type.(*ast.StructType); st {
								idt := Ident{
									Name:     ts.Name.Name,
									Exported: ts.Name.IsExported(),
									HasParam: ts.TypeParams != nil,
									Filename: srcPath,
								}
								res.Structs = append(res.Structs, idt)
								res.AddTypesId(idt.Name)
							}
						}
					}
				}
			}
		case *ast.BlockStmt:
			return false
		}
		return true
	})

	return nil
}
