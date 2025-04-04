// Auto-generated code.
package pkg

import (
	"reflect"

	"github.com/jmoiron/sqlx"
)

var valJmoironSqlx = map[string]reflect.Value{
	//Function(s)
	"BindDriver":        reflect.ValueOf(sqlx.BindDriver),
	"BindNamed":         reflect.ValueOf(sqlx.BindNamed),
	"BindType":          reflect.ValueOf(sqlx.BindType),
	"Connect":           reflect.ValueOf(sqlx.Connect),
	"ConnectContext":    reflect.ValueOf(sqlx.ConnectContext),
	"Get":               reflect.ValueOf(sqlx.Get),
	"GetContext":        reflect.ValueOf(sqlx.GetContext),
	"In":                reflect.ValueOf(sqlx.In),
	"LoadFile":          reflect.ValueOf(sqlx.LoadFile),
	"LoadFileContext":   reflect.ValueOf(sqlx.LoadFileContext),
	"MapScan":           reflect.ValueOf(sqlx.MapScan),
	"MustConnect":       reflect.ValueOf(sqlx.MustConnect),
	"MustExec":          reflect.ValueOf(sqlx.MustExec),
	"MustExecContext":   reflect.ValueOf(sqlx.MustExecContext),
	"MustOpen":          reflect.ValueOf(sqlx.MustOpen),
	"Named":             reflect.ValueOf(sqlx.Named),
	"NamedExec":         reflect.ValueOf(sqlx.NamedExec),
	"NamedExecContext":  reflect.ValueOf(sqlx.NamedExecContext),
	"NamedQuery":        reflect.ValueOf(sqlx.NamedQuery),
	"NamedQueryContext": reflect.ValueOf(sqlx.NamedQueryContext),
	"NewDb":             reflect.ValueOf(sqlx.NewDb),
	"Open":              reflect.ValueOf(sqlx.Open),
	"Preparex":          reflect.ValueOf(sqlx.Preparex),
	"PreparexContext":   reflect.ValueOf(sqlx.PreparexContext),
	"Rebind":            reflect.ValueOf(sqlx.Rebind),
	"Select":            reflect.ValueOf(sqlx.Select),
	"SelectContext":     reflect.ValueOf(sqlx.SelectContext),
	"SliceScan":         reflect.ValueOf(sqlx.SliceScan),
	"StructScan":        reflect.ValueOf(sqlx.StructScan),
	//Variable(s)
	"NameMapper": reflect.ValueOf(sqlx.NameMapper),
	//Constant(s)
	"AT":       reflect.ValueOf(sqlx.AT),
	"DOLLAR":   reflect.ValueOf(sqlx.DOLLAR),
	"NAMED":    reflect.ValueOf(sqlx.NAMED),
	"QUESTION": reflect.ValueOf(sqlx.QUESTION),
	"UNKNOWN":  reflect.ValueOf(sqlx.UNKNOWN),
}

var typJmoironSqlx = map[string]reflect.Type{
	//Struct(s)
	"Conn":      reflect.TypeOf(sqlx.Conn{}),
	"DB":        reflect.TypeOf(sqlx.DB{}),
	"NamedStmt": reflect.TypeOf(sqlx.NamedStmt{}),
	"Row":       reflect.TypeOf(sqlx.Row{}),
	"Rows":      reflect.TypeOf(sqlx.Rows{}),
	"Stmt":      reflect.TypeOf(sqlx.Stmt{}),
	"Tx":        reflect.TypeOf(sqlx.Tx{}),
}
