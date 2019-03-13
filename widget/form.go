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
	"github.com/adverax/echo"
	"github.com/adverax/echo/data"
	"github.com/adverax/echo/generic"
	"io"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"
)

// MultiForm contains list of models.
type MultiForm struct {
	Id     string      // Form identifier
	Name   string      // Form name
	Method string      // Form method
	Action interface{} // Form action
	Models echo.Models // Primary form model
	Hidden bool        // Form is hidden
}

func (w *MultiForm) Render(
	ctx echo.Context,
) (interface{}, error) {
	if w.Hidden {
		return nil, nil
	}

	res := make(map[string]interface{}, 8)

	if len(w.Models) != 0 {
		model, err := RenderModels(ctx, w.Models)
		if err != nil {
			return nil, err
		}
		res["Model"] = model
	}

	if w.Id != "" {
		res["Id"] = w.Id
	}

	if w.Name != "" {
		res["Name"] = w.Name
	}

	if w.Action != nil {
		action, err := RenderLink(ctx, w.Action)
		if err != nil {
			return nil, err
		}
		res["Action"] = action
	}

	if w.Method == "" {
		res["Method"] = echo.POST
	} else {
		res["Method"] = w.Method
	}

	return res, nil
}

// Form based on a Model
type Form struct {
	Id     string      // Form identifier
	Name   string      // Form name
	Method string      // Form method
	Action interface{} // Form action
	Model  echo.Model  // Primary form model
	Hidden bool        // Form is hidden
}

func (w *Form) Render(
	ctx echo.Context,
) (interface{}, error) {
	if w.Hidden {
		return nil, nil
	}

	res := make(map[string]interface{}, 8)

	if w.Model != nil {
		model, err := RenderModel(ctx, w.Model)
		if err != nil {
			return nil, err
		}
		res["Model"] = model
	}

	if w.Id != "" {
		res["Id"] = w.Id
	}

	if w.Name != "" {
		res["Name"] = w.Name
	}

	if w.Action != nil {
		action, err := RenderLink(ctx, w.Action)
		if err != nil {
			return nil, err
		}
		res["Action"] = action
	}

	if w.Method == "" {
		res["Method"] = echo.POST
	} else {
		res["Method"] = w.Method
	}

	return res, nil
}

type FormFieldFilterFunc func(value string) string

type FormField struct {
	Id       string                // Field identifier
	Name     string                // Field name
	Label    interface{}           // Field label
	Disabled bool                  // Field disabled
	Hidden   bool                  // Field is hidden (not rendered)
	Filter   FormFieldFilterFunc   // Custom filter
	Codec    echo.Codec            // Field codec (optional)
	val      interface{}           // Internal representation of value
	value    string                // External representation of value
	errors   echo.ValidationErrors // Field errors
}

func (field *FormField) GetName() string {
	return field.Name
}

func (field *FormField) GetVal() interface{} {
	return field.val
}

func (field *FormField) GetSigned() int64 {
	res, _ := generic.ConvertToInt64(field.val)
	return res
}

func (field *FormField) GetUnsigned() uint64 {
	res, _ := generic.ConvertToUint64(field.val)
	return res
}

func (field *FormField) GetDecimal() float64 {
	res, _ := generic.ConvertToFloat64(field.val)
	return res
}

func (field *FormField) GetString() string {
	res, _ := generic.ConvertToString(field.val)
	return res
}

func (field *FormField) GetBoolean() bool {
	res, _ := generic.ConvertToBoolean(field.val)
	return res
}

// Set internal value of the field
func (field *FormField) SetVal(ctx echo.Context, value interface{}) {
	field.val = value
	if field.Codec == nil {
		field.value, _ = generic.ConvertToString(value)
	} else {
		field.value, _ = field.Codec.Decode(ctx, value)
	}
}

// Get external value of the field
func (field *FormField) GetValue() string {
	return field.value
}

// Check external value of the field
func (field *FormField) SetValue(ctx echo.Context, value string) error {
	if field.Filter != nil {
		value = field.Filter(value)
	}

	field.value = value
	if field.Codec == nil {
		field.val = value
		return nil
	}

	val, err := field.Codec.Encode(ctx, value)
	if err == nil {
		field.val = val
	} else {
		switch e := err.(type) {
		case echo.ValidationError:
			field.AddError(e)
		case echo.ValidationErrors:
			if len(e) == 0 {
				field.val = val
			} else {
				for _, ee := range e {
					field.AddError(ee)
				}
			}
		default:
			return err
		}
	}

	return err
}

// Get disabled flag
func (field *FormField) GetDisabled() bool {
	return field.Disabled
}

// Get hidden flag
func (field *FormField) GetHidden() bool {
	return field.Hidden
}

// Append new error
func (field *FormField) AddError(message echo.ValidationError) {
	field.errors = append(field.errors, message)
}

// Get list of field errors
func (field *FormField) GetErrors() echo.ValidationErrors {
	return field.errors
}

// Test for errors
func (field *FormField) HasErrors() bool {
	return len(field.errors) != 0
}

// Reset field to initial state
func (field *FormField) Reset(ctx echo.Context) error {
	field.errors = nil
	field.value = ""
	field.val = nil
	return nil
}

func (field *FormField) render(
	ctx echo.Context,
) (map[string]interface{}, error) {
	res := make(map[string]interface{}, 16)
	if field.Id != "" {
		res["Id"] = field.Id
	}
	if field.Name != "" {
		res["Name"] = field.Name
	}

	if field.Label != nil {
		label, err := RenderWidget(ctx, field.Label)
		if err != nil {
			return nil, err
		}
		if label != nil {
			res["Label"] = label
		}
	}

	if field.Disabled {
		res["Disabled"] = true
	}

	errors, err := RenderValidationErrors(ctx, field.GetErrors())
	if err != nil {
		return nil, err
	}
	if errors != nil {
		res["Errors"] = errors
	}

	return res, nil
}

func (field *FormField) validateRequired(value string, required bool) bool {
	if required && value == "" {
		field.AddError(MessageConstraintRequired)
		return false
	}
	return true
}

func (field *FormField) validatePattern(value, pattern string, required bool) bool {
	if (required || value != "") && pattern != "" {
		matched, _ := regexp.MatchString(pattern, value)
		if !matched {
			field.AddError(MessageConstraintPattern)
			return false
		}
	}
	return true
}

func (field *FormField) validateMaxLength(value string, maxLength int) bool {
	if maxLength != 0 && utf8.RuneCountInString(value) > maxLength {
		field.AddError(
			&echo.Cause{
				Msg:  uint32(MessageConstraintMaxLength),
				Args: []interface{}{maxLength},
			},
		)
		return false
	}
	return true
}

// FormTextInput represent html entities: <input type="text"> or <input type="password">.
type FormTextInput struct {
	FormField
	Required    bool        // Field is required
	Pattern     string      // Field pattern
	Placeholder interface{} // Field placeholder
	MaxLength   int         // Max length of value
}

func (w *FormTextInput) Render(
	ctx echo.Context,
) (interface{}, error) {
	if w.Hidden {
		return nil, nil
	}

	res, err := w.FormField.render(ctx)
	if err != nil {
		return nil, err
	}

	if w.Required {
		res["Required"] = true
	}

	if w.Pattern != "" {
		res["Pattern"] = w.Pattern
	}

	if w.Placeholder != nil {
		placeholder, err := RenderWidget(ctx, w.Placeholder)
		if err != nil {
			return nil, err
		}
		res["Placeholder"] = placeholder
	}

	if w.MaxLength != 0 {
		res["MaxLen"] = w.MaxLength
	}

	value := w.GetValue()
	if value != "" {
		res["Value"] = value
	}

	return res, nil
}

func (w *FormTextInput) SetValue(
	ctx echo.Context,
	value string,
) error {
	err := w.FormField.SetValue(ctx, value)
	if err != nil {
		return err
	}
	w.FormField.validateRequired(value, w.Required)
	w.FormField.validatePattern(value, w.Pattern, w.Required)
	w.FormField.validateMaxLength(value, w.MaxLength)
	return nil
}

// FormTextArea represent html entity <textarea>.
type FormTextArea struct {
	FormField
	Required    bool        // Field is required
	Pattern     string      // Field pattern
	Placeholder interface{} // Field placeholder
	MaxLength   int         // Field max length
	ReadOnly    bool        // Field is read only
	Rows        int         // Max count of visible rows
}

func (w *FormTextArea) Render(
	ctx echo.Context,
) (interface{}, error) {
	if w.Hidden {
		return nil, nil
	}

	res, err := w.FormField.render(ctx)
	if err != nil {
		return nil, err
	}

	if w.Required {
		res["Required"] = true
	}

	if w.ReadOnly {
		res["Readonly"] = true
	}

	if w.Pattern != "" {
		res["Pattern"] = w.Pattern
	}

	if w.Placeholder != nil {
		placeholder, err := RenderWidget(ctx, w.Placeholder)
		if err != nil {
			return nil, err
		}
		res["Placeholder"] = placeholder
	}

	if w.MaxLength != 0 {
		res["MaxLen"] = w.MaxLength
	}

	if w.Rows != 0 {
		res["Rows"] = w.Rows
	}

	value := w.GetValue()
	if value != "" {
		res["Value"] = value
	}

	return res, nil
}

func (w *FormTextArea) SetValue(
	ctx echo.Context,
	value string,
) error {
	err := w.FormField.SetValue(ctx, value)
	if err != nil {
		return err
	}
	w.FormField.validateRequired(value, w.Required)
	w.FormField.validatePattern(value, w.Pattern, w.Required)
	w.FormField.validateMaxLength(value, w.MaxLength)
	return nil
}

// FormHidden represent html entity <input type="hidden">.
type FormHidden struct {
	FormField
	Required  bool   // Field is required
	Pattern   string // Field pattern
	MaxLength int    // Field max length
}

func (w *FormHidden) Render(
	ctx echo.Context,
) (interface{}, error) {
	if w.Hidden {
		return nil, nil
	}

	res, err := w.FormField.render(ctx)
	if err != nil {
		return nil, err
	}
	res["Value"] = w.GetValue()

	return res, nil
}

func (w *FormHidden) SetValue(
	ctx echo.Context,
	value string,
) error {
	err := w.FormField.SetValue(ctx, value)
	if err != nil {
		return err
	}
	w.FormField.validateRequired(value, w.Required)
	w.FormField.validatePattern(value, w.Pattern, w.Required)
	w.FormField.validateMaxLength(value, w.MaxLength)
	return nil
}

// FormFielInput represent html entity <input type="file">.
type FormFileInput struct {
	FormField
	Accept string // Accept filter
}

func (w *FormFileInput) Render(
	ctx echo.Context,
) (interface{}, error) {
	if w.Hidden {
		return nil, nil
	}

	res, err := w.FormField.render(ctx)
	if err != nil {
		return nil, err
	}

	if w.Accept != "" {
		res["Accept"] = w.Accept
	}

	return res, nil
}

func (w *FormFileInput) Upload(
	ctx echo.Context,
	path string, // file path for store
	name string, // file name for store
) (fileName string, err error) {
	r := ctx.Request()
	r.ParseMultipartForm(32 << 20)

	file, handler, err := r.FormFile(w.Name)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if name == "" {
		name = handler.Filename
	}

	fileName = path + "/" + name
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

// FormSelector represent html entity <select> or other same widget.
// Notice: Field Codec must be ignored (internal and external representations are same).
type FormSelector struct {
	FormField
	Required bool     // Value is required
	Items    data.Set // Field items
}

func (w *FormSelector) Render(
	ctx echo.Context,
) (interface{}, error) {
	if w.Hidden {
		return nil, nil
	}

	res, err := w.FormField.render(ctx)
	if err != nil {
		return nil, err
	}

	if w.Required {
		res["Required"] = true
	}

	options, err := w.renderItems(ctx, w.FormField.GetString())
	if err != nil {
		return nil, err
	}

	if !w.Required {
		label, err := ctx.Echo().Locale.Message(ctx, uint32(MessageSelectorEmpty))
		if err != nil {
			return nil, err
		}
		options = append(
			[]interface{}{
				map[string]interface{}{
					"Value": "",
					"Label": label,
				},
			},
			options...,
		)
	}

	res["Items"] = options

	return res, nil
}

func (w *FormSelector) SetValue(
	ctx echo.Context,
	value string,
) error {
	aVal, aValue := w.val, w.value
	defer func() {
		if w.HasErrors() {
			w.val, w.value = aVal, aValue
		}
	}()

	err := w.FormField.SetValue(ctx, value)
	if err != nil {
		return err
	}

	if w.FormField.validateRequired(value, w.Required) {
		if !w.Items.Has(w.value) {
			w.AddError(echo.ValidationErrorInvalidValue)
		}
	}

	return nil
}

func (w *FormSelector) renderItems(
	ctx echo.Context,
	selected string,
) ([]interface{}, error) {
	rows := make([]interface{}, 0, w.Items.Length())
	err := w.Items.Enumerate(
		ctx,
		func(key, value string) error {
			row := make(map[string]interface{}, 4)
			row["Value"] = key
			row["Label"] = value
			if key == selected {
				row["Selected"] = true
			}
			rows = append(rows, row)
			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

// FormCheckBox represent html entity <input type="checkbox">.
// Example:
//   subscribe := &widget.FormCheckBox{
//       Field: echo.Field{
//           Name: "Subscribe",
//           Codec: widget.CheckBoxCodec,
//       },
//   }
type FormCheckBox struct {
	FormField
	Placeholder interface{} // Placeholder text
}

func (w *FormCheckBox) Render(
	ctx echo.Context,
) (interface{}, error) {
	if w.Hidden {
		return nil, nil
	}

	res, err := w.FormField.render(ctx)
	if err != nil {
		return nil, err
	}

	value := w.GetValue()
	if value != "" && value != "0" {
		res["Checked"] = true
	}

	if w.Placeholder != nil {
		placeholder, err := RenderWidget(ctx, w.Placeholder)
		if err != nil {
			return nil, err
		}
		res["Placeholder"] = placeholder
	}

	return res, nil
}

func (w *FormCheckBox) SetValue(
	ctx echo.Context,
	value string,
) error {
	if value == "" {
		value = "off"
	}
	err := w.FormField.SetValue(ctx, value)
	if err != nil {
		return err
	}
	return nil
}

// FormSubmit represents action Submit
type FormSubmit struct {
	FormField
	Required bool     // Value is required
	Items    data.Set // Optional set of values
}

func (w *FormSubmit) Render(
	ctx echo.Context,
) (interface{}, error) {
	if w.Hidden {
		return nil, nil
	}

	if w.Items == nil {
		// Simple version
		res, err := w.FormField.render(ctx)
		if err != nil {
			return nil, err
		}

		if w.value != "" {
			res["Value"] = w.value
		}

		return res, nil
	}

	// Complex version
	res := make(map[string]interface{}, w.Items.Length())
	err := w.Items.Enumerate(
		ctx,
		func(key, value string) error {
			btn, err := w.FormField.render(ctx)
			if err != nil {
				return err
			}
			btn["Value"] = key
			btn["Label"] = value
			res[key] = btn
			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (w *FormSubmit) SetValue(
	ctx echo.Context,
	value string,
) error {
	aVal, aValue := w.val, w.value
	defer func() {
		if w.HasErrors() {
			w.val, w.value = aVal, aValue
		}
	}()

	keeper := w.GetValue()

	err := w.FormField.SetValue(ctx, value)
	if err != nil {
		return err
	}

	if w.FormField.validateRequired(value, w.Required) {
		if w.Items != nil {
			if !w.Items.Has(w.value) {
				w.AddError(echo.ValidationErrorInvalidValue)
			}
		} else if w.val != keeper {
			w.AddError(echo.ValidationErrorInvalidValue)
		}
	}

	return nil
}

// checkBoxCodec is auxiliary helper for handle html checkbox data.
type checkBoxCodec struct{}

func (codec *checkBoxCodec) Encode(
	ctx echo.Context,
	value string,
) (interface{}, error) {
	value = strings.ToLower(value)
	val := value == "1" || value == "true" || value == "yes" || value == "on"
	return val, nil
}

func (codec *checkBoxCodec) Decode(
	ctx echo.Context,
	value interface{},
) (string, error) {
	val, ok := generic.ConvertToBoolean(value)
	if !ok {
		val = false
	}
	if val {
		return "on", nil
	}
	return "off", nil
}

func (codec *checkBoxCodec) IsEmpty(value interface{}) bool {
	return false
}

var (
	CheckBoxCodec echo.Codec = &checkBoxCodec{}
)
