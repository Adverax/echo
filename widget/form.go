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
		model, err := w.Models.Render(ctx)
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
		model, err := w.Model.Render(ctx)
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
	val    interface{}           // Internal representation of value
	value  []string              // External representation of value
	errors echo.ValidationErrors // Field errors
}

func (field *field) GetInt() int {
	res, _ := generic.ConvertToInt(field.val)
	return res
}

func (field *field) GetInt8() int8 {
	res, _ := generic.ConvertToInt8(field.val)
	return res
}

func (field *field) GetInt16() int16 {
	res, _ := generic.ConvertToInt16(field.val)
	return res
}

func (field *field) GetInt32() int32 {
	res, _ := generic.ConvertToInt32(field.val)
	return res
}

func (field *field) GetInt64() int64 {
	res, _ := generic.ConvertToInt64(field.val)
	return res
}

func (field *field) GetUint() uint {
	res, _ := generic.ConvertToUint(field.val)
	return res
}

func (field *field) GetUint8() uint8 {
	res, _ := generic.ConvertToUint8(field.val)
	return res
}

func (field *field) GetUint16() uint16 {
	res, _ := generic.ConvertToUint16(field.val)
	return res
}

func (field *field) GetUint32() uint32 {
	res, _ := generic.ConvertToUint32(field.val)
	return res
}

func (field *field) GetUint64() uint64 {
	res, _ := generic.ConvertToUint64(field.val)
	return res
}

func (field *field) GetFloat32() float32 {
	res, _ := generic.ConvertToFloat32(field.val)
	return res
}

func (field *field) GetFloat64() float64 {
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
func (field *field) SetVal(
	ctx echo.Context,
	value interface{},
) {
	field.val = value
	val, _ := generic.ConvertToString(value)
	field.value = []string{val}
}

// Set external value of the field.
func (field *field) SetValue(
	ctx echo.Context,
	value []string,
) error {
	aValue := simpleValue(value)
	field.value = []string{aValue}
	field.val = aValue
	return nil
}

// Validate is prototype for all descendants.
func (field *field) Validate(
	ctx echo.Context,
) error {
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
		label, err := echo.RenderWidget(ctx, label)
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
	w.val = value
	var val string
	if w.Codec == nil {
		val, _ = generic.ConvertToString(value)
	} else {
		if w.Required || !generic.IsEmpty(value) {
			val, _ = w.Codec.Decode(ctx, value)
		}
	}
	w.value = []string{val}
}

func (w *FormText) SetValue(
	ctx echo.Context,
	value []string,
) error {
	if w.ReadOnly {
		return nil
	}

	value = filterValue(w.Filter, value)
	aValue := simpleValue(value)

	w.value = []string{aValue}
	if w.Codec == nil {
		w.val = aValue
		return nil
	}

	var val interface{}
	var err error
	if !w.Required && aValue == "" {
		val, _ = w.Codec.Empty(ctx)
	} else {
		val, err = w.Codec.Encode(ctx, aValue)
	}
	if err == nil {
		w.val = val
		return nil
	}

	switch e := err.(type) {
	case echo.ValidationError:
		w.AddError(e)
	case echo.ValidationErrors:
		if len(e) == 0 {
			w.val = val
		} else {
			for _, ee := range e {
				w.AddError(ee)
			}
		}
	default:
		if err != data.ErrNoMatch {
			return err
		}

		w.AddError(echo.ValidationErrorInvalidValue)
	}

	return nil
}

func (w *FormText) Validate(
	ctx echo.Context,
) error {
	if w.ReadOnly {
		return nil
	}

	value := simpleValue(w.value)
	w.validateRequired(value, w.Required)
	w.validatePattern(value, w.Pattern, w.Required)
	w.validateMaxLength(value, w.MaxLength)
	return nil
}

func (w *FormText) Reset(ctx echo.Context) error {
	w.field.reset()
	w.init(ctx)
	return nil
}

func (w *FormText) init(ctx echo.Context) {
	if w.val == nil && w.Default != nil {
		w.SetVal(ctx, w.Default)
	}
}

func (w *FormText) Render(
	ctx echo.Context,
) (interface{}, error) {
	if w.Hidden {
		return nil, nil
	}

	w.init(ctx)

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
		placeholder, err := echo.RenderWidget(ctx, w.Placeholder)
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

func (w *FormSelect) SetValue(
	ctx echo.Context,
	value []string,
) error {
	value = filterValue(w.Filter, value)
	return w.field.SetValue(ctx, value)
}

func (w *FormSelect) Validate(
	ctx echo.Context,
) error {
	value := simpleValue(w.value)
	if w.validateRequired(value, w.Required) && value != "" {
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

	w.init(ctx)

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
		empty := make(map[string]interface{}, 2)
		label, err := ctx.Echo().Locale.Message(ctx, uint32(MessageSelectorEmpty))
		if err != nil {
			return nil, err
		}
		empty["Label"] = label
		value := w.GetValue()
		selected := len(value) != 0 && (value[0] == "" || value[0] == "0")
		if selected {
			empty["Selected"] = true
		}
		res["Empty"] = empty
	}

	return res, nil
}

func (w *FormSelect) Reset(ctx echo.Context) error {
	w.field.reset()
	w.init(ctx)
	return nil
}

func (w *FormSelect) init(ctx echo.Context) {
	if w.val == nil && w.Default != nil {
		w.SetVal(ctx, w.Default)
	}
}

// FormFlag represents single html entity <input tpe="ckeckbox>
// Example:
//   subscribe := widget.FormFlag{
//     Name: "Subscribe",
//     Items: ....,
//   }
type FormFlag struct {
	field
	Id          string              // Field identifier
	Name        string              // Field name
	Label       interface{}         // Field label
	Disabled    bool                // Field disabled
	Hidden      bool                // Field is hidden (not rendered)
	Default     interface{}         // Default value
	Placeholder interface{}         // Placeholder text
	Filter      FormFieldFilterFunc // Custom filter
}

func (w *FormFlag) GetName() string {
	return w.Name
}

func (w *FormFlag) GetDisabled() bool {
	return w.Disabled
}

func (w *FormFlag) GetHidden() bool {
	return w.Hidden
}

func (w *FormFlag) SetVal(ctx echo.Context, value interface{}) {
	val, _ := generic.ConvertToBoolean(value)
	w.field.SetVal(ctx, val)
}

func (w *FormFlag) SetValue(
	ctx echo.Context,
	value []string,
) error {
	value = filterValue(w.Filter, value)

	w.value = value
	w.val = len(value) != 0 && value[0] == "1"
	return nil
}

func (w *FormFlag) Render(
	ctx echo.Context,
) (interface{}, error) {
	if w.Hidden {
		return nil, nil
	}

	w.init(ctx)

	res, err := w.render(ctx, w.Id, w.Name, w.Label, w.Disabled)
	if err != nil {
		return nil, err
	}

	val, _ := generic.ConvertToBoolean(w.val)
	if val {
		res["Selected"] = true
	}

	res["Value"] = "1"

	if w.Placeholder != nil {
		placeholder, err := echo.RenderWidget(ctx, w.Placeholder)
		if err != nil {
			return nil, err
		}
		res["Placeholder"] = placeholder
	}

	return res, nil
}

func (w *FormFlag) Reset(ctx echo.Context) error {
	w.field.reset()
	w.init(ctx)
	return nil
}

func (w *FormFlag) init(ctx echo.Context) {
	if w.val == nil && w.Default != nil {
		w.SetVal(ctx, w.Default)
	}
}

// FormFlags represent html entity <input type="checkbox">.
// Example:
//   subscribe := &widget.FormFlags{
//       Name: "Subscribe",
//       Items: ...,
//   }
type FormFlags struct {
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

func (w *FormFlags) GetName() string {
	return w.Name
}

func (w *FormFlags) GetDisabled() bool {
	return w.Disabled
}

func (w *FormFlags) GetHidden() bool {
	return w.Hidden
}

func (w *FormFlags) SetVal(ctx echo.Context, value interface{}) {
	switch v := value.(type) {
	case []string:
		w.value = v
	default:
		vv, _ := generic.ConvertToString(v)
		w.value = []string{vv}
	}
	w.val = w.value
}

func (w *FormFlags) SetValue(
	ctx echo.Context,
	value []string,
) error {
	value = filterValue(w.Filter, value)

	w.value = value

	keys, err := echo.DataSetKeys(ctx, w.Items)
	if err != nil {
		return err
	}

	w.val = intersect(value, keys)
	return nil
}

func (w *FormFlags) Render(
	ctx echo.Context,
) (interface{}, error) {
	if w.Hidden {
		return nil, nil
	}

	w.init(ctx)

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
	res["Items"] = items

	if w.Placeholder != nil {
		placeholder, err := echo.RenderWidget(ctx, w.Placeholder)
		if err != nil {
			return nil, err
		}
		res["Placeholder"] = placeholder
	}

	return res, nil
}

func (w *FormFlags) Reset(ctx echo.Context) error {
	w.field.reset()
	w.init(ctx)
	return nil
}

func (w *FormFlags) init(ctx echo.Context) {
	if w.val == nil && w.Default != nil {
		w.SetVal(ctx, w.Default)
	}
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

func (w *FormSubmit) SetValue(
	ctx echo.Context,
	value []string,
) error {
	value = filterValue(w.Filter, value)
	return w.field.SetValue(ctx, value)
}

func (w *FormSubmit) Validate(
	ctx echo.Context,
) error {
	v := simpleValue(w.value)
	if w.validateRequired(v, w.Required) {
		if w.Items == nil {
			if w.Default == nil {
				return nil
			}
			a, _ := generic.ConvertToString(w.Default)
			b, _ := generic.ConvertToString(v)
			if a != b {
				w.AddError(echo.ValidationErrorInvalidValue)
			}
			return nil
		}

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

	w.init(ctx)

	res, err := w.render(ctx, w.Id, w.Name, w.Label, w.Disabled)
	if err != nil {
		return nil, err
	}

	if w.Required {
		res["Required"] = true
	}

	if w.Items == nil {
		if val := w.GetString(); val != "" {
			res["Value"] = val
		}
		return res, nil
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

func (w *FormSubmit) init(ctx echo.Context) {
	if w.val == nil && w.Default != nil {
		w.SetVal(ctx, w.Default)
	}
}

func (w *FormSubmit) Reset(ctx echo.Context) error {
	w.field.reset()
	w.init(ctx)
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

func (w *FormHidden) Validate(
	ctx echo.Context,
) error {
	value := simpleValue(w.value)
	w.validateRequired(value, w.Required)
	w.validatePattern(value, w.Pattern, w.Required)
	w.validateMaxLength(value, w.MaxLength)
	return nil
}

func (w *FormHidden) Render(
	ctx echo.Context,
) (interface{}, error) {
	if w.Hidden {
		return nil, nil
	}

	w.init(ctx)

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
	w.init(ctx)
	return nil
}

func (w *FormHidden) init(ctx echo.Context) {
	if w.val == nil && w.Default != nil {
		w.SetVal(ctx, w.Default)
	}
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

func (w *FormFile) Validate(
	ctx echo.Context,
) error {
	value := simpleValue(w.value)
	w.validateRequired(value, w.Required)
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
	if filter == nil {
		return value
	}

	res := make([]string, len(value))
	for i, val := range value {
		res[i] = filter(val)
	}
	return res
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
