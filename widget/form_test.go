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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFormComponent_Render(t *testing.T) {
	type Test struct {
		src echo.ModelField
		val interface{}
		dst string
	}

	cities := echo.NewDataSet(
		map[string]string{
			"1": "London",
			"2": "Paris",
		},
		true,
	)

	tests := map[string]Test{
		"FormHidden": {
			src: &FormText{
				Name: "Challenge",
			},
			val: "1234567890",
			dst: `{"Name":"Challenge","Value":"1234567890"}`,
		},

		"FormText: simple": {
			src: &FormText{
				Name:  "Comments",
				Label: "Comments",
			},
			dst: `{"Label":"Comments","Name":"Comments"}`,
		},
		"FormText: complex": {
			src: &FormText{
				Name:        "Comments",
				Label:       "Comments",
				Disabled:    true,
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

		"FormSelect: simple": {
			src: &FormSelect{
				Name:     "City",
				Label:    "City",
				Disabled: false,
				Required: true,
				Items:    cities,
			},
			val: 2,
			dst: `{"Items":[{"Label":"London","Value":"1"},{"Label":"Paris","Selected":true,"Value":"2"}],"Label":"City","Name":"City","Required":true}`,
		},
		"FormSelect: complex": {
			src: &FormSelect{
				Name:     "City",
				Label:    "City",
				Disabled: false,
				Required: false,
				Items:    cities,
			},
			dst: `{"Empty":{"Label":"(Empty)"},"Items":[{"Label":"London","Value":"1"},{"Label":"Paris","Value":"2"}],"Label":"City","Name":"City"}`,
		},

		"FormFlag": {
			src: &FormFlag{
				Name:  "Option",
				Label: "Input option",
			},
			val: true,
			dst: `{"Label":"Input option","Name":"Option","Selected":true,"Value":"1"}`,
		},

		"FormFlags": {
			src: &FormFlags{
				Name:  "Option",
				Label: "Input option",
				Items: cities,
			},
			val: []string{"1", "2"},
			dst: `{"Items":[{"Label":"London","Selected":true,"Value":"1"},{"Label":"Paris","Selected":true,"Value":"2"}],"Label":"Input option","Name":"Option"}`,
		},

		"FormSubmit: simple": {
			src: &FormSubmit{
				Label: "Accept",
			},
			dst: `{"Label":"Accept"}`,
		},

		"FormFile": {
			src: &FormFile{
				Name:   "File",
				Label:  "File",
				Accept: "*.doc",
			},
			dst: `{"Accept":"*.doc","Label":"File","Name":"File"}`,
		},
	}

	e := echo.New()
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			c := e.NewContext(nil, nil)
			if test.val != nil {
				test.src.SetVal(c, test.val)
			}
			tree, err := echo.RenderWidget(c, test.src)
			require.NoError(t, err)
			dst, err := json.Marshal(tree)
			require.NoError(t, err)
			require.Equal(t, test.dst, string(dst))
		})
	}
}

func TestFormComponent_SetValueAndValidate(t *testing.T) {
	type Test struct {
		// input
		field echo.ModelField
		src   []string
		// Output
		dst    []string
		val    interface{}
		errors echo.ValidationErrors
	}

	cities := echo.NewDataSet(
		map[string]string{
			"1": "London",
			"2": "Paris",
		},
		true,
	)

	tests := map[string]Test{
		"FormText: Default must be accepted": {
			field: &FormText{},
			src:   []string{"123"},
			dst:   []string{"123"},
			val:   "123",
		},

		"FormText: Required value with data must be accepted": {
			field: &FormText{
				Required: true,
			},
			src: []string{"123"},
			dst: []string{"123"},
			val: "123",
		},

		"FormText: Required value without data must be rejected": {
			field: &FormText{
				Required: true,
			},
			dst: []string{""},
			val: "",
			errors: echo.ValidationErrors{
				MessageConstraintRequired,
			},
		},

		"FormText: Value matched pattern must be accepted": {
			field: &FormText{
				Pattern: "[a-z]+",
			},
			src: []string{"abc"},
			dst: []string{"abc"},
			val: "abc",
		},

		"FormText: Value don't matched pattern must be rejected": {
			field: &FormText{
				Pattern: "[a-z]+",
			},
			src: []string{"123"},
			dst: []string{"123"},
			val: "123",
			errors: echo.ValidationErrors{
				MessageConstraintPattern,
			},
		},

		"FormText: Value with length smaller than allowed must be accepted": {
			field: &FormText{
				MaxLength: 64,
			},
			src: []string{"12345"},
			dst: []string{"12345"},
			val: "12345",
		},

		"FormText: Too length value must produce error": {
			field: &FormText{
				MaxLength: 3,
			},
			src: []string{"12345"},
			dst: []string{"12345"},
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
			src:   []string{"123"},
			dst:   []string{"123"},
			val:   "123",
		},

		"FormHidden: Required value with data must be accepted": {
			field: &FormHidden{
				Required: true,
			},
			src: []string{"123"},
			dst: []string{"123"},
			val: "123",
		},

		"FormHidden: Required value without data must be rejected": {
			field: &FormHidden{
				Required: true,
			},
			dst: []string{""},
			val: "",
			errors: echo.ValidationErrors{
				MessageConstraintRequired,
			},
		},

		"FormHidden: Value matched pattern must be accepted": {
			field: &FormHidden{
				Pattern: "[a-z]+",
			},
			src: []string{"abc"},
			dst: []string{"abc"},
			val: "abc",
		},

		"FormHidden: Value don't matched pattern must be rejected": {
			field: &FormHidden{
				Pattern: "[a-z]+",
			},
			src: []string{"123"},
			dst: []string{"123"},
			val: "123",
			errors: echo.ValidationErrors{
				MessageConstraintPattern,
			},
		},

		"FormHidden: Value with length smaller than allowed must be accepted": {
			field: &FormHidden{
				MaxLength: 64,
			},
			src: []string{"12345"},
			dst: []string{"12345"},
			val: "12345",
		},

		"FormHidden: Too length value must produce error": {
			field: &FormHidden{
				MaxLength: 3,
			},
			src: []string{"12345"},
			dst: []string{"12345"},
			val: "12345",
			errors: echo.ValidationErrors{
				&echo.Cause{
					Msg:  uint32(MessageConstraintMaxLength),
					Args: []interface{}{3},
				},
			},
		},

		"FormSelect: Default": {
			field: &FormSelect{
				Items: cities,
			},
			src: []string{"1"},
			dst: []string{"1"},
			val: "1",
		},
		"FormSelect: Required value with data must be accepted": {
			field: &FormSelect{
				Items:    cities,
				Required: true,
			},
			src: []string{"1"},
			dst: []string{"1"},
			val: "1",
		},
		"FormSelect: Required value without data must be rejected": {
			field: &FormSelect{
				Items:    cities,
				Required: true,
			},
			dst: []string{""},
			val: "",
			errors: echo.ValidationErrors{
				MessageConstraintRequired,
			},
		},

		"FormMultiSelect: Normal": {
			field: &FormFlags{
				Items: cities,
			},
			src: []string{"1", "100"},
			dst: []string{"1", "100"},
			val: []string{"1"},
		},
	}

	e := echo.New()
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			//t.Parallel()
			c := e.NewContext(nil, nil)
			err := test.field.SetValue(c, test.src)
			require.NoError(t, err)
			err = test.field.Validate(c)
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

	username := &FormText{
		Name: "Username",
	}

	model := echo.Model{"Username": username}
	err := model.BindFrom(
		ctx,
		map[string][]string{
			"Username": {"Bob"},
		},
	)
	assert.NoError(t, err)
	assert.Equal(t, "Bob", username.GetString())
}

func TestModel_AssignFrom(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	ctx := e.NewContext(req, httptest.NewRecorder())

	username := &FormText{
		Name: "username",
	}

	region := &FormText{
		Name: "region",
	}

	model := echo.Model{
		"username": username,
		"region":   region,
	}

	var rec = struct {
		Username string
		Password string
	}{
		Username: "Bob",
	}

	err := model.AssignFrom(ctx, rec, nil)
	assert.NoError(t, err)
	assert.Equal(t, "Bob", username.GetString())
}

func TestModel_AssignTo(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	ctx := e.NewContext(req, httptest.NewRecorder())

	username := &FormText{
		Name: "Username",
	}

	region := &FormText{
		Name: "Region",
	}

	model := echo.Model{
		"username": username,
		"region":   region,
	}

	username.SetVal(ctx, "Bob")

	var rec struct {
		Username string
		Password string
	}

	err := model.AssignTo(ctx, &rec, nil)
	assert.NoError(t, err)
	assert.Equal(t, "Bob", rec.Username)
}
