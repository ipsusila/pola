package pkg

import (
	"context"
	"reflect"
)

var valContext = map[string]reflect.Value{
	"Canceled":          reflect.ValueOf(context.Canceled),
	"DeadlineExceeded":  reflect.ValueOf(context.DeadlineExceeded),
	"AfterFunc":         reflect.ValueOf(context.AfterFunc),
	"Cause":             reflect.ValueOf(context.Cause),
	"WithCancel":        reflect.ValueOf(context.WithCancel),
	"WithCancelCause":   reflect.ValueOf(context.WithCancelCause),
	"WithDeadline":      reflect.ValueOf(context.WithDeadline),
	"WithDeadlineCause": reflect.ValueOf(context.WithDeadlineCause),
	"WithTimeout":       reflect.ValueOf(context.WithTimeout),
	"WithTimeoutCause":  reflect.ValueOf(context.WithTimeoutCause),
	"Background":        reflect.ValueOf(context.Background),
	"TODO":              reflect.ValueOf(context.TODO),
	"WithValue":         reflect.ValueOf(context.WithValue),
	"WithoutCancel":     reflect.ValueOf(context.WithoutCancel),
}
