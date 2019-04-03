// Copyright 2019 Adverax. All Rights Reserved.
// This file is part of project
//
//      http://github.com/adverax/echo
//
// Licensed under the MIT (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://github.com/adverax/echo/blob/master/LICENSE
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package echo

import (
	"fmt"
	"strconv"

	"github.com/adverax/echo/generic"
)

// Encoder used for encode value to internal representation
type Encoder interface {
	Encode(ctx Context, value string) (interface{}, error)
}

// Decoder used for decode value from internal representation
type Decoder interface {
	Decode(ctx Context, value interface{}) (string, error)
}

// Coder and Decoder
type Codec interface {
	Encoder
	Decoder
	// Get internal empty value
	Empty(ctx Context) (interface{}, error)
}

type ValidatorString interface {
	Validate(ctx Context, value string) error
}

type ValidatorStringFunc func(ctx Context, value string) error

func (f ValidatorStringFunc) Validate(ctx Context, value string) error {
	return f(ctx, value)
}

type String struct {
	Validator ValidatorString
}

func (codec *String) Empty(ctx Context) (interface{}, error) {
	return "", nil
}

func (codec *String) Encode(ctx Context, value string) (interface{}, error) {
	if codec.Validator != nil {
		err := codec.Validator.Validate(ctx, value)
		if err != nil {
			return nil, err
		}
	}

	return value, nil
}

func (codec *String) Decode(ctx Context, value interface{}) (string, error) {
	val, _ := generic.ConvertToString(value)
	return val, nil
}

func (codec *String) Format(ctx Context, value interface{}) (val interface{}, err error) {
	return formatDefault(ctx, codec, value)
}

type ValidatorInt interface {
	Validate(ctx Context, value int) error
}

type ValidatorIntFunc func(ctx Context, value int) error

func (f ValidatorIntFunc) Validate(ctx Context, value int) error {
	return f(ctx, value)
}

type Int struct {
	Min       int
	Max       int
	Validator ValidatorInt
}

func (codec *Int) Empty(ctx Context) (interface{}, error) {
	return int(0), nil
}

func (codec *Int) Encode(ctx Context, value string) (interface{}, error) {
	val, ok := generic.ConvertToInt(value)
	if !ok {
		return nil, ValidationErrors{
			ValidationErrorInvalidValue,
		}
	}

	var errs ValidationErrors
	if codec.Min < codec.Max {
		if val < codec.Min {
			errs = append(
				errs,
				NewValidationErrorMustBeNotBelow(
					fmt.Sprintf("%d", codec.Min),
				),
			)
		}

		if val > codec.Max {
			errs = append(
				errs,
				NewValidationErrorMustBeNotAbove(
					fmt.Sprintf("%d", codec.Max),
				),
			)
		}
	}

	validator := codec.Validator
	if validator != nil {
		var err error
		errs, err = AppendValidationError(
			errs,
			codec.Validator.Validate(ctx, val),
		)
		if err != nil {
			return nil, err
		}
	}

	if len(errs) == 0 {
		return val, nil
	} else {
		return 0, errs
	}
}

func (codec *Int) Decode(ctx Context, value interface{}) (string, error) {
	val, _ := generic.ConvertToInt(value)
	return strconv.FormatInt(int64(val), 10), nil
}

func (codec *Int) Format(ctx Context, value interface{}) (val interface{}, err error) {
	return formatDefault(ctx, codec, value)
}

type ValidatorInt8 interface {
	Validate(ctx Context, value int8) error
}

type ValidatorInt8Func func(ctx Context, value int8) error

func (f ValidatorInt8Func) Validate(ctx Context, value int8) error {
	return f(ctx, value)
}

type Int8 struct {
	Min       int8
	Max       int8
	Validator ValidatorInt8
}

func (codec *Int8) Empty(ctx Context) (interface{}, error) {
	return int8(0), nil
}

func (codec *Int8) Encode(ctx Context, value string) (interface{}, error) {
	val, ok := generic.ConvertToInt8(value)
	if !ok {
		return nil, ValidationErrors{
			ValidationErrorInvalidValue,
		}
	}

	var errs ValidationErrors
	if codec.Min < codec.Max {
		if val < codec.Min {
			errs = append(
				errs,
				NewValidationErrorMustBeNotBelow(
					fmt.Sprintf("%d", codec.Min),
				),
			)
		}

		if val > codec.Max {
			errs = append(
				errs,
				NewValidationErrorMustBeNotAbove(
					fmt.Sprintf("%d", codec.Max),
				),
			)
		}
	}

	validator := codec.Validator
	if validator != nil {
		var err error
		errs, err = AppendValidationError(
			errs,
			codec.Validator.Validate(ctx, val),
		)
		if err != nil {
			return nil, err
		}
	}

	if len(errs) == 0 {
		return val, nil
	} else {
		return 0, errs
	}
}

func (codec *Int8) Decode(ctx Context, value interface{}) (string, error) {
	val, _ := generic.ConvertToInt8(value)
	return strconv.FormatInt(int64(val), 10), nil
}

func (codec *Int8) Format(ctx Context, value interface{}) (val interface{}, err error) {
	return formatDefault(ctx, codec, value)
}

type ValidatorInt16 interface {
	Validate(ctx Context, value int16) error
}

type ValidatorInt16Func func(ctx Context, value int16) error

func (f ValidatorInt16Func) Validate(ctx Context, value int16) error {
	return f(ctx, value)
}

type Int16 struct {
	Min       int16
	Max       int16
	Validator ValidatorInt16
}

func (codec *Int16) Empty(ctx Context) (interface{}, error) {
	return int16(0), nil
}

func (codec *Int16) Encode(ctx Context, value string) (interface{}, error) {
	val, ok := generic.ConvertToInt16(value)
	if !ok {
		return nil, ValidationErrors{
			ValidationErrorInvalidValue,
		}
	}

	var errs ValidationErrors
	if codec.Min < codec.Max {
		if val < codec.Min {
			errs = append(
				errs,
				NewValidationErrorMustBeNotBelow(
					fmt.Sprintf("%d", codec.Min),
				),
			)
		}

		if val > codec.Max {
			errs = append(
				errs,
				NewValidationErrorMustBeNotAbove(
					fmt.Sprintf("%d", codec.Max),
				),
			)
		}
	}

	validator := codec.Validator
	if validator != nil {
		var err error
		errs, err = AppendValidationError(
			errs,
			codec.Validator.Validate(ctx, val),
		)
		if err != nil {
			return nil, err
		}
	}

	if len(errs) == 0 {
		return val, nil
	} else {
		return 0, errs
	}
}

func (codec *Int16) Decode(ctx Context, value interface{}) (string, error) {
	val, _ := generic.ConvertToInt16(value)
	return strconv.FormatInt(int64(val), 10), nil
}

func (codec *Int16) Format(ctx Context, value interface{}) (val interface{}, err error) {
	return formatDefault(ctx, codec, value)
}

type ValidatorInt32 interface {
	Validate(ctx Context, value int32) error
}

type ValidatorInt32Func func(ctx Context, value int32) error

func (f ValidatorInt32Func) Validate(ctx Context, value int32) error {
	return f(ctx, value)
}

type Int32 struct {
	Min       int32
	Max       int32
	Validator ValidatorInt32
}

func (codec *Int32) Empty(ctx Context) (interface{}, error) {
	return int32(0), nil
}

func (codec *Int32) Encode(ctx Context, value string) (interface{}, error) {
	val, ok := generic.ConvertToInt32(value)
	if !ok {
		return nil, ValidationErrors{
			ValidationErrorInvalidValue,
		}
	}

	var errs ValidationErrors
	if codec.Min < codec.Max {
		if val < codec.Min {
			errs = append(
				errs,
				NewValidationErrorMustBeNotBelow(
					fmt.Sprintf("%d", codec.Min),
				),
			)
		}

		if val > codec.Max {
			errs = append(
				errs,
				NewValidationErrorMustBeNotAbove(
					fmt.Sprintf("%d", codec.Max),
				),
			)
		}
	}

	validator := codec.Validator
	if validator != nil {
		var err error
		errs, err = AppendValidationError(
			errs,
			codec.Validator.Validate(ctx, val),
		)
		if err != nil {
			return nil, err
		}
	}

	if len(errs) == 0 {
		return val, nil
	} else {
		return 0, errs
	}
}

func (codec *Int32) Decode(ctx Context, value interface{}) (string, error) {
	val, _ := generic.ConvertToInt32(value)
	return strconv.FormatInt(int64(val), 10), nil
}

func (codec *Int32) Format(ctx Context, value interface{}) (val interface{}, err error) {
	return formatDefault(ctx, codec, value)
}

type ValidatorInt64 interface {
	Validate(ctx Context, value int64) error
}

type ValidatorInt64Func func(ctx Context, value int64) error

func (f ValidatorInt64Func) Validate(ctx Context, value int64) error {
	return f(ctx, value)
}

type Int64 struct {
	Min       int64
	Max       int64
	Validator ValidatorInt64
}

func (codec *Int64) Empty(ctx Context) (interface{}, error) {
	return int64(0), nil
}

func (codec *Int64) Encode(ctx Context, value string) (interface{}, error) {
	val, ok := generic.ConvertToInt64(value)
	if !ok {
		return nil, ValidationErrors{
			ValidationErrorInvalidValue,
		}
	}

	var errs ValidationErrors
	if codec.Min < codec.Max {
		if val < codec.Min {
			errs = append(
				errs,
				NewValidationErrorMustBeNotBelow(
					fmt.Sprintf("%d", codec.Min),
				),
			)
		}

		if val > codec.Max {
			errs = append(
				errs,
				NewValidationErrorMustBeNotAbove(
					fmt.Sprintf("%d", codec.Max),
				),
			)
		}
	}

	validator := codec.Validator
	if validator != nil {
		var err error
		errs, err = AppendValidationError(
			errs,
			codec.Validator.Validate(ctx, val),
		)
		if err != nil {
			return nil, err
		}
	}

	if len(errs) == 0 {
		return val, nil
	} else {
		return 0, errs
	}
}

func (codec *Int64) Decode(ctx Context, value interface{}) (string, error) {
	val, _ := generic.ConvertToInt64(value)
	return strconv.FormatInt(int64(val), 10), nil
}

func (codec *Int64) Format(ctx Context, value interface{}) (val interface{}, err error) {
	return formatDefault(ctx, codec, value)
}

type ValidatorUint8 interface {
	Validate(ctx Context, value uint8) error
}

type ValidatorUint8Func func(ctx Context, value uint8) error

func (f ValidatorUint8Func) Validate(ctx Context, value uint8) error {
	return f(ctx, value)
}

type ValidatorUint interface {
	Validate(ctx Context, value uint) error
}

type ValidatorUintFunc func(ctx Context, value uint) error

func (f ValidatorUintFunc) Validate(ctx Context, value uint) error {
	return f(ctx, value)
}

type Uint struct {
	Min       uint
	Max       uint
	Validator ValidatorUint
}

func (codec *Uint) Empty(ctx Context) (interface{}, error) {
	return uint(0), nil
}

func (codec *Uint) Encode(ctx Context, value string) (interface{}, error) {
	val, ok := generic.ConvertToUint(value)
	if !ok {
		return nil, ValidationErrors{
			ValidationErrorInvalidValue,
		}
	}

	var errs ValidationErrors
	if codec.Min < codec.Max {
		if val < codec.Min {
			errs = append(
				errs,
				NewValidationErrorMustBeNotBelow(
					fmt.Sprintf("%d", codec.Min),
				),
			)
		}

		if val > codec.Max {
			errs = append(
				errs,
				NewValidationErrorMustBeNotAbove(
					fmt.Sprintf("%d", codec.Max),
				),
			)
		}
	}

	validator := codec.Validator
	if validator != nil {
		var err error
		errs, err = AppendValidationError(
			errs,
			codec.Validator.Validate(ctx, val),
		)
		if err != nil {
			return nil, err
		}
	}

	if len(errs) == 0 {
		return val, nil
	} else {
		return 0, errs
	}
}

func (codec *Uint) Decode(ctx Context, value interface{}) (string, error) {
	val, _ := generic.ConvertToUint(value)
	return strconv.FormatUint(uint64(val), 10), nil
}

func (codec *Uint) Format(ctx Context, value interface{}) (val interface{}, err error) {
	return formatDefault(ctx, codec, value)
}

type Uint8 struct {
	Min       uint8
	Max       uint8
	Validator ValidatorUint8
}

func (codec *Uint8) Empty(ctx Context) (interface{}, error) {
	return uint8(0), nil
}

func (codec *Uint8) Encode(ctx Context, value string) (interface{}, error) {
	val, ok := generic.ConvertToUint8(value)
	if !ok {
		return nil, ValidationErrors{
			ValidationErrorInvalidValue,
		}
	}

	var errs ValidationErrors
	if codec.Min < codec.Max {
		if val < codec.Min {
			errs = append(
				errs,
				NewValidationErrorMustBeNotBelow(
					fmt.Sprintf("%d", codec.Min),
				),
			)
		}

		if val > codec.Max {
			errs = append(
				errs,
				NewValidationErrorMustBeNotAbove(
					fmt.Sprintf("%d", codec.Max),
				),
			)
		}
	}

	validator := codec.Validator
	if validator != nil {
		var err error
		errs, err = AppendValidationError(
			errs,
			codec.Validator.Validate(ctx, val),
		)
		if err != nil {
			return nil, err
		}
	}

	if len(errs) == 0 {
		return val, nil
	} else {
		return 0, errs
	}
}

func (codec *Uint8) Decode(ctx Context, value interface{}) (string, error) {
	val, _ := generic.ConvertToUint8(value)
	return strconv.FormatUint(uint64(val), 10), nil
}

func (codec *Uint8) Format(ctx Context, value interface{}) (val interface{}, err error) {
	return formatDefault(ctx, codec, value)
}

type ValidatorUint16 interface {
	Validate(ctx Context, value uint16) error
}

type ValidatorUint16Func func(ctx Context, value uint16) error

func (f ValidatorUint16Func) Validate(ctx Context, value uint16) error {
	return f(ctx, value)
}

type Uint16 struct {
	Min       uint16
	Max       uint16
	Validator ValidatorUint16
}

func (codec *Uint16) Empty(ctx Context) (interface{}, error) {
	return uint16(0), nil
}

func (codec *Uint16) Encode(ctx Context, value string) (interface{}, error) {
	val, ok := generic.ConvertToUint16(value)
	if !ok {
		return nil, ValidationErrors{
			ValidationErrorInvalidValue,
		}
	}

	var errs ValidationErrors
	if codec.Min < codec.Max {
		if val < codec.Min {
			errs = append(
				errs,
				NewValidationErrorMustBeNotBelow(
					fmt.Sprintf("%d", codec.Min),
				),
			)
		}

		if val > codec.Max {
			errs = append(
				errs,
				NewValidationErrorMustBeNotAbove(
					fmt.Sprintf("%d", codec.Max),
				),
			)
		}
	}

	validator := codec.Validator
	if validator != nil {
		var err error
		errs, err = AppendValidationError(
			errs,
			codec.Validator.Validate(ctx, val),
		)
		if err != nil {
			return nil, err
		}
	}

	if len(errs) == 0 {
		return val, nil
	} else {
		return 0, errs
	}
}

func (codec *Uint16) Decode(ctx Context, value interface{}) (string, error) {
	val, _ := generic.ConvertToUint16(value)
	return strconv.FormatUint(uint64(val), 10), nil
}

func (codec *Uint16) Format(ctx Context, value interface{}) (val interface{}, err error) {
	return formatDefault(ctx, codec, value)
}

type ValidatorUint32 interface {
	Validate(ctx Context, value uint32) error
}

type ValidatorUint32Func func(ctx Context, value uint32) error

func (f ValidatorUint32Func) Validate(ctx Context, value uint32) error {
	return f(ctx, value)
}

type Uint32 struct {
	Min       uint32
	Max       uint32
	Validator ValidatorUint32
}

func (codec *Uint32) Empty(ctx Context) (interface{}, error) {
	return uint32(0), nil
}

func (codec *Uint32) Encode(ctx Context, value string) (interface{}, error) {
	val, ok := generic.ConvertToUint32(value)
	if !ok {
		return nil, ValidationErrors{
			ValidationErrorInvalidValue,
		}
	}

	var errs ValidationErrors
	if codec.Min < codec.Max {
		if val < codec.Min {
			errs = append(
				errs,
				NewValidationErrorMustBeNotBelow(
					fmt.Sprintf("%d", codec.Min),
				),
			)
		}

		if val > codec.Max {
			errs = append(
				errs,
				NewValidationErrorMustBeNotAbove(
					fmt.Sprintf("%d", codec.Max),
				),
			)
		}
	}

	validator := codec.Validator
	if validator != nil {
		var err error
		errs, err = AppendValidationError(
			errs,
			codec.Validator.Validate(ctx, val),
		)
		if err != nil {
			return nil, err
		}
	}

	if len(errs) == 0 {
		return val, nil
	} else {
		return 0, errs
	}
}

func (codec *Uint32) Decode(ctx Context, value interface{}) (string, error) {
	val, _ := generic.ConvertToUint32(value)
	return strconv.FormatUint(uint64(val), 10), nil
}

func (codec *Uint32) Format(ctx Context, value interface{}) (val interface{}, err error) {
	return formatDefault(ctx, codec, value)
}

type ValidatorUint64 interface {
	Validate(ctx Context, value uint64) error
}

type ValidatorUint64Func func(ctx Context, value uint64) error

func (f ValidatorUint64Func) Validate(ctx Context, value uint64) error {
	return f(ctx, value)
}

type Uint64 struct {
	Min       uint64
	Max       uint64
	Validator ValidatorUint64
}

func (codec *Uint64) Empty(ctx Context) (interface{}, error) {
	return uint64(0), nil
}

func (codec *Uint64) Encode(ctx Context, value string) (interface{}, error) {
	val, ok := generic.ConvertToUint64(value)
	if !ok {
		return nil, ValidationErrors{
			ValidationErrorInvalidValue,
		}
	}

	var errs ValidationErrors
	if codec.Min < codec.Max {
		if val < codec.Min {
			errs = append(
				errs,
				NewValidationErrorMustBeNotBelow(
					fmt.Sprintf("%d", codec.Min),
				),
			)
		}

		if val > codec.Max {
			errs = append(
				errs,
				NewValidationErrorMustBeNotAbove(
					fmt.Sprintf("%d", codec.Max),
				),
			)
		}
	}

	validator := codec.Validator
	if validator != nil {
		var err error
		errs, err = AppendValidationError(
			errs,
			codec.Validator.Validate(ctx, val),
		)
		if err != nil {
			return nil, err
		}
	}

	if len(errs) == 0 {
		return val, nil
	} else {
		return 0, errs
	}
}

func (codec *Uint64) Decode(ctx Context, value interface{}) (string, error) {
	val, _ := generic.ConvertToUint64(value)
	return strconv.FormatUint(uint64(val), 10), nil
}

func (codec *Uint64) Format(ctx Context, value interface{}) (val interface{}, err error) {
	return formatDefault(ctx, codec, value)
}

type ValidatorFloat32 interface {
	Validate(ctx Context, value float32) error
}

type ValidatorFloat32Func func(ctx Context, value float32) error

func (f ValidatorFloat32Func) Validate(ctx Context, value float32) error {
	return f(ctx, value)
}

type Float32 struct {
	Min       float32
	Max       float32
	Validator ValidatorFloat32
}

func (codec *Float32) Empty(ctx Context) (interface{}, error) {
	return float32(0), nil
}

func (codec *Float32) Encode(ctx Context, value string) (interface{}, error) {
	val, ok := generic.ConvertToFloat32(value)
	if !ok {
		return nil, ValidationErrors{
			ValidationErrorInvalidValue,
		}
	}

	var errs ValidationErrors
	if codec.Min < codec.Max {
		if val < codec.Min {
			errs = append(
				errs,
				NewValidationErrorMustBeNotBelow(
					fmt.Sprintf("%g", codec.Min),
				),
			)
		}

		if val > codec.Max {
			errs = append(
				errs,
				NewValidationErrorMustBeNotAbove(
					fmt.Sprintf("%g", codec.Max),
				),
			)
		}
	}

	validator := codec.Validator
	if validator != nil {
		var err error
		errs, err = AppendValidationError(
			errs,
			codec.Validator.Validate(ctx, val),
		)
		if err != nil {
			return nil, err
		}
	}

	if len(errs) == 0 {
		return val, nil
	} else {
		return 0, errs
	}
}

func (codec *Float32) Decode(ctx Context, value interface{}) (string, error) {
	val, _ := generic.ConvertToFloat32(value)
	return fmt.Sprintf("%g", val), nil
}

func (codec *Float32) Format(ctx Context, value interface{}) (val interface{}, err error) {
	return formatDefault(ctx, codec, value)
}

type ValidatorFloat64 interface {
	Validate(ctx Context, value float64) error
}

type ValidatorFloat64Func func(ctx Context, value float64) error

func (f ValidatorFloat64Func) Validate(ctx Context, value float64) error {
	return f(ctx, value)
}

type Float64 struct {
	Min       float64
	Max       float64
	Validator ValidatorFloat64
}

func (codec *Float64) Empty(ctx Context) (interface{}, error) {
	return float64(0), nil
}

func (codec *Float64) Encode(ctx Context, value string) (interface{}, error) {
	val, ok := generic.ConvertToFloat64(value)
	if !ok {
		return nil, ValidationErrors{
			ValidationErrorInvalidValue,
		}
	}

	var errs ValidationErrors
	if codec.Min < codec.Max {
		if val < codec.Min {
			errs = append(
				errs,
				NewValidationErrorMustBeNotBelow(
					fmt.Sprintf("%g", codec.Min),
				),
			)
		}

		if val > codec.Max {
			errs = append(
				errs,
				NewValidationErrorMustBeNotAbove(
					fmt.Sprintf("%g", codec.Max),
				),
			)
		}
	}

	validator := codec.Validator
	if validator != nil {
		var err error
		errs, err = AppendValidationError(
			errs,
			codec.Validator.Validate(ctx, val),
		)
		if err != nil {
			return nil, err
		}
	}

	if len(errs) == 0 {
		return val, nil
	} else {
		return 0, errs
	}
}

func (codec *Float64) Decode(ctx Context, value interface{}) (string, error) {
	val, _ := generic.ConvertToFloat64(value)
	return fmt.Sprintf("%g", val), nil
}

func (codec *Float64) Format(ctx Context, value interface{}) (val interface{}, err error) {
	return formatDefault(ctx, codec, value)
}

// Abstract value formatter
type Formatter interface {
	Format(ctx Context, value interface{}) (val interface{}, err error)
}

// OptionalFormatter is formatter for any optional value.
// OptionalFormatter is wrapper for inner formatter.
// Example:
//   formatter := &OptionalFormatter{
//     Formatter: &Signed{},
//   }
type OptionalFormatter struct {
	Formatter
}

func (formatter *OptionalFormatter) Format(ctx Context, value interface{}) (val interface{}, err error) {
	if generic.IsEmpty(value) {
		return "", nil
	}

	return formatter.Formatter.Format(ctx, value)
}

// Formatter, that based on codec.Decode method
type BaseFormatter struct {
	Decoder
}

func (w *BaseFormatter) Format(
	ctx Context,
	value interface{},
) (interface{}, error) {
	c := w.Decoder
	if c == nil {
		c = StringCodec
	}

	return formatDefault(ctx, w.Decoder, value)
}

// Abstract value converter
type Converter interface {
	Codec
	Formatter
}

type converter struct {
	Codec
	Formatter
}

func NewConverter(codec Codec) Converter {
	return &converter{
		Codec:     codec,
		Formatter: &BaseFormatter{Decoder: codec},
	}
}

// Abstract pair converter
type pairs = DataSet

type PairConverter interface {
	pairs
	Formatter
}

type pairConverter struct {
	pairs
	Formatter
}

func NewPairConverter(codec DataSet) PairConverter {
	return &pairConverter{
		pairs: codec,
		Formatter: &BaseFormatter{
			Decoder: codec,
		},
	}
}

var (
	StringCodec              = new(String)
	IntCodec                 = new(Int)
	Int8Codec                = new(Int8)
	Int16Codec               = new(Int16)
	Int32Codec               = new(Int32)
	Int64Codec               = new(Int64)
	UintCodec                = new(Uint)
	Uint8Codec               = new(Uint8)
	Uint16Codec              = new(Uint16)
	Uint32Codec              = new(Uint32)
	Uint64Codec              = new(Uint64)
	Float32Codec             = new(Float32)
	Float64Codec             = new(Float64)
	BoolCodec                = new(Uint) // Override in real application
	OptionalIntFormatter     = &OptionalFormatter{Formatter: IntCodec}
	OptionalInt8Formatter    = &OptionalFormatter{Formatter: Int8Codec}
	OptionalInt16Formatter   = &OptionalFormatter{Formatter: Int16Codec}
	OptionalInt32Formatter   = &OptionalFormatter{Formatter: Int32Codec}
	OptionalInt64Formatter   = &OptionalFormatter{Formatter: Int64Codec}
	OptionalUintFormatter    = &OptionalFormatter{Formatter: UintCodec}
	OptionalUint8Formatter   = &OptionalFormatter{Formatter: Uint8Codec}
	OptionalUint16Formatter  = &OptionalFormatter{Formatter: Uint16Codec}
	OptionalUint32Formatter  = &OptionalFormatter{Formatter: Uint32Codec}
	OptionalUint64Formatter  = &OptionalFormatter{Formatter: Uint64Codec}
	OptionalFloat32Formatter = &OptionalFormatter{Formatter: Float32Codec}
	OptionalFloat64Formatter = &OptionalFormatter{Formatter: Float64Codec}
)

func formatDefault(
	ctx Context,
	decoder Decoder,
	value interface{},
) (interface{}, error) {
	if decoder == nil {
		decoder = StringCodec
	}

	val, err := decoder.Decode(ctx, value)
	if err != nil {
		return "", err
	}

	return val, nil
}
