package pkg

import (
	"reflect"

	"github.com/mattn/anko/vm"
)

// list of packages
var (
	Pkgs = map[string]map[string]reflect.Value{
		"database/sql": valDatabaseSql,
		//"context":      valContext,
		//"jmoiron/sqlx": valJmoironSqlx,
	}

	PkgTypes = map[string]map[string]reflect.Type{
		"database/sql": typDatabaseSql,
		//"jmoiron/sqlx": typJmoironSqlx,
	}
)

// NewImporter for given packages
func NewImporter(pkgs ...string) vm.Importer {
	imp := vm.NewPackagesImporter(nil, nil)
	imp.AppendMap(Pkgs, PkgTypes, pkgs...)

	return imp
}
