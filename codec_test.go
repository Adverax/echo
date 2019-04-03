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
			codec: &String{
				Validator: ValidatorStringFunc(
					func(ctx Context, value string) error {
						return nil
					},
				),
			},
			value:  "aaa",
			result: "aaa",
		},
		"TEXT: The value don't matched custom validator's must be rejected": {
			codec: &String{
				Validator: ValidatorStringFunc(
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

		// INT
		"INT: The invalid value must be rejected": {
			codec: &Int{},
			value: "abcd",
			error: true,
		},
		"INT: The valid value must be accepted": {
			codec:  &Int{},
			value:  "123",
			result: int(123),
		},
		"INT: The value less than allowed min must be rejected": {
			codec:  &Int{Min: 10, Max: 100},
			value:  "1",
			result: int(1),
			error:  true,
		},
		"INT: The value greater than allowed max must be rejected": {
			codec:  &Int{Min: 10, Max: 100},
			value:  "1000",
			result: int(1000),
			error:  true,
		},
		"INT: The value inside allowed range must be accepted": {
			codec:  &Int{Min: 10, Max: 100},
			value:  "50",
			result: int(50),
		},
		"INT: The value matched custom validator's must be accepted": {
			codec: &Int{
				Validator: ValidatorIntFunc(
					func(ctx Context, value int) error {
						return nil
					},
				),
			},
			value:  "123",
			result: int(123),
		},
		"INT: The value don't matched custom validator's must be rejected": {
			codec: &Int{
				Validator: ValidatorIntFunc(
					func(ctx Context, value int) error {
						return ValidationErrors{
							ValidationErrorInvalidValue,
						}
					},
				),
			},
			value: "123",
			error: true,
		},

		// UINT
		"UINT: The invalid value must be rejected": {
			codec: &Uint{},
			value: "abcd",
			error: true,
		},
		"UINT: The valid value must be accepted": {
			codec:  &Uint{},
			value:  "123",
			result: uint(123),
		},
		"UINT: The value less than allowed min must be rejected": {
			codec:  &Uint{Min: 10, Max: 100},
			value:  "1",
			result: uint(1),
			error:  true,
		},
		"UINT: The value greater than allowed max must be rejected": {
			codec:  &Uint{Min: 10, Max: 100},
			value:  "1000",
			result: uint(1000),
			error:  true,
		},
		"UINT: The value inside allowed range must be accepted": {
			codec:  &Uint{Min: 10, Max: 100},
			value:  "50",
			result: uint(50),
		},
		"UINT: The value matched custom validator's must be accepted": {
			codec: &Uint{
				Validator: ValidatorUintFunc(
					func(ctx Context, value uint) error {
						return nil
					},
				),
			},
			value:  "123",
			result: uint64(123),
		},
		"UINT: The value don't matched custom validator's must be rejected": {
			codec: &Uint{
				Validator: ValidatorUintFunc(
					func(ctx Context, value uint) error {
						return ValidationErrors{
							ValidationErrorInvalidValue,
						}
					},
				),
			},
			value: "123",
			error: true,
		},

		// FLOAT64
		"FLOAT64: The invalid value must be rejected": {
			codec: &Float64{},
			value: "abcd",
			error: true,
		},
		"FLOAT64: The valid value must be accepted": {
			codec:  &Float64{},
			value:  "123.5",
			result: float64(123.5),
		},
		"FLOAT64: The value less than allowed min must be rejected": {
			codec:  &Float64{Min: 10, Max: 100},
			value:  "1",
			result: float64(1),
			error:  true,
		},
		"FLOAT64: The value greater than allowed max must be rejected": {
			codec:  &Float64{Min: 10, Max: 100},
			value:  "1000",
			result: float64(1000),
			error:  true,
		},
		"FLOAT64: The value inside allowed range must be accepted": {
			codec:  &Float64{Min: 10, Max: 100},
			value:  "50",
			result: float64(50),
		},
		"FLOAT64: The value matched custom validator's must be accepted": {
			codec: &Float64{
				Validator: ValidatorFloat64Func(
					func(ctx Context, value float64) error {
						return nil
					},
				),
			},
			value:  "123",
			result: float64(123),
		},
		"FLOAT64: The value don't matched custom validator's must be rejected": {
			codec: &Float64{
				Validator: ValidatorFloat64Func(
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
