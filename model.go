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
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/adverax/echo/generic"
	"net/http"
	"reflect"
	"strings"
)

/*
Example:
func myLoginHandler(ctx echo.Context) error {
	username := XXX{} // Create field
	password := XXX{} // Create field

	model := echo.Model{username, password}
	rec := struct{
        Username string
        Password string
    }{
		Username: "Default name",
        Password: "Default password,
	}

	if ctx.Request == post {
		err := model.Bind(&data, &rec)
		if err != nil {
			return err
		}

		if model.IsValid() {
            // Record is valid
			...
			return nil
		}
	}

	// Show form
	...
	return nil
*/

type Mapper interface {
	Execute(name string) (string, bool)
}

type MapperFunc func(name string) (string, bool)

func (fn MapperFunc) Execute(name string) (string, bool) {
	return fn(name)
}

type DictMapper map[string]string

func (mapper DictMapper) Execute(name string) (string, bool) {
	res, ok := mapper[name]
	return res, ok
}

// Abstract field. Implemented by descendants of field.
type ModelField interface {
	// Field has no errors
	IsValid() bool
	// Get list of validation errors
	GetErrors() ValidationErrors
	// Get name of field
	GetName() string
	// Get internal representation of value
	GetVal() interface{}
	// Set internal representation pf value
	SetVal(ctx Context, value interface{})
	// Get external representation pf value
	GetValue() string
	// Set external representation pf value
	SetValue(ctx Context, value string) error
	// Get flag disabled
	GetDisabled() bool
	// Get flag hidden
	GetHidden() bool
	// Delete all errors and uncheck field and set default value
	Reset(ctx Context) error
}

type ValidatorFunc func() error

type ModelValidator func(me *Model) error

type Model map[string]interface{}

func (model Model) Clone() Model {
	res := make(Model, 2*len(model))
	for key, val := range model {
		res[key] = val
	}
	return res
}

// Bind works with not structured data only.
func (model Model) Bind(
	ctx Context,
	validators ...ValidatorFunc,
) error {
	req := ctx.Request()

	/*
		if req.ContentLength == 0 {
			if req.Method == http.MethodGet || req.Method == http.MethodDelete {
				if err = model.load(ctx, c.QueryParams(), "query"); err != nil {
					return NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
				}
				return
			}
			return NewHTTPError(http.StatusBadRequest, "Request body can't be empty")
		}
	*/

	var params map[string][]string
	ctype := req.Header.Get(HeaderContentType)
	switch {
	case strings.HasPrefix(ctype, MIMEApplicationJSON):
		var raw map[string]string
		if err := json.NewDecoder(req.Body).Decode(raw); err != nil {
			if ute, ok := err.(*json.UnmarshalTypeError); ok {
				return NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unmarshal type error: expected=%v, got=%v, field=%v, offset=%v", ute.Type, ute.Value, ute.Field, ute.Offset)).SetInternal(err)
			} else if se, ok := err.(*json.SyntaxError); ok {
				return NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Syntax error: offset=%v, error=%v", se.Offset, se.Error())).SetInternal(err)
			}
			return NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
		}
		params = MakeModelParams(raw)
	case strings.HasPrefix(ctype, MIMEApplicationXML), strings.HasPrefix(ctype, MIMETextXML):
		var raw map[string]string
		if err := xml.NewDecoder(req.Body).Decode(raw); err != nil {
			if ute, ok := err.(*xml.UnsupportedTypeError); ok {
				return NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unsupported type error: type=%v, error=%v", ute.Type, ute.Error())).SetInternal(err)
			} else if se, ok := err.(*xml.SyntaxError); ok {
				return NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Syntax error: line=%v, error=%v", se.Line, se.Error())).SetInternal(err)
			}
			return NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
		}
	case strings.HasPrefix(ctype, MIMEApplicationForm), strings.HasPrefix(ctype, MIMEMultipartForm):
		var err error
		params, err = ctx.FormParams()
		if err != nil {
			return NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
		}
	default:
		return ErrUnsupportedMediaType
	}

	if err := model.BindFrom(ctx, params, validators...); err != nil {
		return NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	return nil
}

func (model Model) BindFrom(
	ctx Context,
	data map[string][]string,
	validators ...ValidatorFunc,
) error {
	for _, item := range model {
		if field, ok := item.(ModelField); ok {
			err := field.Reset(ctx)
			if err != nil {
				return err
			}

			if field.GetDisabled() || field.GetHidden() {
				continue
			}

			name := field.GetName()
			value, ok := data[name]
			if ok && len(value) != 0 {
				err := field.SetValue(ctx, value[0])
				if err != nil {
					return err
				}
			}
		}
	}

	// Apply custom validators
	for _, validator := range validators {
		if validator == nil {
			continue
		}
		err := validator()
		if err != nil {
			return err
		}
	}

	return nil
}

func (model Model) IsValid() bool {
	for _, item := range model {
		if field, ok := item.(ModelField); ok {
			if field == nil {
				continue
			}
			if !field.IsValid() {
				return false
			}
		}
	}

	return true
}

func (model Model) AssignFrom(
	ctx Context,
	src interface{},
	mapper Mapper,
) error {
	if mapper == nil {
		mapper = DefaultMapper
	}
	rec := reflect.ValueOf(src)
	if rec.Kind() == reflect.Ptr {
		rec = rec.Elem()
	}
	if rec.Kind() != reflect.Struct {
		return fmt.Errorf("invalid type of source")
	}

	for _, item := range model {
		if field, ok := item.(ModelField); ok {
			if field == nil {
				continue
			}
			if name, ok := mapper.Execute(field.GetName()); ok {
				f := rec.FieldByName(name)
				if f.Kind() != reflect.Invalid {
					if f.CanInterface() {
						field.SetVal(ctx, f.Interface())
					}
				}
			}
		}
	}

	return nil
}

func (model Model) AssignTo(
	ctx Context,
	dst interface{},
	mapper Mapper,
) error {
	if mapper == nil {
		mapper = DefaultMapper
	}

	rec := reflect.ValueOf(dst).Elem()
	if rec.Kind() != reflect.Struct {
		return fmt.Errorf("invalid type of source")
	}

	for _, item := range model {
		if field, ok := item.(ModelField); ok {
			if field == nil {
				continue
			}
			if name, ok := mapper.Execute(field.GetName()); ok {
				f := rec.FieldByName(name)
				if f.Kind() != reflect.Invalid {
					if f.CanSet() {
						dst := f.Addr().Interface()
						_ = generic.ConvertAssign(dst, field.GetVal())
					}
				}
			}
		}
	}

	return nil
}

type Models []Model

func (models Models) Bind(
	ctx Context,
) error {
	for _, model := range models {
		err := model.Bind(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (models Models) IsValid() bool {
	for _, model := range models {
		if model != nil && model.IsValid() {
			return true
		}
	}
	return false
}

// Create name of field of band
func MakeMultiModelName(key, name string) string {
	return fmt.Sprintf("[%s].%s", key, name)
}

// Create model parameters from map.
func MakeModelParams(raw map[string]string) map[string][]string {
	params := make(map[string][]string, len(raw))
	for k, v := range raw {
		params[k] = []string{v}
	}
	return params
}
