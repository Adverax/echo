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

		if model.HasErrors() {
			// Show form
			...
			return nil
		}
	}

	// Record valid
	...
	return nil
*/

// Abstract field. Implemented by descendands of field.
type ModelField interface {
	// Field has errors
	HasErrors() bool
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

type ModelValidator func(me *Model) error

type Model map[string]interface{}

func (model Model) Clone() Model {
	res := make(Model, 2*len(model))
	for key, val := range model {
		res[key] = val
	}
	return res
}

func (model Model) Bind(
	ctx Context,
	rec interface{},
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
		}*/

	ctype := req.Header.Get(HeaderContentType)
	switch {
	case strings.HasPrefix(ctype, MIMEApplicationJSON):
		if err := json.NewDecoder(req.Body).Decode(rec); err != nil {
			if ute, ok := err.(*json.UnmarshalTypeError); ok {
				return NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unmarshal type error: expected=%v, got=%v, field=%v, offset=%v", ute.Type, ute.Value, ute.Field, ute.Offset)).SetInternal(err)
			} else if se, ok := err.(*json.SyntaxError); ok {
				return NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Syntax error: offset=%v, error=%v", se.Offset, se.Error())).SetInternal(err)
			}
			return NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
		}
	case strings.HasPrefix(ctype, MIMEApplicationXML), strings.HasPrefix(ctype, MIMETextXML):
		if err := xml.NewDecoder(req.Body).Decode(rec); err != nil {
			if ute, ok := err.(*xml.UnsupportedTypeError); ok {
				return NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unsupported type error: type=%v, error=%v", ute.Type, ute.Error())).SetInternal(err)
			} else if se, ok := err.(*xml.SyntaxError); ok {
				return NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Syntax error: line=%v, error=%v", se.Line, se.Error())).SetInternal(err)
			}
			return NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
		}
	case strings.HasPrefix(ctype, MIMEApplicationForm), strings.HasPrefix(ctype, MIMEMultipartForm):
		params, err := ctx.FormParams()
		if err != nil {
			return NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
		}
		if err = model.BindFrom(ctx, params, rec); err != nil {
			return NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
		}
	default:
		return ErrUnsupportedMediaType
	}

	return nil
}

func (model Model) BindFrom(
	ctx Context,
	data map[string][]string,
	rec interface{},
) error {
	if rec == nil {
		// Skip record assignment
		for _, item := range model {
			if field, ok := item.(ModelField); ok {
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
		return nil
	}

	// Use record assignment
	record := reflect.ValueOf(rec)
	if record.Type().Kind() == reflect.Ptr {
		record = record.Elem()
	}

	for _, item := range model {
		if field, ok := item.(ModelField); ok {
			f := record.FieldByName(field.GetName())
			if f.Kind() == reflect.Invalid {
				continue
			}

			if f.CanInterface() {
				field.SetVal(ctx, f.Interface())
			}

			if f.CanSet() {
				value, ok := data[field.GetName()]
				if ok && len(value) != 0 {
					err := field.SetValue(ctx, value[0])
					if err != nil {
						return err
					}

					dst := f.Addr().Interface()
					generic.ConvertAssign(dst, field.GetVal())
				}
			}
		}
	}

	return nil
}

func (model Model) HasErrors() bool {
	for _, item := range model {
		if field, ok := item.(ModelField); ok {
			if field == nil {
				continue
			}
			if field.HasErrors() {
				return true
			}
		}
	}

	return false
}

type Models []Model

func (models Models) Bind(
	ctx Context,
	recs []interface{},
) error {
	for i, model := range models {
		err := model.Bind(ctx, recs[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (models Models) HasErrors() bool {
	for _, model := range models {
		if model != nil && model.HasErrors() {
			return true
		}
	}
	return false
}
