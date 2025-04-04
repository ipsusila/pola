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

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Ident struct {
	Name     string
	Exported bool
	HasParam bool
	Filename string
}

type InspectionResult struct {
	Prefix  string // prefix given
	PkgName string // package name
	Structs []Ident
	Funcs   []Ident
	Consts  []Ident
	Vars    []Ident

	mVals map[string]bool
	mTyps map[string]bool
}

func NewInspectionResult() *InspectionResult {
	return &InspectionResult{
		mVals: make(map[string]bool),
		mTyps: make(map[string]bool),
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
	return "val" + r.prefixAndName()
}
func (r *InspectionResult) TypeName() string {
	return "typ" + r.prefixAndName()
}
func (r *InspectionResult) prefixAndName() string {
	camel := cases.Title(language.Und, cases.NoLower)
	return camel.String(r.Prefix) + camel.String(r.PkgName)
}
func (r *InspectionResult) FPrint(w io.Writer, targetPkg string, impName ...string) {
	//vals := strings.Builder{}
	prtIdent := func(w io.Writer, comment, kind, brace string, idts []Ident) {
		if len(idts) == 0 {
			return
		}

		hasItem := false
		sb := strings.Builder{}
		for _, id := range idts {
			if !id.Exported || id.HasParam {
				continue
			}
			item := fmt.Sprintf(`    "%s": reflect.%s(%s.%s%s),`, id.Name, kind, r.PkgName, id.Name, brace)
			fmt.Fprintln(&sb, item)
			hasItem = true
		}
		if hasItem {
			fmt.Fprintln(w, "    //"+comment)
			fmt.Fprint(w, sb.String())
		}
	}
	vals := strings.Builder{}
	prtIdent(&vals, "Function(s)", "ValueOf", "", r.Funcs)
	prtIdent(&vals, "Variable(s)", "ValueOf", "", r.Vars)
	prtIdent(&vals, "Constant(s)", "ValueOf", "", r.Consts)
	if vals.Len() == 0 {
		return
	}

	// generate importStmt
	impStmt := ""
	for i, imp := range impName {
		if i == 0 {
			impStmt += fmt.Sprintln()
		}
		impStmt += fmt.Sprintln(`    "` + imp + `"`)
	}

	// add values
	const tplVals = `// Auto-generated code.
package %s

import (
    "reflect"

	%s
)

var %s = map[string]reflect.Value{
%s}
`
	fmt.Fprintf(w, tplVals, targetPkg, impStmt, r.ValueName(), vals.String())

	// Add types
	typs := strings.Builder{}
	prtIdent(&typs, "Struct(s)", "TypeOf", "{}", r.Structs)

	// format string
	const tplTyps = `
var %s = map[string]reflect.Type{
%s}
`
	//strTyps := fmt.Sprintf(t)
	fmt.Fprintf(w, tplTyps, r.TypeName(), typs.String())
}

func InspectDir(res *InspectionResult, dir, prefix string) error {
	absDir, err := filepath.Abs(dir)
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
			err := InspectGoFile(res, filepath.Join(absDir, e.Name()), prefix)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func InspectGoFile(res *InspectionResult, srcPath, prefix string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, srcPath, nil, 0)
	if err != nil {
		return err
	}

	// package name
	if name := f.Name; name != nil {
		if name.Name == "main" {
			return nil
		}
		res.PkgName = name.Name
	}
	res.Prefix = prefix

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
