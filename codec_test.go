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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Encode(t *testing.T) {
	type Test struct {
		codec  Codec
		value  string
		result interface{}
		error  bool
	}

	tests := map[string]Test{
		// TEXT
		"TEXT: The value matched custom validator's must be accepted": {
			codec: &Text{
				Validator: ValidatorTextFunc(
					func(ctx Context, value string) error {
						return nil
					},
				),
			},
			value:  "aaa",
			result: "aaa",
		},
		"TEXT: The value don't matched custom validator's must be rejected": {
			codec: &Text{
				Validator: ValidatorTextFunc(
					func(ctx Context, value string) error {
						return ValidationErrors{
							ValidationErrorInvalidValue,
						}
					},
				),
			},
			value: "aaa",
			error: true,
		},

		// SIGNED
		"SIGNED: The invalid value must be rejected": {
			codec: &Signed{},
			value: "abcd",
			error: true,
		},
		"SIGNED: The valid value must be accepted": {
			codec:  &Signed{},
			value:  "123",
			result: int64(123),
		},
		"SIGNED: The value less than allowed min must be rejected": {
			codec: &Signed{Min: 10, Max: 100},
			value: "1",
			error: true,
		},
		"SIGNED: The value greater than allowed max must be rejected": {
			codec: &Signed{Min: 10, Max: 100},
			value: "1000",
			error: true,
		},
		"SIGNED: The value inside allowed range must be accepted": {
			codec:  &Signed{Min: 10, Max: 100},
			value:  "50",
			result: int64(50),
		},
		"SIGNED: The value matched custom validator's must be accepted": {
			codec: &Signed{
				Validator: ValidatorSignedFunc(
					func(ctx Context, value int64) error {
						return nil
					},
				),
			},
			value:  "123",
			result: int64(123),
		},
		"SIGNEDThe value don't matched custom validator's must be rejected": {
			codec: &Signed{
				Validator: ValidatorSignedFunc(
					func(ctx Context, value int64) error {
						return ValidationErrors{
							ValidationErrorInvalidValue,
						}
					},
				),
			},
			value: "123",
			error: true,
		},

		// UNSIGNED
		"UNSIGNED: The invalid value must be rejected": {
			codec: &Unsigned{},
			value: "abcd",
			error: true,
		},
		"UNSIGNED: The valid value must be accepted": {
			codec:  &Unsigned{},
			value:  "123",
			result: uint64(123),
		},
		"UNSIGNED: The value less than allowed min must be rejected": {
			codec: &Unsigned{Min: 10, Max: 100},
			value: "1",
			error: true,
		},
		"UNSIGNED: The value greater than allowed max must be rejected": {
			codec: &Unsigned{Min: 10, Max: 100},
			value: "1000",
			error: true,
		},
		"UNSIGNED: The value inside allowed range must be accepted": {
			codec:  &Unsigned{Min: 10, Max: 100},
			value:  "50",
			result: uint64(50),
		},
		"UNSIGNED: The value matched custom validator's must be accepted": {
			codec: &Unsigned{
				Validator: ValidatorUnsignedFunc(
					func(ctx Context, value uint64) error {
						return nil
					},
				),
			},
			value:  "123",
			result: uint64(123),
		},
		"UNSIGNED: The value don't matched custom validator's must be rejected": {
			codec: &Unsigned{
				Validator: ValidatorUnsignedFunc(
					func(ctx Context, value uint64) error {
						return ValidationErrors{
							ValidationErrorInvalidValue,
						}
					},
				),
			},
			value: "123",
			error: true,
		},

		// DECIMAL
		"DECIMAL: The invalid value must be rejected": {
			codec: &Decimal{},
			value: "abcd",
			error: true,
		},
		"DECIMAL: The valid value must be accepted": {
			codec:  &Decimal{},
			value:  "123.5",
			result: float64(123.5),
		},
		"DECIMAL: The value less than allowed min must be rejected": {
			codec: &Decimal{Min: 10, Max: 100},
			value: "1",
			error: true,
		},
		"DECIMAL: The value greater than allowed max must be rejected": {
			codec: &Decimal{Min: 10, Max: 100},
			value: "1000",
			error: true,
		},
		"DECIMAL: The value inside allowed range must be accepted": {
			codec:  &Decimal{Min: 10, Max: 100},
			value:  "50",
			result: float64(50),
		},
		"DECIMAL: The value matched custom validator's must be accepted": {
			codec: &Decimal{
				Validator: ValidatorDecimalFunc(
					func(ctx Context, value float64) error {
						return nil
					},
				),
			},
			value:  "123",
			result: float64(123),
		},
		"DECIMAL: The value don't matched custom validator's must be rejected": {
			codec: &Decimal{
				Validator: ValidatorDecimalFunc(
					func(ctx Context, value float64) error {
						return ValidationErrors{
							ValidationErrorInvalidValue,
						}
					},
				),
			},
			value: "123",
			error: true,
		},

		// OPTIONAL
		"OPTIONAL: The valid value must be accepted": {
			codec: &Optional{
				Codec: &Signed{},
			},
			value:  "123",
			result: int64(123),
		},
		"OPTIONAL: The default value must be accepted": {
			codec: &Optional{
				Codec: &Signed{},
			},
			value:  "123",
			result: int64(0),
		},
	}

	e := New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			result, err := test.codec.Encode(c, test.value)
			if test.error {
				assert.Error(t, err)
			} else {
				if assert.NoError(t, err) {
					assert.Equal(t, test.result, result)
				}
			}
		})
	}
}
