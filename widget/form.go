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
	"io"
	"os"
	"regexp"
	"sort"
	"unicode/utf8"

	"github.com/adverax/echo"
	"github.com/adverax/echo/data"
	"github.com/adverax/echo/generic"
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

type field struct {
	val         interface{}           // Internal representation of value
	value       []string              // External representation of value
	errors      echo.ValidationErrors // Field errors
	initialized bool                  // Field is initialized
}

func (field *field) GetSigned() int64 {
	res, _ := generic.ConvertToInt64(field.val)
	return res
}

func (field *field) GetUnsigned() uint64 {
	res, _ := generic.ConvertToUint64(field.val)
	return res
}

func (field *field) GetDecimal() float64 {
	res, _ := generic.ConvertToFloat64(field.val)
	return res
}

func (field *field) GetString() string {
	res, _ := generic.ConvertToString(field.val)
	return res
}

func (field *field) GetBoolean() bool {
	res, _ := generic.ConvertToBoolean(field.val)
	return res
}

func (field *field) GetVal() interface{} {
	return field.val
}

func (field *field) GetValue() []string {
	return field.value
}

// Set internal value of the field
func (field *field) setVal(
	ctx echo.Context,
	value interface{},
	codec echo.Codec,
) {
	field.val = value
	field.initialized = true
	var val string
	if codec == nil {
		val, _ = generic.ConvertToString(value)
	} else {
		val, _ = codec.Decode(ctx, value)
	}
	field.value = []string{val}
}

// Set external value of the field.
// This is base version, that can be override by descendants.
func (field *field) setValue(
	ctx echo.Context,
	value []string,
	codec echo.Codec,
) error {
	aValue := simpleValue(value)

	field.initialized = true
	field.value = []string{aValue}
	if codec == nil {
		field.val = aValue
		return nil
	}

	val, err := codec.Encode(ctx, aValue)
	if err == nil {
		field.val = val
		return nil
	}

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
		if err != data.ErrNoMatch {
			return err
		}

		field.AddError(echo.ValidationErrorInvalidValue)
	}

	return nil
}

func (field *field) render(
	ctx echo.Context,
	id string,
	name string,
	label interface{},
	disabled bool,
) (map[string]interface{}, error) {
	res := make(map[string]interface{}, 16)
	if id != "" {
		res["Id"] = id
	}
	if name != "" {
		res["Name"] = name
	}

	if label != nil {
		label, err := RenderWidget(ctx, label)
		if err != nil {
			return nil, err
		}
		if label != nil {
			res["Label"] = label
		}
	}

	if disabled {
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

// Append a new error
func (field *field) AddError(message echo.ValidationError) {
	field.errors = append(field.errors, message)
}

// Get list of field errors
func (field *field) GetErrors() echo.ValidationErrors {
	return field.errors
}

// Test for errors
func (field *field) IsValid() bool {
	return len(field.errors) == 0
}

// Reset field to initial state
func (field *field) reset() {
	field.errors = nil
	field.value = nil
	field.val = nil
	field.initialized = false
}

func (field *field) validateRequired(value string, required bool) bool {
	if required && value == "" {
		field.AddError(MessageConstraintRequired)
		return false
	}
	return true
}

func (field *field) validatePattern(value, pattern string, required bool) bool {
	if (required || value != "") && pattern != "" {
		matched, _ := regexp.MatchString(pattern, value)
		if !matched {
			field.AddError(MessageConstraintPattern)
			return false
		}
	}
	return true
}

func (field *field) validateMaxLength(value string, maxLength int) bool {
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

// FormText represent html entity <textarea> or <input type="text"> or <input type="password">.
type FormText struct {
	field
	Id          string              // Field identifier
	Name        string              // Field name
	Label       interface{}         // Field label
	Disabled    bool                // Field disabled
	Hidden      bool                // Field is hidden (not rendered)
	Filter      FormFieldFilterFunc // Custom filter
	Codec       echo.Codec          // Field codec (optional)
	Default     interface{}         // Default value
	Required    bool                // Field is required
	Pattern     string              // Field pattern
	Placeholder interface{}         // Field placeholder
	MaxLength   int                 // Field max length
	ReadOnly    bool                // Field is read only
	Rows        int                 // Max count of visible rows
}

func (w *FormText) GetName() string {
	return w.Name
}

func (w *FormText) GetDisabled() bool {
	return w.Disabled
}

func (w *FormText) GetHidden() bool {
	return w.Hidden
}

func (w *FormText) SetVal(ctx echo.Context, value interface{}) {
	w.field.setVal(ctx, value, w.Codec)
}

func (w *FormText) SetValue(
	ctx echo.Context,
	value []string,
) error {
	value = filterValue(w.Filter, value)
	err := w.setValue(ctx, value, w.Codec)
	if err != nil {
		return err
	}

	aValue := simpleValue(value)
	w.validateRequired(aValue, w.Required)
	w.validatePattern(aValue, w.Pattern, w.Required)
	w.validateMaxLength(aValue, w.MaxLength)
	return nil
}

func (w *FormText) Reset(ctx echo.Context) error {
	w.field.reset()
	if w.Default != nil {
		w.SetVal(ctx, w.Default)
	}
	return nil
}

func (w *FormText) Render(
	ctx echo.Context,
) (interface{}, error) {
	if w.Hidden {
		return nil, nil
	}

	res, err := w.field.render(ctx, w.Id, w.Name, w.Label, w.Disabled)
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
	if len(value) != 0 {
		res["Value"] = value[0]
	}

	return res, nil
}

// FormSelect represent html entity <select> or other same widget.
// Notice: Field Codec must be ignored (internal and external representations are same).
type FormSelect struct {
	field
	Id       string              // Field identifier
	Name     string              // Field name
	Label    interface{}         // Field label
	Disabled bool                // Field disabled
	Hidden   bool                // Field is hidden (not rendered)
	Filter   FormFieldFilterFunc // Custom filter
	Default  interface{}         // Default value
	Required bool                // Value is required
	Items    echo.DataSet        // Field items
}

func (w *FormSelect) GetName() string {
	return w.Name
}

func (w *FormSelect) GetDisabled() bool {
	return w.Disabled
}

func (w *FormSelect) GetHidden() bool {
	return w.Hidden
}

func (w *FormSelect) SetVal(ctx echo.Context, value interface{}) {
	w.field.setVal(ctx, value, nil)
}

func (w *FormSelect) SetValue(
	ctx echo.Context,
	value []string,
) error {
	aVal, aValue := w.val, w.value
	defer func() {
		if !w.IsValid() {
			w.val, w.value = aVal, aValue
		}
	}()

	value = filterValue(w.Filter, value)
	err := w.setValue(ctx, value, nil)
	if err != nil {
		return err
	}

	v := simpleValue(value)
	if w.validateRequired(v, w.Required) {
		_, err := w.Items.Decode(ctx, w.val)
		if err != nil {
			if err != data.ErrNoMatch {
				return err
			}
			w.AddError(echo.ValidationErrorInvalidValue)
		}
	}

	return nil
}

func (w *FormSelect) Render(
	ctx echo.Context,
) (interface{}, error) {
	if w.Hidden {
		return nil, nil
	}

	res, err := w.render(ctx, w.Id, w.Name, w.Label, w.Disabled)
	if err != nil {
		return nil, err
	}

	if w.Required {
		res["Required"] = true
	}

	selected := map[string]bool{
		w.GetString(): true,
	}
	items, err := RenderDataSet(ctx, w.Items, selected)
	if err != nil {
		return nil, err
	}
	if items != nil {
		res["Items"] = items
	}

	if !w.Required {
		label, err := ctx.Echo().Locale.Message(ctx, uint32(MessageSelectorEmpty))
		if err != nil {
			return nil, err
		}
		res["Empty"] = label
	}

	return res, nil
}

func (w *FormSelect) Reset(ctx echo.Context) error {
	w.field.reset()
	if w.Default != nil {
		w.SetVal(ctx, w.Default)
	}
	return nil
}

// FormMultiSelect represent html entity <input type="checkbox">.
// Example:
//   subscribe := &widget.FormMultiSelect{
//       Name: "Subscribe",
//       Items: ...,
//   }
type FormMultiSelect struct {
	field
	Id          string              // Field identifier
	Name        string              // Field name
	Label       interface{}         // Field label
	Disabled    bool                // Field disabled
	Hidden      bool                // Field is hidden (not rendered)
	Filter      FormFieldFilterFunc // Custom filter
	Default     interface{}         // Default value
	Items       echo.DataSet        // List labels for allowed values
	Placeholder interface{}         // Placeholder text
}

func (w *FormMultiSelect) GetName() string {
	return w.Name
}

func (w *FormMultiSelect) GetDisabled() bool {
	return w.Disabled
}

func (w *FormMultiSelect) GetHidden() bool {
	return w.Hidden
}

func (w *FormMultiSelect) SetVal(ctx echo.Context, value interface{}) {
	w.initialized = true
	if v, ok := value.([]string); ok {
		w.value = v
	} else {
		v, _ := generic.ConvertToString(value)
		w.value = []string{v}
	}
	w.val = w.value
}

func (w *FormMultiSelect) SetValue(
	ctx echo.Context,
	value []string,
) error {
	value = filterValue(w.Filter, value)

	w.initialized = true
	w.value = value

	keys, err := echo.DataSetKeys(ctx, w.getItems())
	if err != nil {
		return err
	}

	w.val = intersect(value, keys)
	return nil
}

func (w *FormMultiSelect) Render(
	ctx echo.Context,
) (interface{}, error) {
	if w.Hidden {
		return nil, nil
	}

	res, err := w.render(ctx, w.Id, w.Name, w.Label, w.Disabled)
	if err != nil {
		return nil, err
	}

	vs := w.GetValue()
	selected := make(map[string]bool, len(vs))
	for _, v := range vs {
		selected[v] = true
	}
	items, err := RenderDataSet(ctx, w.Items, selected)
	if err != nil {
		return nil, err
	}
	if items != nil {
		res["Items"] = items
	} else {
		value := w.value
		if len(value) != 0 {
			res["Value"] = value[0]
		}
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

func (w *FormMultiSelect) Reset(ctx echo.Context) error {
	w.field.reset()
	if w.Default != nil {
		w.SetVal(ctx, w.Default)
	}
	return nil
}

func (w *FormMultiSelect) getItems() echo.DataSet {
	if w.Items != nil {
		return w.Items
	}

	return defaultMultiSelectItems
}

// FormSubmit represents action Submit
type FormSubmit struct {
	field
	Id       string              // Field identifier
	Name     string              // Field name
	Label    interface{}         // Field label
	Disabled bool                // Field disabled
	Hidden   bool                // Field is hidden (not rendered)
	Filter   FormFieldFilterFunc // Custom filter
	Default  interface{}         // Default value
	Required bool                // Value is required
	Items    echo.DataSet        // Field items
}

func (w *FormSubmit) GetName() string {
	return w.Name
}

func (w *FormSubmit) GetDisabled() bool {
	return w.Disabled
}

func (w *FormSubmit) GetHidden() bool {
	return w.Hidden
}

func (w *FormSubmit) SetVal(ctx echo.Context, value interface{}) {
	w.field.setVal(ctx, value, nil)
}

func (w *FormSubmit) SetValue(
	ctx echo.Context,
	value []string,
) error {
	aVal, aValue := w.val, w.value
	defer func() {
		if !w.IsValid() {
			w.val, w.value = aVal, aValue
		}
	}()

	value = filterValue(w.Filter, value)
	err := w.setValue(ctx, value, nil)
	if err != nil {
		return err
	}

	v := simpleValue(value)
	if w.validateRequired(v, w.Required) {
		_, err := w.Items.Decode(ctx, w.val)
		if err != nil {
			if err != data.ErrNoMatch {
				return err
			}
			w.AddError(echo.ValidationErrorInvalidValue)
		}
	}

	return nil
}

func (w *FormSubmit) Render(
	ctx echo.Context,
) (interface{}, error) {
	if w.Hidden {
		return nil, nil
	}

	res, err := w.render(ctx, w.Id, w.Name, w.Label, w.Disabled)
	if err != nil {
		return nil, err
	}

	if w.Required {
		res["Required"] = true
	}

	Submited := map[string]bool{
		w.GetString(): true,
	}
	items, err := RenderDataSet(ctx, w.Items, Submited)
	if err != nil {
		return nil, err
	}
	if items != nil {
		res["Items"] = items
	}

	return res, nil
}

func (w *FormSubmit) Reset(ctx echo.Context) error {
	w.field.reset()
	if w.Default != nil {
		w.SetVal(ctx, w.Default)
	}
	return nil
}

// FormHidden represent html entity <input type="hidden">.
type FormHidden struct {
	field
	Id        string      // Field identifier
	Name      string      // Field name
	Hidden    bool        // Field is hidden (not rendered)
	Default   interface{} // Default value
	Required  bool        // Field is required
	Pattern   string      // Field pattern
	MaxLength int         // Field max length
}

func (w *FormHidden) GetName() string {
	return w.Name
}

func (w *FormHidden) GetDisabled() bool {
	return false
}

func (w *FormHidden) GetHidden() bool {
	return w.Hidden
}

func (w *FormHidden) SetVal(ctx echo.Context, value interface{}) {
	w.field.setVal(ctx, value, nil)
}

func (w *FormHidden) SetValue(
	ctx echo.Context,
	value []string,
) error {
	err := w.field.setValue(ctx, value, nil)
	if err != nil {
		return err
	}

	aValue := simpleValue(value)
	w.validateRequired(aValue, w.Required)
	w.validatePattern(aValue, w.Pattern, w.Required)
	w.validateMaxLength(aValue, w.MaxLength)
	return nil
}

func (w *FormHidden) Render(
	ctx echo.Context,
) (interface{}, error) {
	if w.Hidden {
		return nil, nil
	}

	res, err := w.render(ctx, w.Id, w.Name, nil, false)
	if err != nil {
		return nil, err
	}

	value := w.GetValue()
	if len(value) != 0 {
		res["Value"] = value[0]
	}

	return res, nil
}

func (w *FormHidden) Reset(ctx echo.Context) error {
	w.field.reset()
	if w.Default != nil {
		w.SetVal(ctx, w.Default)
	}
	return nil
}

// FormFile represent html entity <input type="file">.
type FormFile struct {
	field
	Id       string      // Field identifier
	Name     string      // Field name
	Label    interface{} // Field label
	Disabled bool        // Field disabled
	Hidden   bool        // Field is hidden (not rendered)
	Accept   string      // Accept filter
	Required bool        // Field is required
}

func (w *FormFile) GetName() string {
	return w.Name
}

func (w *FormFile) GetDisabled() bool {
	return w.Disabled
}

func (w *FormFile) GetHidden() bool {
	return w.Hidden
}

func (w *FormFile) SetVal(ctx echo.Context, value interface{}) {
	w.field.setVal(ctx, value, nil)
}

func (w *FormFile) SetValue(
	ctx echo.Context,
	value []string,
) error {
	err := w.field.setValue(ctx, value, nil)
	if err != nil {
		return err
	}

	aValue := simpleValue(value)
	w.validateRequired(aValue, w.Required)
	// todo: validate accept filter
	return nil
}

func (w *FormFile) Render(
	ctx echo.Context,
) (interface{}, error) {
	if w.Hidden {
		return nil, nil
	}

	res, err := w.render(ctx, w.Id, w.Name, w.Label, w.Disabled)
	if err != nil {
		return nil, err
	}

	if w.Accept != "" {
		res["Accept"] = w.Accept
	}

	return res, nil
}

func (w *FormFile) Upload(
	ctx echo.Context,
	path string, // file path for store
	name string, // file name for store
) (fileName string, err error) {
	_, err = ctx.MultipartForm()
	if err != nil {
		return "", err
	}

	file, handler, err := ctx.Request().FormFile(w.Name)
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

func (w *FormFile) Reset(ctx echo.Context) error {
	w.field.reset()
	return nil
}

func simpleValue(value []string) string {
	if len(value) == 0 {
		return ""
	}

	return value[0]
}

func filterValue(filter FormFieldFilterFunc, value []string) []string {
	if filter != nil {
		for i, val := range value {
			value[i] = filter(val)
		}
	}

	return value
}

func intersect(as, bs []string) []string {
	sort.Strings(as)
	sort.Strings(bs)

	la := len(as)
	lb := len(bs)
	var lc int
	if la > lb {
		lc = lb
	} else {
		lc = la
	}
	if lc == 0 {
		return nil
	}

	a := 0
	b := 0
	cs := make([]string, 0, lc)
	for a < la && b < lb {
		if as[a] < bs[b] {
			a++
		} else if as[a] > bs[b] {
			b++
		} else {
			cs = append(cs, as[a])
			a++
			b++
		}
	}

	if len(cs) == 0 {
		return nil
	}

	return cs
}

var defaultMultiSelectItems = echo.NewDataSet(
	map[string]string{
		"on":  "on",
		"off": "off",
	},
	false,
)
