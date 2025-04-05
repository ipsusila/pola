package pola

import (
	"errors"
	"fmt"
	"go/build"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/mod/modfile"
)

var (
	ErrInvalidSemVer = errors.New("invalid Semantic Version")
	reSemVer         = regexp.MustCompile(`[0-9]{1,}\.[0-9]{1,}\.[0-9]{1,}`)
	reCharBeg        = regexp.MustCompile(`^[0-9a-zA-Z\!\~]`)
)

type GoPackageHint struct {
	Domain       string
	RelPath      string
	ImportPath   string
	SrcDir       string
	IsStdPkg     bool
	ModGoVersion string
	Version      string
	HasGoMod     bool
	Valid        bool
}

type GoPackageHints struct {
	BaseDir   string
	GoVersion string
	Hints     []*GoPackageHint
}

func (g *GoPackageHint) validVarName(c rune) rune {
	if c == '_' {
		return c
	}
	if '0' <= c && c <= '9' {
		return c
	}
	if 'a' <= c && c <= 'z' {
		return c
	}
	if 'A' <= c && c <= 'Z' {
		return c
	}

	return '_'
}

func (g *GoPackageHint) SanitizeImportPath() {
	up := false
	sb := strings.Builder{}
	for _, c := range g.ImportPath {
		if c == '!' {
			up = true
		} else if up {
			sb.WriteRune(unicode.ToUpper(c))
			up = false
		} else {
			sb.WriteRune(c)
		}
	}
	g.ImportPath = sb.String()
}

func (g *GoPackageHint) ID() string {
	if g.IsStdPkg {
		return g.ImportPath + "@v" + g.Version
	}
	return g.RelPath
}
func (g *GoPackageHint) GoVariableName() string {
	up := true
	sb := strings.Builder{}
	for _, c := range g.RelPath {
		if c == '!' || c == os.PathSeparator {
			up = true
		} else if up {
			up = false
			sb.WriteRune(unicode.ToUpper(g.validVarName(c)))
		} else {
			sb.WriteRune(g.validVarName(c))
		}
	}
	return sb.String()
}

func (g *GoPackageHint) OutputFilename() string {
	ps := '_'
	if g.IsStdPkg {
		ps = '.'
	}
	up := false
	sb := strings.Builder{}
	for _, c := range g.RelPath {
		if c == '@' {
			sb.WriteRune('_')
		} else if c == os.PathSeparator || c == '/' {
			sb.WriteRune(ps)
		} else if c == '!' {
			up = true
		} else if up {
			up = false
			sb.WriteRune(unicode.ToUpper(c))
		} else {
			sb.WriteRune(c)
		}
	}
	if !g.IsStdPkg {
		sb.WriteString("_v")
		sb.WriteString(g.Version)
	}
	sb.WriteString(".go")
	return sb.String()
}

func pkgIgnoreDir(name string) bool {
	if !reCharBeg.MatchString(name) {
		return true
	}
	excludes := []string{"vendor", "internal", "cmd", "go", "testdata", "builtin", "reflect", "syscall", "unsafe"}
	for _, ex := range excludes {
		if name == ex {
			return true
		}
	}
	return false
}

func IsGoSrcDir(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, ent := range entries {
		if !ent.IsDir() {
			if ent.Name() == "go.mod" {
				return true
			}
			if strings.ToLower(filepath.Ext(ent.Name())) == ".go" {
				return true
			}
		}
	}
	return false
}

func GoPath() string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	return gopath
}
func GoRoot() string {
	goroot := os.Getenv("GOROOT")
	if goroot == "" {
		goroot = build.Default.GOROOT
	}
	return goroot
}
func GoVersion() string {
	gover := os.Getenv("GOVERSION")
	if gover == "" {
		gover = reSemVer.FindString(GoRoot())
	}
	return gover
}

func GetGoModHint(base, importDirective, dir string) (*GoPackageHint, error) {
	if !IsGoSrcDir(dir) {
		return nil, nil
	}

	name := "go.mod"
	gmFile := filepath.Join(dir, name)
	data, err := os.ReadFile(gmFile)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	hint := GoPackageHint{
		IsStdPkg: false,
		SrcDir:   dir,
		HasGoMod: err == nil,
	}

	if base == "" {
		hint.Version = "0.0.0"
		hint.ImportPath = importDirective
		items := strings.Split(importDirective, "/")
		if len(items) > 1 {
			hint.RelPath = strings.Join(items[1:], "/")
		} else {
			hint.RelPath = importDirective
		}
	} else if rel, err := filepath.Rel(base, dir); err == nil {
		hint.RelPath = rel

		version := reSemVer.FindString(rel)
		if version != "" {
			hint.Version = version
		}
		if idx := strings.IndexRune(rel, '@'); idx > 0 {
			hint.ImportPath = rel[:idx]
			hint.SanitizeImportPath()
		}
	}

	// parse go mod if exists
	if hint.HasGoMod {
		gm, err := modfile.Parse(gmFile, data, nil)
		if err != nil {
			return nil, err
		}
		//return nil, err
		if gm.Go != nil {
			hint.ModGoVersion = gm.Go.Version
		}
		if gm.Module != nil {
			hint.ImportPath = gm.Module.Mod.Path
		}

		// std path
		// TODO: is correct (?)
		//hint.IsStdPkg = hint.ImportPath == "std"
	}

	return &hint, nil
}

func GetGoPackageHints(tgtVer, tgtImport, pkgSrcDir string, rootDir ...string) (*GoPackageHints, error) {
	// layout > go/pkg/mod/github.com/!burnt!sushi/toml@v1.5.0
	// go.mod > module github.com/BurntSushi/toml
	// no go.mod > /go/pkg/mod/github.com/k0kubun/pp@v3.0.1+incompatible
	goBase, err := filepath.Abs(GoPath())
	if err != nil {
		return nil, err
	}
	if len(rootDir) > 0 && rootDir[0] != "" {
		goBase = rootDir[0]
	}
	// get domain information
	domain := "unknown"
	items := strings.Split(tgtImport, "/")
	if len(items) >= 2 {
		domain = items[0]
	}

	var srcBaseDir, modBase string
	if pkgSrcDir == "" {
		// get srcBaseDir for specific domain
		modBase = filepath.Join(goBase, "pkg", "mod")
		srcBaseDir = filepath.Join(modBase, domain)
	} else {
		modBase = ""
		absPath, err := filepath.Abs(pkgSrcDir)
		if err != nil {
			return nil, err
		}
		srcBaseDir = absPath
	}

	hints := GoPackageHints{
		BaseDir:   srcBaseDir,
		GoVersion: GoVersion(),
	}

	// get matching version and import path only
	wderr := filepath.WalkDir(srcBaseDir, func(path string, d fs.DirEntry, err error) error {
		if d == nil {
			return os.ErrNotExist
		} else if d.Name() == "." || !d.IsDir() {
			return nil
		}

		if pkgIgnoreDir(d.Name()) {
			return filepath.SkipDir
		}

		hint, err := GetGoModHint(modBase, tgtImport, path)
		if err != nil {
			return err
		} else if hint != nil {
			match := hint.ImportPath == tgtImport && (tgtVer == "" || hint.Version == tgtVer)
			if match {
				hint.Domain = domain
				hints.Hints = append(hints.Hints, hint)
			}
			return filepath.SkipDir
		}
		return nil
	})

	return &hints, wderr
}

func GetStdGoPackageHints(version string, rootDir ...string) (*GoPackageHints, error) {
	sdkBase := ""
	if len(rootDir) > 0 && rootDir[0] != "" {
		sdkBase = rootDir[0]
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		sdkBase = home
	}
	var srcBaseDir string
	if version != "" {
		// sanitize version
		// remove non digit
		version = reSemVer.FindString(version)
		if version == "" {
			return nil, ErrInvalidSemVer
		}
		srcBaseDir = filepath.Join(sdkBase, "sdk", "go"+version, "src")
	} else {
		sdkBase = GoRoot()
		version = GoVersion()
		srcBaseDir = filepath.Join(sdkBase, "src")
	}
	fmt.Println("BaseSrc:", srcBaseDir, ">version:", version)

	hints := GoPackageHints{
		BaseDir:   sdkBase,
		GoVersion: version,
	}
	wderr := filepath.WalkDir(srcBaseDir, func(path string, d fs.DirEntry, err error) error {
		if d == nil {
			return os.ErrNotExist
		} else if d.Name() == "." || !d.IsDir() {
			return nil
		}

		if pkgIgnoreDir(d.Name()) {
			return filepath.SkipDir
		}
		rel, err := filepath.Rel(srcBaseDir, path)
		if err != nil {
			return err
		}
		if rel == "." || rel == ".." {
			return nil
		}

		if IsGoSrcDir(path) {
			hint := GoPackageHint{
				RelPath:      rel,
				ImportPath:   rel,
				SrcDir:       path,
				Version:      version,
				IsStdPkg:     true,
				ModGoVersion: "", // TODO
			}
			hints.Hints = append(hints.Hints, &hint)
		}
		return nil
	})

	return &hints, wderr
}
