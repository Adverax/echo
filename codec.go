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
}

// Advanced codec with items enumeration.
type PairCodec interface {
	Codec
	PairEnumerator
	DataSetProvider
}

// Pair enumerator
type PairEnumerator interface {
	Enumerate(ctx Context, action PairConsumer) error
}

type PairConsumer func(key string, value string) error

/*func (fn PairEnumeratorFunc) Enumerate(ctx Context, action PairEnumeratorFunc) error {
	return fn(ctx, action)
}*/

type ValidatorText interface {
	Validate(ctx Context, value string) error
}

type ValidatorFuncText func(ctx Context, value string) error

func (f ValidatorFuncText) Validate(ctx Context, value string) error {
	return f(ctx, value)
}

type Text struct {
	Validator ValidatorText
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

func (codec *Text) IsEmpty(value interface{}) bool {
	val, _ := generic.ConvertToString(value)
	return val == ""
}

type ValidatorSigned interface {
	Validate(ctx Context, value int64) error
}

type ValidatorFuncSigned func(ctx Context, value int64) error

func (f ValidatorFuncSigned) Validate(ctx Context, value int64) error {
	return f(ctx, value)
}

type Signed struct {
	Min       int64
	Max       int64
	Validator ValidatorSigned
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

func (codec *Signed) IsEmpty(value interface{}) bool {
	val, _ := generic.ConvertToInt64(value)
	return val == 0
}

type ValidatorUnsigned interface {
	Validate(ctx Context, value uint64) error
}

type ValidatorFuncUnsigned func(ctx Context, value uint64) error

func (f ValidatorFuncUnsigned) Validate(ctx Context, value uint64) error {
	return f(ctx, value)
}

type Unsigned struct {
	Min       uint64
	Max       uint64
	Validator ValidatorUnsigned
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

func (codec *Unsigned) IsEmpty(value interface{}) bool {
	val, _ := generic.ConvertToUint64(value)
	return val == 0
}

type ValidatorDecimal interface {
	Validate(ctx Context, value float64) error
}

type ValidatorFuncDecimal func(ctx Context, value float64) error

func (f ValidatorFuncDecimal) Validate(ctx Context, value float64) error {
	return f(ctx, value)
}

type Decimal struct {
	Min       float64
	Max       float64
	Validator ValidatorDecimal
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

func (codec *Decimal) IsEmpty(value interface{}) bool {
	val, _ := generic.ConvertToFloat64(value)
	return val == 0
}

// Abstract value formatter
type Formatter interface {
	Format(ctx Context, value interface{}) (val interface{}, err error)
}

// Formatter, that based on codec.Decode method
type BaseFormatter struct {
	Decoder
	ShowEmpty bool
}

func (w *BaseFormatter) Format(
	ctx Context,
	value interface{},
) (interface{}, error) {
	c := w.Decoder
	if c == nil {
		c = TextCodec
	}

	if !w.ShowEmpty {
		if cc, ok := c.(Empty); ok {
			if cc.IsEmpty(value) {
				return "", nil
			}
		}
	}

	val, err := c.Decode(ctx, value)
	if err != nil {
		return "", err
	}

	return val, nil
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
type PairConverter interface {
	PairCodec
	Formatter
}

type pairConverter struct {
	PairCodec
	Formatter
}

func NewPairConverter(codec PairCodec) PairConverter {
	return &pairConverter{
		PairCodec: codec,
		Formatter: &BaseFormatter{Decoder: codec},
	}
}

var (
	TextCodec     Codec = &Text{}
	SignedCodec   Codec = &Signed{}
	UnsignedCodec Codec = &Unsigned{}
	DecimalCodec  Codec = &Decimal{}
	BoolCodec     Codec = &Unsigned{} // Override in real application
)
