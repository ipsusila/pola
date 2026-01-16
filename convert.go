package pola

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	timeLayouts = []string{
		time.RFC3339,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339Nano,
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
		time.DateTime,
		time.DateOnly,
		time.TimeOnly,
		time.Layout,
		"2006/01/02 15:04:05",
		"2006/01/02",
	}
)

// ToBool convert any value to boolean
// Part of these code is taken from github.com/mattn/anko
func ToBool(v any) bool {
	rv := reflect.ValueOf(v)
	if !rv.IsValid() {
		return false
	}
	nt := reflect.TypeOf(true)
	if rv.Type().ConvertibleTo(nt) {
		return rv.Convert(nt).Bool()
	}
	if rv.Type().ConvertibleTo(reflect.TypeOf(1.0)) && rv.Convert(reflect.TypeOf(1.0)).Float() > 0.0 {
		return true
	}
	if rv.Kind() == reflect.String {
		s := strings.ToLower(rv.String())
		if s == "y" || s == "yes" {
			return true
		}
		b, err := strconv.ParseBool(s)
		if err == nil {
			return b
		}
	}
	return false
}

// ToInt convert any convertible value to int64
func ToInt(v any) (int64, bool) {
	rv := reflect.ValueOf(v)
	if !rv.IsValid() {
		return 0, false
	}
	nt := reflect.TypeOf(1)
	if rv.Type().ConvertibleTo(nt) {
		return rv.Convert(nt).Int(), true
	}
	if rv.Kind() == reflect.String {
		i, err := strconv.ParseInt(rv.String(), 10, 64)
		if err == nil {
			return i, true
		}
		f, err := strconv.ParseFloat(rv.String(), 64)
		if err == nil {
			return int64(f), true
		}
	}
	if rv.Kind() == reflect.Bool {
		if b, ok := v.(bool); ok {
			if b {
				return 1, true
			} else {
				return 0, true
			}
		}
	}
	return 0, false
}

// ToFloat convert any convertible vaue to float64.
func ToFloat(v any) (float64, bool) {
	rv := reflect.ValueOf(v)
	if !rv.IsValid() {
		return 0, false
	}
	nt := reflect.TypeOf(1.0)
	if rv.Type().ConvertibleTo(nt) {
		return rv.Convert(nt).Float(), true
	}
	if rv.Kind() == reflect.String {
		f, err := strconv.ParseFloat(rv.String(), 64)
		if err == nil {
			return f, true
		}
	}
	if rv.Kind() == reflect.Bool {
		if b, ok := v.(bool); ok {
			if b {
				return 1.0, true
			} else {
				return 0.0, true
			}
		}
	}
	return 0.0, false
}

// ToString convert any value to string representation
func ToString(v any) string {
	switch s := v.(type) {
	case string:
		return s
	case *string:
		return *s
	case []byte:
		return string(s)
	case fmt.Stringer:
		return s.String()
	}
	return fmt.Sprint(v)
}

// ToTime convert any to time
func ToTime(v any, opt ...*time.Location) (time.Time, bool) {
	if tm, ok := v.(time.Time); ok {
		return tm, true
	}

	// if string/[]byte/fmt.Stringer
	loc := time.Local
	if len(opt) > 0 && opt[0] != nil {
		loc = opt[0]
	}
	switch sv := v.(type) {
	case string:
		for _, layout := range timeLayouts {
			if tm, err := time.ParseInLocation(layout, sv, loc); err == nil {
				return tm, true
			}
		}
	case []byte:
		for _, layout := range timeLayouts {
			if tm, err := time.ParseInLocation(layout, string(sv), loc); err == nil {
				return tm, true
			}
		}
	case fmt.Stringer:
		for _, layout := range timeLayouts {
			if tm, err := time.ParseInLocation(layout, sv.String(), loc); err == nil {
				return tm, true
			}
		}
	default:
		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.String {
			for _, layout := range timeLayouts {
				if tm, err := time.ParseInLocation(layout, rv.String(), loc); err == nil {
					return tm, true
				}
			}
		}
	}

	return time.Time{}, false
}

// ToDuration convert any valid duration representation to time.Duration
func ToDuration(v any) (time.Duration, bool) {
	switch d := v.(type) {
	case time.Duration:
		return d, true
	case float32:
		return time.Duration(float64(d) * float64(time.Second)), true
	case float64:
		return time.Duration(d * float64(time.Second)), true
	case string:
		if dur, err := time.ParseDuration(d); err == nil {
			return dur, true
		}
	case []byte:
		if dur, err := time.ParseDuration(string(d)); err == nil {
			return dur, true
		}
	default:
		if sec, ok := ToInt(v); ok {
			return time.Duration(sec * int64(time.Second)), true
		}
	}
	return time.Duration(0), false
}
