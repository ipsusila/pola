package pkg

import (
	"database/sql"
	"reflect"
)

var valDatabaseSql = map[string]reflect.Value{
	"Drivers":              reflect.ValueOf(sql.Drivers),
	"Register":             reflect.ValueOf(sql.Register),
	"Open":                 reflect.ValueOf(sql.Open),
	"OpenDB":               reflect.ValueOf(sql.OpenDB),
	"Named":                reflect.ValueOf(sql.Named),
	"ErrConnDone":          reflect.ValueOf(sql.ErrConnDone),
	"ErrNoRows":            reflect.ValueOf(sql.ErrNoRows),
	"ErrTxDone":            reflect.ValueOf(sql.ErrTxDone),
	"LevelDefault":         reflect.ValueOf(sql.LevelDefault),
	"LevelReadUncommitted": reflect.ValueOf(sql.LevelReadUncommitted),
	"LevelReadCommitted":   reflect.ValueOf(sql.LevelReadCommitted),
	"LevelWriteCommitted":  reflect.ValueOf(sql.LevelWriteCommitted),
	"LevelRepeatableRead":  reflect.ValueOf(sql.LevelRepeatableRead),
	"LevelSnapshot":        reflect.ValueOf(sql.LevelSnapshot),
	"LevelSerializable":    reflect.ValueOf(sql.LevelSerializable),
	"LevelLinearizable":    reflect.ValueOf(sql.LevelLinearizable),
}

var typDatabaseSql = map[string]reflect.Type{
	"ColumnType":  reflect.TypeOf(sql.ColumnType{}),
	"Conn":        reflect.TypeOf(sql.Conn{}),
	"DB":          reflect.TypeOf(sql.DB{}),
	"DBStats":     reflect.TypeOf(sql.DBStats{}),
	"NamedArg":    reflect.TypeOf(sql.NamedArg{}),
	"NullBool":    reflect.TypeOf(sql.NullBool{}),
	"NullByte":    reflect.TypeOf(sql.NullByte{}),
	"NullFloat64": reflect.TypeOf(sql.NullFloat64{}),
	"NullInt16":   reflect.TypeOf(sql.NullInt16{}),
	"NullInt32":   reflect.TypeOf(sql.NullInt32{}),
	"NullInt64":   reflect.TypeOf(sql.NullInt64{}),
	"NullString":  reflect.TypeOf(sql.NullString{}),
	"NullTime":    reflect.TypeOf(sql.NullTime{}),
	"NullInt":     reflect.TypeOf(sql.Null[int]{}),
	"Out":         reflect.TypeOf(sql.Out{}),
	"RawBytes":    reflect.TypeOf(sql.RawBytes{}),
	"Row":         reflect.TypeOf(sql.Row{}),
	"Rows":        reflect.TypeOf(sql.Rows{}),
	"Stmt":        reflect.TypeOf(sql.Stmt{}),
	"Tx":          reflect.TypeOf(sql.Tx{}),
	"TxOptions":   reflect.TypeOf(sql.TxOptions{}),
}
