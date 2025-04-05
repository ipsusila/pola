//go:build ignore

// 1. Usage: go run genpkg.go -std-ver="1.23.8" -out-pkg-dir="tst" -out-pkg-name="std"
// This will generate Anko packages for go version 1.23.8 under "tst" dir, with package name of `std`
// Don't forget to install corresponding Go version. See https://go.dev/doc/manage-install
//
// 2. Usage: go run genpkg.go -pkg-ver="" -pkg-import="github.com/BurntSushi/toml" -out-pkg-dir="tst" -out-pkg-name="usr"
// This wil generate Anko package for library specified in `pkg-import`
//
// 3. Usage: go run genpkg.go -out-pkg-dir="tst" -out-pkg-name="std"
// This will generate Anko package using system go tools.
//
// 4. Usage: go run genpkg.go -pkg-src-dir="." -pkg-import="github.com/ipsusila/pola/es/ank" -out-pkg-dir="tst" -out-pkg-name="usr"
package main

import (
	"flag"
	"log/slog"

	"github.com/ipsusila/pola/es/ank"
)

func main() {
	var (
		stdVer     = flag.String("std-ver", "", "version of the Go SDK")
		pkgRootDir = flag.String("pkg-root-dir", "", "Root directory where go `sdk` or GOPATH is stored. Usually HOME directory")

		pkgVer    = flag.String("pkg-ver", "", "custom package version to be imported")
		pkgImport = flag.String("pkg-import", "", "Import directive of given package")
		pkgSrcDir = flag.String("pkg-src-dir", "", "Import package from given source directory")

		outPkgDir  = flag.String("out-pkg-dir", "pkg", "Root directory for storing imported packages")
		outPkgName = flag.String("out-pkg-name", "imp", "Package name for the imported packages")
	)
	flag.Parse()

	// default parameters
	tgtOutPkgDir := *outPkgDir
	if tgtOutPkgDir == "" {
		tgtOutPkgDir = "pkg"
	}
	tgtOutPkg := *outPkgName
	if tgtOutPkg == "" {
		tgtOutPkg = "imp"
	}
	rootDirs := []string{}
	if *pkgRootDir != "" {
		rootDirs = append(rootDirs, *pkgRootDir)
	}

	if *stdVer != "" {
		err := ank.GenerateStdPackages(*stdVer, tgtOutPkgDir, tgtOutPkg, rootDirs...)
		if err != nil {
			slog.Error("Error importing/generating standard package",
				"goVer", *stdVer,
				"rootPkg", tgtOutPkgDir,
				"targetPkg", tgtOutPkg,
				"error", err.Error(),
			)
		} else {
			slog.Info("Finished importing/generating standard package",
				"goVer", *stdVer,
				"rootPkg", tgtOutPkgDir,
				"targetPkg", tgtOutPkg,
			)
		}
	} else if *pkgImport != "" {
		// do custom import
		rootDirs := []string{}
		if *pkgRootDir != "" {
			rootDirs = append(rootDirs, *pkgRootDir)
		}
		err := ank.GenerateCustomPackages(*pkgVer, *pkgImport, *pkgSrcDir, tgtOutPkgDir, tgtOutPkg, rootDirs...)
		if err != nil {
			slog.Error("Error importing/generating user package",
				"pkgVer", *pkgVer,
				"pkgName", *pkgImport,
				"rootPkg", tgtOutPkgDir,
				"targetPkg", tgtOutPkg,
				"error", err.Error(),
			)
		} else {
			slog.Info("Finished importing/generating user package",
				"pkgVer", *pkgVer,
				"pkgName", *pkgImport,
				"rootPkg", tgtOutPkgDir,
				"targetPkg", tgtOutPkg,
			)
		}
	} else {
		// sistem go
		err := ank.GenerateStdPackages("", tgtOutPkgDir, tgtOutPkg, rootDirs...)
		if err != nil {
			slog.Error("Error importing/generating standard system package",
				"goVer", *stdVer,
				"rootPkg", tgtOutPkgDir,
				"targetPkg", tgtOutPkg,
				"error", err.Error(),
			)
		} else {
			slog.Info("Finished importing/generating standard system package",
				"goVer", *stdVer,
				"rootPkg", tgtOutPkgDir,
				"targetPkg", tgtOutPkg,
			)
		}
	}
}
