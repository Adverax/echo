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

package widget

import (
	"encoding/json"
	"github.com/adverax/echo"
	"github.com/adverax/echo/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRenderFormElement(t *testing.T) {
	type Test struct {
		src echo.ModelField
		val interface{}
		dst string
	}

	cities := data.NewSet(
		map[string]string{
			"1": "London",
			"2": "Paris",
		},
		true,
	)

	tests := map[string]Test{
		"FormHidden": {
			src: &FormTextArea{
				FormField: FormField{
					Name: "Challenge",
				},
			},
			val: "1234567890",
			dst: `{"Name":"Challenge","Value":"1234567890"}`,
		},

		"FormTextInput: simple": {
			src: &FormTextInput{
				FormField: FormField{
					Name:  "Username",
					Label: "User name",
				},
			},
			dst: `{"Label":"User name","Name":"Username"}`,
		},
		"FormTextInput: complex": {
			src: &FormTextInput{
				FormField: FormField{
					Name:     "Username",
					Label:    "User name",
					Disabled: true,
				},
				Required:    true,
				Pattern:     `[A-Za-z]{3,64}`,
				Placeholder: "Input your name",
				MaxLength:   64,
			},
			val: "Tom",
			dst: `{"Disabled":true,"Label":"User name","MaxLen":64,"Name":"Username","Pattern":"[A-Za-z]{3,64}","Placeholder":"Input your name","Required":true,"Value":"Tom"}`,
		},

		"FormTextArea: simple": {
			src: &FormTextArea{
				FormField: FormField{
					Name:  "Comments",
					Label: "Comments",
				},
			},
			dst: `{"Label":"Comments","Name":"Comments"}`,
		},
		"FormTextArea: complex": {
			src: &FormTextArea{
				FormField: FormField{
					Name:     "Comments",
					Label:    "Comments",
					Disabled: true,
				},
				Required:    true,
				Pattern:     `[A-Za-z ]+`,
				Placeholder: "Enter your comments",
				MaxLength:   1024,
				ReadOnly:    true,
				Rows:        16,
			},
			val: "My comments",
			dst: `{"Disabled":true,"Label":"Comments","MaxLen":1024,"Name":"Comments","Pattern":"[A-Za-z ]+","Placeholder":"Enter your comments","Readonly":true,"Required":true,"Rows":16,"Value":"My comments"}`,
		},

		"FormFileInput": {
			src: &FormFileInput{
				FormField: FormField{
					Name:  "File",
					Label: "File",
				},
				Accept: "*.doc",
			},
			dst: `{"Accept":"*.doc","Label":"File","Name":"File"}`,
		},

		"FormSelector: simple": {
			src: &FormSelector{
				FormField: FormField{
					Name:     "City",
					Label:    "City",
					Disabled: false,
				},
				Required: true,
				Items:    cities,
			},
			val: 2,
			dst: `{"Items":[{"Label":"London","Value":"1"},{"Label":"Paris","Selected":true,"Value":"2"}],"Label":"City","Name":"City","Required":true}`,
		},
		"FormSelector: complex": {
			src: &FormSelector{
				FormField: FormField{
					Name:     "City",
					Label:    "City",
					Disabled: false,
				},
				Required: false,
				Items:    cities,
			},
			dst: `{"Items":[{"Label":"(Empty)","Value":""},{"Label":"London","Value":"1"},{"Label":"Paris","Value":"2"}],"Label":"City","Name":"City"}`,
		},

		"FormCheckBox": {
			src: &FormCheckBox{
				FormField: FormField{
					Name:  "Gender",
					Label: "Mail",
					Codec: CheckBoxCodec,
				},
			},
			val: true,
			dst: `{"Checked":true,"Label":"Mail","Name":"Gender"}`,
		},

		"FormSubmit": {
			src: &FormSubmit{
				FormField: FormField{
					Name:  "Accept",
					Label: "Accept",
				},
			},
			val: "accept",
			dst: `{"Label":"Accept","Name":"Accept","Value":"accept"}`,
		},
	}

	e := echo.New()
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			c := e.NewContext(nil, nil)
			if test.val != nil {
				test.src.SetVal(c, test.val)
			}
			tree, err := RenderWidget(c, test.src)
			require.NoError(t, err)
			dst, err := json.Marshal(tree)
			require.NoError(t, err)
			require.Equal(t, test.dst, string(dst))
		})
	}
}

func TestFormComponent_SetValue(t *testing.T) {
	type Test struct {
		// input
		field echo.ModelField
		src   string
		// Output
		dst    string
		val    interface{}
		errors echo.ValidationErrors
	}

	cities := data.NewSet(
		map[string]string{
			"1": "London",
			"2": "Paris",
		},
		true,
	)

	tests := map[string]Test{
		"FormTextInput: Default must be accepted": {
			field: &FormTextInput{},
			src:   "123",
			dst:   "123",
			val:   "123",
		},

		"FormTextInput: Required value with data must be accepted": {
			field: &FormTextInput{
				Required: true,
			},
			src: "123",
			dst: "123",
			val: "123",
		},

		"FormTextInput: Required value without data must be rejected": {
			field: &FormTextInput{
				Required: true,
			},
			val: "",
			errors: echo.ValidationErrors{
				MessageConstraintRequired,
			},
		},

		"FormTextInput: Value matched pattern must be accepted": {
			field: &FormTextInput{
				Pattern: "[a-z]+",
			},
			src: "abc",
			dst: "abc",
			val: "abc",
		},

		"FormTextInput: Value don't matched pattern must be rejected": {
			field: &FormTextInput{
				Pattern: "[a-z]+",
			},
			src: "123",
			dst: "123",
			val: "123",
			errors: echo.ValidationErrors{
				MessageConstraintPattern,
			},
		},

		"FormTextInput: Value with length smaller than allowed must be accepted": {
			field: &FormTextInput{
				MaxLength: 64,
			},
			src: "12345",
			dst: "12345",
			val: "12345",
		},

		"FormTextInput: Too length value must produce error": {
			field: &FormTextInput{
				MaxLength: 3,
			},
			src: "12345",
			dst: "12345",
			val: "12345",
			errors: echo.ValidationErrors{
				&echo.Cause{
					Msg:  uint32(MessageConstraintMaxLength),
					Args: []interface{}{3},
				},
			},
		},

		"FormTextArea: Default must be accepted": {
			field: &FormTextArea{},
			src:   "123",
			dst:   "123",
			val:   "123",
		},

		"FormTextArea: Required value with data must be accepted": {
			field: &FormTextArea{
				Required: true,
			},
			src: "123",
			dst: "123",
			val: "123",
		},

		"FormTextArea: Required value without data must be rejected": {
			field: &FormTextArea{
				Required: true,
			},
			val: "",
			errors: echo.ValidationErrors{
				MessageConstraintRequired,
			},
		},

		"FormTextArea: Value matched pattern must be accepted": {
			field: &FormTextArea{
				Pattern: "[a-z]+",
			},
			src: "abc",
			dst: "abc",
			val: "abc",
		},

		"FormTextArea: Value don't matched pattern must be rejected": {
			field: &FormTextArea{
				Pattern: "[a-z]+",
			},
			src: "123",
			dst: "123",
			val: "123",
			errors: echo.ValidationErrors{
				MessageConstraintPattern,
			},
		},

		"FormTextArea: Value with length smaller than allowed must be accepted": {
			field: &FormTextArea{
				MaxLength: 64,
			},
			src: "12345",
			dst: "12345",
			val: "12345",
		},

		"FormTextArea: Too length value must produce error": {
			field: &FormTextArea{
				MaxLength: 3,
			},
			src: "12345",
			dst: "12345",
			val: "12345",
			errors: echo.ValidationErrors{
				&echo.Cause{
					Msg:  uint32(MessageConstraintMaxLength),
					Args: []interface{}{3},
				},
			},
		},

		"FormHidden: Default must be accepted": {
			field: &FormHidden{},
			src:   "123",
			dst:   "123",
			val:   "123",
		},

		"FormHidden: Required value with data must be accepted": {
			field: &FormHidden{
				Required: true,
			},
			src: "123",
			dst: "123",
			val: "123",
		},

		"FormHidden: Required value without data must be rejected": {
			field: &FormHidden{
				Required: true,
			},
			val: "",
			errors: echo.ValidationErrors{
				MessageConstraintRequired,
			},
		},

		"FormHidden: Value matched pattern must be accepted": {
			field: &FormHidden{
				Pattern: "[a-z]+",
			},
			src: "abc",
			dst: "abc",
			val: "abc",
		},

		"FormHidden: Value don't matched pattern must be rejected": {
			field: &FormHidden{
				Pattern: "[a-z]+",
			},
			src: "123",
			dst: "123",
			val: "123",
			errors: echo.ValidationErrors{
				MessageConstraintPattern,
			},
		},

		"FormHidden: Value with length smaller than allowed must be accepted": {
			field: &FormHidden{
				MaxLength: 64,
			},
			src: "12345",
			dst: "12345",
			val: "12345",
		},

		"FormHidden: Too length value must produce error": {
			field: &FormHidden{
				MaxLength: 3,
			},
			src: "12345",
			dst: "12345",
			val: "12345",
			errors: echo.ValidationErrors{
				&echo.Cause{
					Msg:  uint32(MessageConstraintMaxLength),
					Args: []interface{}{3},
				},
			},
		},

		"FormSelector: Default": {
			field: &FormSelector{
				Items: cities,
			},
			src: "1",
			dst: "1",
			val: "1",
		},
		"FormSelector: Required value with data must be accepted": {
			field: &FormSelector{
				Items:    cities,
				Required: true,
			},
			src: "1",
			dst: "1",
			val: "1",
		},
		"FormSelector: Required value without data must be rejected": {
			field: &FormSelector{
				Items:    cities,
				Required: true,
			},
			errors: echo.ValidationErrors{
				MessageConstraintRequired,
			},
		},

		"FormCheckBox: On": {
			field: &FormCheckBox{
				FormField: FormField{
					Codec: CheckBoxCodec,
				},
			},
			src: "on",
			dst: "on",
			val: true,
		},
		"FormCheckBox: Off": {
			field: &FormCheckBox{
				FormField: FormField{
					Codec: CheckBoxCodec,
				},
			},
			src: "off",
			dst: "off",
			val: false,
		},

		"FormSubmit: Default": {
			field: &FormSubmit{},
			val:   "",
		},
		/*"FormSubmit: Simple valid value must be accepted": {
			field: &FormSubmit{
				Field: echo.Field{
					Name: "Accept",
				},
				Required: true,
			},
			src: "accept",
			dst: "accept",
			val: "accept",
		},*/
		"FormSubmit: Simple invalid value must be rejected": {
			field: &FormSubmit{
				FormField: FormField{
					Name: "Accept",
				},
				Required: true,
			},
			src: "unknown",
			errors: echo.ValidationErrors{
				echo.ValidationErrorInvalidValue,
			},
		},
		"FormSubmit: Complex valid value must be accepted": {
			field: &FormSubmit{
				FormField: FormField{
					Name: "Accept",
				},
				Required: true,
				Items:    cities,
			},
			src: "1",
			dst: "1",
			val: "1",
		},
		"FormSubmit: Complex invalid value must be rejected": {
			field: &FormSubmit{
				FormField: FormField{
					Name: "Accept",
				},
				Required: true,
				Items:    cities,
			},
			src: "unknown",
			errors: echo.ValidationErrors{
				echo.ValidationErrorInvalidValue,
			},
		},
	}

	e := echo.New()
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			c := e.NewContext(nil, nil)
			err := test.field.SetValue(c, test.src)
			require.NoError(t, err)
			require.Equal(t, test.dst, test.field.GetValue(), "Invalid external representation")
			require.Equal(t, test.val, test.field.GetVal(), "Invalid internal representation")
			require.Equal(t, test.errors, test.field.GetErrors(), "Invalid errors")
		})
	}
}

func TestModel_Bind(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	ctx := e.NewContext(req, httptest.NewRecorder())

	rec := struct {
		Username string
	}{
		Username: "Default",
	}

	username := &FormField{
		Name: "Username",
	}

	model := echo.Model{"Username": username}
	err := model.BindFrom(
		ctx,
		map[string][]string{
			"Username": {"Bob"},
		},
		&rec,
	)
	assert.NoError(t, err)
	assert.Equal(t, "Bob", rec.Username)
}
