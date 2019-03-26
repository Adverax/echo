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

type ValidatorText interface {
	Validate(ctx Context, value string) error
}

type ValidatorTextFunc func(ctx Context, value string) error

func (f ValidatorTextFunc) Validate(ctx Context, value string) error {
	return f(ctx, value)
}

type Text struct {
	Validator ValidatorText
}

func (codec *Text) Empty(ctx Context) (interface{}, error) {
	return "", nil
}

func (codec *Text) Encode(ctx Context, value string) (interface{}, error) {
	if codec.Validator != nil {
		err := codec.Validator.Validate(ctx, value)
		if err != nil {
			return nil, err
		}
	}

	return value, nil
}

func (codec *Text) Decode(ctx Context, value interface{}) (string, error) {
	val, _ := generic.ConvertToString(value)
	return val, nil
}

func (codec *Text) Format(ctx Context, value interface{}) (val interface{}, err error) {
	return formatDefault(ctx, codec, value)
}

type ValidatorSigned interface {
	Validate(ctx Context, value int64) error
}

type ValidatorSignedFunc func(ctx Context, value int64) error

func (f ValidatorSignedFunc) Validate(ctx Context, value int64) error {
	return f(ctx, value)
}

type Signed struct {
	Min       int64
	Max       int64
	Validator ValidatorSigned
}

func (codec *Signed) Empty(ctx Context) (interface{}, error) {
	return int64(0), nil
}

func (codec *Signed) Encode(ctx Context, value string) (interface{}, error) {
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

func (codec *Signed) Decode(ctx Context, value interface{}) (string, error) {
	val, _ := generic.ConvertToInt64(value)
	return strconv.FormatInt(int64(val), 10), nil
}

func (codec *Signed) Format(ctx Context, value interface{}) (val interface{}, err error) {
	return formatDefault(ctx, codec, value)
}

type ValidatorUnsigned interface {
	Validate(ctx Context, value uint64) error
}

type ValidatorUnsignedFunc func(ctx Context, value uint64) error

func (f ValidatorUnsignedFunc) Validate(ctx Context, value uint64) error {
	return f(ctx, value)
}

type Unsigned struct {
	Min       uint64
	Max       uint64
	Validator ValidatorUnsigned
}

func (codec *Unsigned) Empty(ctx Context) (interface{}, error) {
	return uint64(0), nil
}

func (codec *Unsigned) Encode(ctx Context, value string) (interface{}, error) {
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

func (codec *Unsigned) Decode(ctx Context, value interface{}) (string, error) {
	val, _ := generic.ConvertToUint64(value)
	return strconv.FormatUint(uint64(val), 10), nil
}

func (codec *Unsigned) Format(ctx Context, value interface{}) (val interface{}, err error) {
	return formatDefault(ctx, codec, value)
}

type ValidatorDecimal interface {
	Validate(ctx Context, value float64) error
}

type ValidatorDecimalFunc func(ctx Context, value float64) error

func (f ValidatorDecimalFunc) Validate(ctx Context, value float64) error {
	return f(ctx, value)
}

type Decimal struct {
	Min       float64
	Max       float64
	Validator ValidatorDecimal
}

func (codec *Decimal) Empty(ctx Context) (interface{}, error) {
	return float64(0), nil
}

func (codec *Decimal) Encode(ctx Context, value string) (interface{}, error) {
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

func (codec *Decimal) Decode(ctx Context, value interface{}) (string, error) {
	val, _ := generic.ConvertToFloat64(value)
	return fmt.Sprintf("%g", val), nil
}

func (codec *Decimal) Format(ctx Context, value interface{}) (val interface{}, err error) {
	return formatDefault(ctx, codec, value)
}

// Optional is codec for any optional value.
// Optional is wrapper for inner codec.
// Example:
//   codec := &Optional{
//     Codec: &Signed{},
//   }
type Optional struct {
	Codec
}

func (codec *Optional) Encode(ctx Context, value string) (interface{}, error) {
	if value == "" {
		return codec.Codec.Empty(ctx)
	}

	return codec.Codec.Encode(ctx, value)
}

func (codec *Optional) Decode(ctx Context, value interface{}) (string, error) {
	if generic.IsEmpty(value) {
		return "", nil
	}

	return codec.Codec.Decode(ctx, value)
}

func (codec *Optional) Format(ctx Context, value interface{}) (val interface{}, err error) {
	return formatDefault(ctx, codec, value)
}

// Abstract value formatter
type Formatter interface {
	Format(ctx Context, value interface{}) (val interface{}, err error)
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
		c = TextCodec
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
	TextCodec     = new(Text)
	SignedCodec   = new(Signed)
	UnsignedCodec = new(Unsigned)
	DecimalCodec  = new(Decimal)
	BoolCodec     = new(Unsigned) // Override in real application
)

func formatDefault(
	ctx Context,
	decoder Decoder,
	value interface{},
) (interface{}, error) {
	if decoder == nil {
		decoder = TextCodec
	}

	val, err := decoder.Decode(ctx, value)
	if err != nil {
		return "", err
	}

	return val, nil
}
