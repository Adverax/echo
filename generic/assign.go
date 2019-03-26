package generic

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// Data scanner (see standard package database/sql)
type Scanner interface {
	Scan(src interface{}) error
}

type RawBytes []byte

var errNilPtr = errors.New("destination pointer is nil") // embedded in descriptive error

// ConvertAssign is variation of standard routine from package database/sql
// convertAssign copies to dest the value in src, converting it if possible.
// An error is returned if the copy would result in loss of information.
// dest should be a pointer type.

func ConvertAssign(dest, src interface{}) error {
	// Common cases, without reflect.
	switch s := src.(type) {
	case string:
		switch d := dest.(type) {
		case *string:
			if d == nil {
				return errNilPtr
			}
			*d = s
			return nil
		case *[]byte:
			if d == nil {
				return errNilPtr
			}
			*d = []byte(s)
			return nil
		}
	case []byte:
		switch d := dest.(type) {
		case *string:
			if d == nil {
				return errNilPtr
			}
			*d = string(s)
			return nil
		case *interface{}:
			if d == nil {
				return errNilPtr
			}
			*d = cloneBytes(s)
			return nil
		case *[]byte:
			if d == nil {
				return errNilPtr
			}
			*d = cloneBytes(s)
			return nil
		case *RawBytes:
			if d == nil {
				return errNilPtr
			}
			*d = s
			return nil
		}
	case time.Time:
		switch d := dest.(type) {
		case *string:
			*d = s.Format(time.RFC3339Nano)
			return nil
		case *[]byte:
			if d == nil {
				return errNilPtr
			}
			*d = []byte(s.Format(time.RFC3339Nano))
			return nil
		}
	case nil:
		switch d := dest.(type) {
		case *interface{}:
			if d == nil {
				return errNilPtr
			}
			*d = nil
			return nil
		case *[]byte:
			if d == nil {
				return errNilPtr
			}
			*d = nil
			return nil
		case *RawBytes:
			if d == nil {
				return errNilPtr
			}
			*d = nil
			return nil
		}
	case int8:
		switch d := dest.(type) {
		case *int8:
			*d = s
			return nil
		case *int16:
			*d = int16(s)
			return nil
		case *int32:
			*d = int32(s)
			return nil
		case *int64:
			*d = int64(s)
			return nil
		case *float32:
			*d = float32(s)
			return nil
		case *float64:
			*d = float64(s)
			return nil
		}
	case int16:
		switch d := dest.(type) {
		case *int8:
			*d = int8(s)
			return nil
		case *int16:
			*d = s
			return nil
		case *int32:
			*d = int32(s)
			return nil
		case *int64:
			*d = int64(s)
			return nil
		case *float32:
			*d = float32(s)
			return nil
		case *float64:
			*d = float64(s)
			return nil
		}
	case int32:
		switch d := dest.(type) {
		case *int8:
			*d = int8(s)
			return nil
		case *int16:
			*d = int16(s)
			return nil
		case *int32:
			*d = s
			return nil
		case *int64:
			*d = int64(s)
			return nil
		case *float32:
			*d = float32(s)
			return nil
		case *float64:
			*d = float64(s)
			return nil
		}
	case int64:
		switch d := dest.(type) {
		case *int8:
			*d = int8(s)
			return nil
		case *int16:
			*d = int16(s)
			return nil
		case *int32:
			*d = int32(s)
			return nil
		case *int64:
			*d = s
			return nil
		case *float32:
			*d = float32(s)
			return nil
		case *float64:
			*d = float64(s)
			return nil
		}
	case float32:
		switch d := dest.(type) {
		case *int8:
			v := int8(s)
			if float32(v) == s {
				*d = v
				return nil
			}
		case *int16:
			v := int16(s)
			if float32(v) == s {
				*d = v
				return nil
			}
		case *int32:
			v := int32(s)
			if float32(v) == s {
				*d = v
				return nil
			}
			return nil
		case *int64:
			v := int64(s)
			if float32(v) == s {
				*d = v
				return nil
			}
			return nil
		case *float32:
			*d = s
			return nil
		case *float64:
			*d = float64(s)
			return nil
		}
	case float64:
		switch d := dest.(type) {
		case *int8:
			v := int8(s)
			if float64(v) == s {
				*d = v
				return nil
			}
		case *int16:
			v := int16(s)
			if float64(v) == s {
				*d = v
				return nil
			}
		case *int32:
			v := int32(s)
			if float64(v) == s {
				*d = v
				return nil
			}
			return nil
		case *int64:
			v := int64(s)
			if float64(v) == s {
				*d = v
				return nil
			}
			return nil
		case *float32:
			*d = float32(s)
			return nil
		case *float64:
			*d = s
			return nil
		}
	}

	var sv reflect.Value

	switch d := dest.(type) {
	case *string:
		sv = reflect.ValueOf(src)
		switch sv.Kind() {
		case reflect.Bool,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			*d = asString(src)
			return nil
		}
	case *[]byte:
		sv = reflect.ValueOf(src)
		if b, ok := asBytes(nil, sv); ok {
			*d = b
			return nil
		}
	case *RawBytes:
		sv = reflect.ValueOf(src)
		if b, ok := asBytes([]byte(*d)[:0], sv); ok {
			*d = RawBytes(b)
			return nil
		}
	case *bool:
		bv, err := asBool(src)
		if err == nil {
			*d = bv
		}
		return err
	case *interface{}:
		*d = src
		return nil
	}

	if scanner, ok := dest.(Scanner); ok {
		return scanner.Scan(src)
	}

	dpv := reflect.ValueOf(dest)
	if dpv.Kind() != reflect.Ptr {
		return fmt.Errorf("destination not a pointer")
	}
	if dpv.IsNil() {
		return errNilPtr
	}

	if !sv.IsValid() {
		sv = reflect.ValueOf(src)
	}

	dv := reflect.Indirect(dpv)
	if sv.IsValid() && sv.Type().AssignableTo(dv.Type()) {
		switch b := src.(type) {
		case []byte:
			dv.Set(reflect.ValueOf(cloneBytes(b)))
		default:
			dv.Set(sv)
		}
		return nil
	}

	if dv.Kind() == sv.Kind() && sv.Type().ConvertibleTo(dv.Type()) {
		dv.Set(sv.Convert(dv.Type()))
		return nil
	}

	switch dv.Kind() {
	case reflect.Ptr:
		if src == nil {
			dv.Set(reflect.Zero(dv.Type()))
			return nil
		} else {
			dv.Set(reflect.New(dv.Type().Elem()))
			return ConvertAssign(dv.Interface(), src)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s := asString(src)
		i64, err := strconv.ParseInt(s, 10, dv.Type().Bits())
		if err != nil {
			err = strconvErr(err)
			return fmt.Errorf("converting driver.Value type %T (%q) to a %s: %v", src, s, dv.Kind(), err)
		}
		dv.SetInt(i64)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		s := asString(src)
		u64, err := strconv.ParseUint(s, 10, dv.Type().Bits())
		if err != nil {
			err = strconvErr(err)
			return fmt.Errorf("converting driver.Value type %T (%q) to a %s: %v", src, s, dv.Kind(), err)
		}
		dv.SetUint(u64)
		return nil
	case reflect.Float32, reflect.Float64:
		s := asString(src)
		f64, err := strconv.ParseFloat(s, dv.Type().Bits())
		if err != nil {
			err = strconvErr(err)
			return fmt.Errorf("converting driver.Value type %T (%q) to a %s: %v", src, s, dv.Kind(), err)
		}
		dv.SetFloat(f64)
		return nil
	}

	return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type %T", src, dest)
}

func strconvErr(err error) error {
	if ne, ok := err.(*strconv.NumError); ok {
		return ne.Err
	}
	return err
}

func cloneBytes(b []byte) []byte {
	if b == nil {
		return nil
	} else {
		c := make([]byte, len(b))
		copy(c, b)
		return c
	}
}

func asString(src interface{}) string {
	switch v := src.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	}
	rv := reflect.ValueOf(src)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(rv.Uint(), 10)
	case reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 64)
	case reflect.Float32:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 32)
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool())
	}
	return fmt.Sprintf("%v", src)
}

func asBytes(buf []byte, rv reflect.Value) (b []byte, ok bool) {
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.AppendInt(buf, rv.Int(), 10), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.AppendUint(buf, rv.Uint(), 10), true
	case reflect.Float32:
		return strconv.AppendFloat(buf, rv.Float(), 'g', -1, 32), true
	case reflect.Float64:
		return strconv.AppendFloat(buf, rv.Float(), 'g', -1, 64), true
	case reflect.Bool:
		return strconv.AppendBool(buf, rv.Bool()), true
	case reflect.String:
		s := rv.String()
		return append(buf, s...), true
	}
	return
}

func asBool(src interface{}) (bool, error) {
	switch v := src.(type) {
	case bool:
		return v, nil
	case int8:
		return v != 0, nil
	case int16:
		return v != 0, nil
	case int32:
		return v != 0, nil
	case int64:
		return v != 0, nil
	case uint8:
		return v != 0, nil
	case uint16:
		return v != 0, nil
	case uint32:
		return v != 0, nil
	case uint64:
		return v != 0, nil
	case []byte:
		return len(v) != 0, nil
	case string:
		vv, _ := strconv.ParseBool(v)
		return vv, nil
	default:
		return false, nil
	}
}

func ResetFields(rec interface{}, fields []string) {
	record := reflect.ValueOf(rec)
	for _, name := range fields {
		field := record.FieldByName(name)
		switch field.Type().Kind() {
		case reflect.Int8,
			reflect.Int16,
			reflect.Int32,
			reflect.Int64,
			reflect.Uint8,
			reflect.Uint16,
			reflect.Uint32,
			reflect.Uint64:
			field.SetInt(0)
		case reflect.Float32,
			reflect.Float64:
			field.SetFloat(0)
		case reflect.String:
			field.SetString("")
		}
	}
}
