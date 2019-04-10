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

    if err := model.Resolve(ctx, rec, nil); err != nil {
        if err != echo.ErrModelSealed {
            return err
        }

		// Record is valid
		err := model.AssignTo(ctx, &rec)
		if err != nil {
			return err
		}
		...
		return ctx.Redirect(http.StatusSeeOther, "...")
    }

	// Show form
	...
	return nil
*/

type Mapper interface {
	// Convert external representation to internal representation
	Execute(name string) (string, bool)
}

type MapperFunc func(name string) (string, bool)

func (fn MapperFunc) Execute(name string) (string, bool) {
	return fn(name)
}

type ListMapper []string

func (mapper ListMapper) Execute(name string) (string, bool) {
	name = strings.Title(name)
	for _, val := range mapper {
		if val == name {
			return name, true
		}
	}
	return "", false
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
	// Set internal representation of value
	SetVal(ctx Context, value interface{})
	// Get external representation of value
	GetValue() []string
	// Set external representation of value
	SetValue(ctx Context, value []string) error
	// Validate field and extends field errors
	Validate(ctx Context) error
	// Get internal data as signed value
	GetInt() int
	// Get internal data as signed value
	GetInt8() int8
	// Get internal data as signed value
	GetInt16() int16
	// Get internal data as signed value
	GetInt32() int32
	// Get internal data as signed value
	GetInt64() int64
	// Get internal data as unsigned value
	GetUint() uint
	// Get internal data as unsigned value
	GetUint8() uint8
	// Get internal data as unsigned value
	GetUint16() uint16
	// Get internal data as unsigned value
	GetUint32() uint32
	// Get internal data as unsigned value
	GetUint64() uint64
	// Get internal data as decimal value
	GetFloat32() float32
	// Get internal data as decimal value
	GetFloat64() float64
	// Get internal data as string value
	GetString() string
	// Get internal data as boolean value
	GetBoolean() bool
	// Get flag disabled
	GetDisabled() bool
	// Get flag hidden
	GetHidden() bool
	// Delete all errors and uncheck field and set default value
	Reset(ctx Context) error
}

type ValidatorFunc func() error

type Model map[string]interface{}

func (model Model) Clone() Model {
	res := make(Model, 2*len(model))
	for key, val := range model {
		res[key] = val
	}
	return res
}

// Import and validate data
// Returns ErrModelSealed if model imported and validated.
func (model Model) Resolve(
	ctx Context,
	src interface{}, // Optional data source
	mapper Mapper, // Optional mapper
) error {
	if src != nil {
		err := model.AssignFrom(ctx, src, mapper)
		if err != nil {
			return err
		}
	}

	if ctx.Request().Method != POST {
		return nil
	}

	err := model.Bind(ctx)
	if err != nil {
		return err
	}

	if model.IsValid() {
		return ErrModelSealed
	}

	return nil
}

// Bind works with not structured data only.
func (model Model) Bind(
	ctx Context,
) error {
	req := ctx.Request()

	if req.ContentLength == 0 {
		if req.Method == http.MethodGet || req.Method == http.MethodDelete {
			if err := model.BindFrom(ctx, ctx.QueryParams()); err != nil {
				return NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
			}
			return nil
		}
		return NewHTTPError(http.StatusBadRequest, "Request body can't be empty")
	}

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

	if err := model.BindFrom(ctx, params); err != nil {
		return NewHTTPError(http.StatusBadRequest, err.Error()).SetInternal(err)
	}

	return nil
}

func (model Model) BindFrom(
	ctx Context,
	data map[string][]string,
) error {
	// Bind model
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
				err := field.SetValue(ctx, value)
				if err != nil {
					return err
				}
			}

			err = field.Validate(ctx)
			if err != nil {
				return err
			}
		}
	}

	// Validate model
	for _, item := range model {
		if validator, ok := item.(ValidatorFunc); ok {
			err := validator()
			if err != nil {
				return err
			}
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
				f := access(rec, name)
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
		return fmt.Errorf("invalid type of destination")
	}

	for _, item := range model {
		if field, ok := item.(ModelField); ok {
			if field == nil {
				continue
			}
			if name, ok := mapper.Execute(field.GetName()); ok {
				f := access(rec, name)
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

func (model Model) Render(
	ctx Context,
) (interface{}, error) {
	res := make(map[string]interface{}, len(model)+1)

	for key, item := range model {
		if item == nil {
			continue
		}

		f, err := RenderWidget(ctx, item)
		if err != nil {
			return nil, err
		}

		res[key] = f
	}

	return res, nil
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

func (models Models) Render(
	ctx Context,
) (interface{}, error) {
	res := make([]interface{}, 0, len(models))
	for _, model := range models {
		if model != nil {
			item, err := model.Render(ctx)
			if err != nil {
				return nil, err
			}
			if item != nil {
				res = append(res, item)
			}
		}
	}

	if len(res) == 0 {
		return nil, nil
	}

	return res, nil
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

// Access to field
func access(rec reflect.Value, name string) reflect.Value {
	if !strings.Contains(name, ".") {
		return rec.FieldByName(name)
	}

	path := strings.Split(name, ".")
	for _, item := range path {
		rec = rec.FieldByName(item)
		if rec.Kind() != reflect.Struct {
			return rec
		}
	}

	return rec
}
