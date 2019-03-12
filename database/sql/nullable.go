package sql

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/adverax/echo/generic"
)

// NullString represents a string that may be null.
// NullString implements the Scanner interface so
// it can be used as a scan destination:
//
//  var s NullString
//  err := db.QueryRow("SELECT name FROM foo WHERE id=?", id).Scan(&s)
//  ...
//  if s.Valid {
//     // use s.String
//  } else {
//     // NULL value
//  }
//
type NullString struct {
	String string
	Valid  bool // Valid is true if String is not NULL
}

func (n *NullString) Scan(value interface{}) error {
	if value == nil {
		n.String, n.Valid = "", false
		return nil
	}
	n.Valid = true
	return generic.ConvertAssign(&n.String, value)
}

func (n NullString) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.String, nil
}

func (n NullString) Internal() driver.Value {
	if !n.Valid {
		return nil
	}
	return n.String
}

func (n NullString) External() interface{} {
	if !n.Valid {
		return ""
	}
	return n.String
}

func (n *NullString) MarshalJSON() ([]byte, error) {
	var res []byte
	var err error
	if n.Valid {
		res, err = json.Marshal(n.String)
	} else {
		res, err = json.Marshal(nil)
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (n *NullString) UnmarshalJSON(data []byte) error {
	// Unmarshalling into a pointer will let us detect null
	var x *string
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		n.Valid = true
		n.String = *x
	} else {
		n.Valid = false
	}
	return nil
}

// NullInt represents an int that may be null.
// NullInt implements the Scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullInt struct {
	Int   int
	Valid bool // Valid is true if Int is not NULL
}

func (n *NullInt) Scan(value interface{}) error {
	if value == nil {
		n.Int, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return generic.ConvertAssign(&n.Int, value)
}

func (n NullInt) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Int, nil
}

func (n NullInt) Internal() driver.Value {
	if !n.Valid {
		return nil
	}
	return n.Int
}

func (n NullInt) External() interface{} {
	if !n.Valid {
		return int(0)
	}
	return n.Int
}

func (n *NullInt) MarshalJSON() ([]byte, error) {
	var res []byte
	var err error
	if n.Valid {
		res, err = json.Marshal(n.Int)
	} else {
		res, err = json.Marshal(nil)
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (n *NullInt) UnmarshalJSON(data []byte) error {
	// Unmarshalling into a pointer will let us detect null
	var x *int
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		n.Valid = true
		n.Int = *x
	} else {
		n.Valid = false
	}
	return nil
}

// NullInt8 represents an int8 that may be null.
// NullInt8 implements the Scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullInt8 struct {
	Int8  int8
	Valid bool // Valid is true if Int8 is not NULL
}

func (n *NullInt8) Scan(value interface{}) error {
	if value == nil {
		n.Int8, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return generic.ConvertAssign(&n.Int8, value)
}

func (n NullInt8) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Int8, nil
}

func (n NullInt8) Internal() driver.Value {
	if !n.Valid {
		return nil
	}
	return n.Int8
}

func (n NullInt8) External() interface{} {
	if !n.Valid {
		return int8(0)
	}
	return n.Int8
}

func (n *NullInt8) MarshalJSON() ([]byte, error) {
	var res []byte
	var err error
	if n.Valid {
		res, err = json.Marshal(n.Int8)
	} else {
		res, err = json.Marshal(nil)
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (n *NullInt8) UnmarshalJSON(data []byte) error {
	// Unmarshalling into a pointer will let us detect null
	var x *int8
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		n.Valid = true
		n.Int8 = *x
	} else {
		n.Valid = false
	}
	return nil
}

// NullInt16 represents an int16 that may be null.
// NullInt16 implements the Scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullInt16 struct {
	Int16 int16
	Valid bool // Valid is true if Int16 is not NULL
}

func (n *NullInt16) Scan(value interface{}) error {
	if value == nil {
		n.Int16, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return generic.ConvertAssign(&n.Int16, value)
}

func (n NullInt16) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Int16, nil
}

func (n NullInt16) Internal() driver.Value {
	if !n.Valid {
		return nil
	}
	return n.Int16
}

func (n NullInt16) External() interface{} {
	if !n.Valid {
		return int16(0)
	}
	return n.Int16
}

func (n *NullInt16) MarshalJSON() ([]byte, error) {
	var res []byte
	var err error
	if n.Valid {
		res, err = json.Marshal(n.Int16)
	} else {
		res, err = json.Marshal(nil)
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (n *NullInt16) UnmarshalJSON(data []byte) error {
	// Unmarshalling into a pointer will let us detect null
	var x *int16
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		n.Valid = true
		n.Int16 = *x
	} else {
		n.Valid = false
	}
	return nil
}

// NullInt32 represents an int32 that may be null.
// NullInt32 implements the Scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullInt32 struct {
	Int32 int32
	Valid bool // Valid is true if Int32 is not NULL
}

func (n *NullInt32) Scan(value interface{}) error {
	if value == nil {
		n.Int32, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return generic.ConvertAssign(&n.Int32, value)
}

func (n NullInt32) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Int32, nil
}

func (n NullInt32) Internal() driver.Value {
	if !n.Valid {
		return nil
	}
	return n.Int32
}

func (n NullInt32) External() interface{} {
	if !n.Valid {
		return int32(0)
	}
	return n.Int32
}

func (n *NullInt32) MarshalJSON() ([]byte, error) {
	var res []byte
	var err error
	if n.Valid {
		res, err = json.Marshal(n.Int32)
	} else {
		res, err = json.Marshal(nil)
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (n *NullInt32) UnmarshalJSON(data []byte) error {
	// Unmarshalling into a pointer will let us detect null
	var x *int32
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		n.Valid = true
		n.Int32 = *x
	} else {
		n.Valid = false
	}
	return nil
}

// NullInt64 represents an int64 that may be null.
// NullInt64 implements the Scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullInt64 struct {
	Int64 int64
	Valid bool // Valid is true if Int64 is not NULL
}

func (n *NullInt64) Scan(value interface{}) error {
	if value == nil {
		n.Int64, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return generic.ConvertAssign(&n.Int64, value)
}

func (n NullInt64) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Int64, nil
}

func (n NullInt64) Internal() driver.Value {
	if !n.Valid {
		return nil
	}
	return n.Int64
}

func (n NullInt64) External() interface{} {
	if !n.Valid {
		return int64(0)
	}
	return n.Int64
}

func (n *NullInt64) MarshalJSON() ([]byte, error) {
	var res []byte
	var err error
	if n.Valid {
		res, err = json.Marshal(n.Int64)
	} else {
		res, err = json.Marshal(nil)
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (n *NullInt64) UnmarshalJSON(data []byte) error {
	// Unmarshalling into a pointer will let us detect null
	var x *int64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		n.Valid = true
		n.Int64 = *x
	} else {
		n.Valid = false
	}
	return nil
}

// NullUInt represents an uInt that may be null.
// NullUInt implements the Scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullUint struct {
	Uint  uint
	Valid bool // Valid is true if UInt is not NULL
}

func (n *NullUint) Scan(value interface{}) error {
	if value == nil {
		n.Uint, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return generic.ConvertAssign(&n.Uint, value)
}

func (n NullUint) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Uint, nil
}

func (n NullUint) Internal() driver.Value {
	if !n.Valid {
		return nil
	}
	return n.Uint
}

func (n NullUint) External() interface{} {
	if !n.Valid {
		return uint(0)
	}
	return n.Uint
}

func (n *NullUint) MarshalJSON() ([]byte, error) {
	var res []byte
	var err error
	if n.Valid {
		res, err = json.Marshal(n.Uint)
	} else {
		res, err = json.Marshal(nil)
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (n *NullUint) UnmarshalJSON(data []byte) error {
	// Unmarshalling into a pointer will let us detect null
	var x *uint
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		n.Valid = true
		n.Uint = *x
	} else {
		n.Valid = false
	}
	return nil
}

// NullUInt8 represents an uInt8 that may be null.
// NullUInt8 implements the Scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullUint8 struct {
	Uint8 uint8
	Valid bool // Valid is true if UInt8 is not NULL
}

func (n *NullUint8) Scan(value interface{}) error {
	if value == nil {
		n.Uint8, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return generic.ConvertAssign(&n.Uint8, value)
}

func (n NullUint8) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Uint8, nil
}

func (n NullUint8) Internal() driver.Value {
	if !n.Valid {
		return nil
	}
	return n.Uint8
}

func (n NullUint8) External() interface{} {
	if !n.Valid {
		return uint8(0)
	}
	return n.Uint8
}

func (n *NullUint8) MarshalJSON() ([]byte, error) {
	var res []byte
	var err error
	if n.Valid {
		res, err = json.Marshal(n.Uint8)
	} else {
		res, err = json.Marshal(nil)
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (n *NullUint8) UnmarshalJSON(data []byte) error {
	// Unmarshalling into a pointer will let us detect null
	var x *uint8
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		n.Valid = true
		n.Uint8 = *x
	} else {
		n.Valid = false
	}
	return nil
}

// NullUInt16 represents an uInt16 that may be null.
// NullUInt16 implements the Scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullUint16 struct {
	Uint16 uint16
	Valid  bool // Valid is true if UInt16 is not NULL
}

func (n *NullUint16) Scan(value interface{}) error {
	if value == nil {
		n.Uint16, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return generic.ConvertAssign(&n.Uint16, value)
}

func (n NullUint16) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Uint16, nil
}

func (n NullUint16) Internal() driver.Value {
	if !n.Valid {
		return nil
	}
	return n.Uint16
}

func (n NullUint16) External() interface{} {
	if !n.Valid {
		return uint16(0)
	}
	return n.Uint16
}

func (n *NullUint16) MarshalJSON() ([]byte, error) {
	var res []byte
	var err error
	if n.Valid {
		res, err = json.Marshal(n.Uint16)
	} else {
		res, err = json.Marshal(nil)
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (n *NullUint16) UnmarshalJSON(data []byte) error {
	// Unmarshalling into a pointer will let us detect null
	var x *uint16
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		n.Valid = true
		n.Uint16 = *x
	} else {
		n.Valid = false
	}
	return nil
}

// NullUInt32 represents an uInt32 that may be null.
// NullUInt32 implements the Scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullUint32 struct {
	Uint32 uint32
	Valid  bool // Valid is true if UInt32 is not NULL
}

func (n *NullUint32) Scan(value interface{}) error {
	if value == nil {
		n.Uint32, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return generic.ConvertAssign(&n.Uint32, value)
}

func (n NullUint32) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Uint32, nil
}

func (n NullUint32) Internal() driver.Value {
	if !n.Valid {
		return nil
	}
	return n.Uint32
}

func (n NullUint32) External() interface{} {
	if !n.Valid {
		return uint32(0)
	}
	return n.Uint32
}

func (n *NullUint32) MarshalJSON() ([]byte, error) {
	var res []byte
	var err error
	if n.Valid {
		res, err = json.Marshal(n.Uint32)
	} else {
		res, err = json.Marshal(nil)
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (n *NullUint32) UnmarshalJSON(data []byte) error {
	// Unmarshalling into a pointer will let us detect null
	var x *uint32
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		n.Valid = true
		n.Uint32 = *x
	} else {
		n.Valid = false
	}
	return nil
}

// NullUInt64 represents an uInt64 that may be null.
// NullUInt64 implements the Scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullUint64 struct {
	Uint64 uint64
	Valid  bool // Valid is true if UInt64 is not NULL
}

func (n *NullUint64) Scan(value interface{}) error {
	if value == nil {
		n.Uint64, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return generic.ConvertAssign(&n.Uint64, value)
}

func (n NullUint64) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Uint64, nil
}

func (n NullUint64) Internal() driver.Value {
	if !n.Valid {
		return nil
	}
	return n.Uint64
}

func (n NullUint64) External() interface{} {
	if !n.Valid {
		return uint64(0)
	}
	return n.Uint64
}

func (n *NullUint64) MarshalJSON() ([]byte, error) {
	var res []byte
	var err error
	if n.Valid {
		res, err = json.Marshal(n.Uint64)
	} else {
		res, err = json.Marshal(nil)
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (n *NullUint64) UnmarshalJSON(data []byte) error {
	// Unmarshalling into a pointer will let us detect null
	var x *uint64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		n.Valid = true
		n.Uint64 = *x
	} else {
		n.Valid = false
	}
	return nil
}

// NullFloat64 represents a float64 that may be null.
// NullFloat64 implements the Scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullFloat64 struct {
	Float64 float64
	Valid   bool // Valid is true if Float64 is not NULL
}

func (n *NullFloat64) Scan(value interface{}) error {
	if value == nil {
		n.Float64, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return generic.ConvertAssign(&n.Float64, value)
}

func (n NullFloat64) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Float64, nil
}

func (n NullFloat64) Internal() driver.Value {
	if !n.Valid {
		return nil
	}
	return n.Float64
}

func (n NullFloat64) External() interface{} {
	if !n.Valid {
		return float64(0)
	}
	return n.Float64
}

func (n *NullFloat64) MarshalJSON() ([]byte, error) {
	var res []byte
	var err error
	if n.Valid {
		res, err = json.Marshal(n.Float64)
	} else {
		res, err = json.Marshal(nil)
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (n *NullFloat64) UnmarshalJSON(data []byte) error {
	// Unmarshalling into a pointer will let us detect null
	var x *float64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		n.Valid = true
		n.Float64 = *x
	} else {
		n.Valid = false
	}
	return nil
}

// NullBool represents a bool that may be null.
// NullBool implements the Scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullBool struct {
	Bool  bool
	Valid bool // Valid is true if Bool is not NULL
}

func (n *NullBool) Scan(value interface{}) error {
	if value == nil {
		n.Bool, n.Valid = false, false
		return nil
	}
	n.Valid = true
	return generic.ConvertAssign(&n.Bool, value)
}

func (n NullBool) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Bool, nil
}

func (n NullBool) Internal() driver.Value {
	if !n.Valid {
		return nil
	}
	return n.Bool
}

func (n NullBool) External() interface{} {
	if !n.Valid {
		return false
	}
	return n.Bool
}

func (n *NullBool) MarshalJSON() ([]byte, error) {
	var res []byte
	var err error
	if n.Valid {
		res, err = json.Marshal(n.Bool)
	} else {
		res, err = json.Marshal(nil)
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (n *NullBool) UnmarshalJSON(data []byte) error {
	// Unmarshalling into a pointer will let us detect null
	var x *bool
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		n.Valid = true
		n.Bool = *x
	} else {
		n.Valid = false
	}
	return nil
}

// NullTime represents a time that may be null.
// NullTime implements the Scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

func (n *NullTime) Scan(value interface{}) error {
	var str sql.NullString
	err := str.Scan(value)
	if err != nil {
		n.Time, n.Valid = DefaultTime, false
		return err
	}
	tm, err := time.Parse("2006-01-02 15:04:05", str.String)
	if err != nil {
		n.Time, n.Valid = DefaultTime, false
		return err
	}
	n.Time = tm
	n.Valid = true
	return nil
}

func (n NullTime) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Time, nil
}

func (n NullTime) Internal() driver.Value {
	if !n.Valid {
		return nil
	}
	return n.Time
}

func (n NullTime) External() interface{} {
	if !n.Valid {
		return time.Time{}
	}
	return n.Time
}

func (n *NullTime) MarshalJSON() ([]byte, error) {
	var res []byte
	var err error
	if n.Valid {
		res, err = json.Marshal(n.Time)
	} else {
		res, err = json.Marshal(nil)
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (n *NullTime) UnmarshalJSON(data []byte) error {
	// Unmarshalling into a pointer will let us detect null
	var x *time.Time
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		n.Valid = true
		n.Time = *x
	} else {
		n.Valid = false
	}
	return nil
}

var (
	DefaultTime = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	null        interface{}
)

func ExternalInt(value int) interface{} {
	if value == 0 {
		return null
	}
	return value
}

func ExternalInt8(value int8) interface{} {
	if value == 0 {
		return null
	}
	return value
}

func ExternalInt16(value int16) interface{} {
	if value == 0 {
		return null
	}
	return value
}

func ExternalInt32(value int32) interface{} {
	if value == 0 {
		return null
	}
	return value
}

func ExternalInt64(value int64) interface{} {
	if value == 0 {
		return null
	}
	return value
}

func ExternalUint(value uint) interface{} {
	if value == 0 {
		return null
	}
	return value
}

func ExternalUint8(value uint8) interface{} {
	if value == 0 {
		return null
	}
	return value
}

func ExternalUint16(value uint16) interface{} {
	if value == 0 {
		return null
	}
	return value
}

func ExternalUint32(value uint32) interface{} {
	if value == 0 {
		return null
	}
	return value
}

func ExternalUint64(value uint64) interface{} {
	if value == 0 {
		return null
	}
	return value
}

func ExternalFloat64(value float64) interface{} {
	if value == 0 {
		return null
	}
	return value
}

func ExternalString(value string) interface{} {
	if value == "" {
		return null
	}
	return value
}

func ExternalBool(value bool) interface{} {
	if value == false {
		return null
	}
	return value
}

func ExternalTime(value time.Time) interface{} {
	if value.Unix() <= DefaultTime.Unix() {
		return null
	}
	return value
}
