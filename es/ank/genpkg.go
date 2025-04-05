//go:build ignore

// Usage: go run genpkg.go -std-ver="1.23.8" -out-pkg-dir="tst" -out-pkg-name="std"
// This will generate Anko packages for go version 1.23.8 under "tst" dir, with package name of `std`
// Don't forget to install corresponding Go version. See https://go.dev/doc/manage-install
//
// Usage go run genpkg.go -pkg-ver="" -pkg-import-path="github.com/BurntSushi/toml" -out-pkg-dir="tst" -out-pkg-name="usr"
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

		pkgVer        = flag.String("pkg-ver", "", "custom package version to be imported")
		pkgImportPath = flag.String("pkg-import-path", "", "Import path of given package")

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
	} else if *pkgImportPath != "" {
		// do custom import
		rootDirs := []string{}
		if *pkgRootDir != "" {
			rootDirs = append(rootDirs, *pkgRootDir)
		}
		err := ank.GenerateCustomPackages(*pkgVer, *pkgImportPath, tgtOutPkgDir, tgtOutPkg, rootDirs...)
		if err != nil {
			slog.Error("Error importing/generating user package",
				"pkgVer", *pkgVer,
				"pkgName", *pkgImportPath,
				"rootPkg", tgtOutPkgDir,
				"targetPkg", tgtOutPkg,
				"error", err.Error(),
			)
		} else {
			slog.Info("Finished importing/generating user package",
				"pkgVer", *pkgVer,
				"pkgName", *pkgImportPath,
				"rootPkg", tgtOutPkgDir,
				"targetPkg", tgtOutPkg,
			)
		}
	}
}
