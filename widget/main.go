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
	"fmt"
	"html"
	"html/template"
	"regexp"
	"strconv"
	"time"

	"github.com/adverax/echo"
	"github.com/adverax/echo/generic"
	"net/url"
)

type GuidMaker interface {
	CreateGuid() uint64
}

type DataFunc func() (interface{}, error)

type ExpanderFund func(data map[string]interface{}) error

// Entity, that can be converted to plain text
type Stringer interface {
	String(ctx echo.Context) (string, error)
}

var (
	BoolFormatter     echo.Formatter
	TextFormatter     = &echo.BaseFormatter{Decoder: echo.TextCodec}
	SignedFormatter   = &echo.BaseFormatter{Decoder: echo.SignedCodec}
	UnsignedFormatter = &echo.BaseFormatter{Decoder: echo.UnsignedCodec}
	DecimalFormatter  = &echo.BaseFormatter{Decoder: echo.DecimalCodec}
	DefaultFormatter  = TextFormatter
)

var (
	MessagePagerPrev                = DeclareDefaultMsg(1000, "Prev")
	MessagePagerNext                = DeclareDefaultMsg(1001, "Next")
	MessageListRecords              = DeclareDefaultMsg(1002, "Shows rows from %d to %d of %d")
	MessageListNoRecords            = DeclareDefaultMsg(1003, "No data for display")
	MessageDetailsColumnKey         = DeclareDefaultMsg(1004, "Key")
	MessageDetailsColumnVal         = DeclareDefaultMsg(1005, "Value")
	MessageTableTooltipActionView   = DeclareDefaultMsg(1006, "View")
	MessageTableTooltipActionUpdate = DeclareDefaultMsg(1007, "Update")
	MessageTableTooltipActionDelete = DeclareDefaultMsg(1008, "Delete this row")
	MessageTableConfirmActionDelete = DeclareDefaultMsg(1009, "Are you sure delete row?")
	MessageConstraintInvalid        = DeclareDefaultMsg(1010, "Invalid value")
	MessageConstraintRequired       = DeclareDefaultMsg(1011, "Required value")
	MessageConstraintPattern        = DeclareDefaultMsg(1012, "Value don't match pattern")
	MessageConstraintMaxLength      = DeclareDefaultMsg(1013, "Value too long")
	MessageSelectorEmpty            = DeclareDefaultMsg(1014, "(Empty)")
	MessageMultistepFormPrev        = DeclareDefaultMsg(1015, "Prev")
	MessageMultistepFormNext        = DeclareDefaultMsg(1016, "Next")
)

func DeclareDefaultMsg(msg MESSAGE, message string) MESSAGE {
	echo.DefaultMessages[uint32(msg)] = message
	return msg
}

// Widget is abstract entity, that generated output data
type Widget interface {
	Render(ctx echo.Context) (interface{}, error)
}

type Map map[string]interface{}

func (ws Map) Render(ctx echo.Context) (interface{}, error) {
	return RenderMap(ctx, ws)
}

func (ws Map) Clone() Map {
	res := make(Map, 2*len(ws))
	for key, val := range ws {
		res[key] = val
	}
	return res
}

type List []interface{}

func (ws List) Render(ctx echo.Context) (interface{}, error) {
	return RenderList(ctx, ws)
}

type WidgetFunc func(ctx echo.Context) (interface{}, error)

func (fn WidgetFunc) Render(ctx echo.Context) (interface{}, error) {
	return fn(ctx)
}

// Optional content
type Optional struct {
	Hidden bool
	Value  interface{}
}

func (w *Optional) Render(ctx echo.Context) (interface{}, error) {
	if w.Hidden {
		return nil, nil
	}
	return RenderWidget(ctx, w.Value)
}

// Raw text
type TEXT string

func (w TEXT) Render(ctx echo.Context) (interface{}, error) {
	return TextFormatter.Format(ctx, string(w))
}

func (w TEXT) String(ctx echo.Context) (string, error) {
	return string(w), nil
}

// Boolean value
type BOOLEAN bool

func (w BOOLEAN) Render(ctx echo.Context) (interface{}, error) {
	formatter := BoolFormatter
	if formatter != nil {
		if w {
			return formatter.Format(ctx, 1)
		} else {
			return formatter.Format(ctx, 0)
		}
	}

	if w {
		return "True", nil
	} else {
		return "False", nil
	}
}

// Signed integer value
type SIGNED int64

func (w SIGNED) Render(ctx echo.Context) (interface{}, error) {
	return SignedFormatter.Format(ctx, int64(w))
}

// Unsigned integer value
type UNSIGNED uint64

func (w UNSIGNED) Render(ctx echo.Context) (interface{}, error) {
	return UnsignedFormatter.Format(ctx, uint64(w))
}

// Decimal (float) value
type DECIMAL float64

func (w DECIMAL) Render(ctx echo.Context) (interface{}, error) {
	return DecimalFormatter.Format(ctx, float64(w))
}

// Message
type MESSAGE uint32

func (w MESSAGE) Error() string {
	return "Validation error: "
}

func (w MESSAGE) Render(ctx echo.Context) (interface{}, error) {
	msg, err := w.String(ctx)
	if err != nil {
		return "", err
	}
	return msg, nil
}

func (w MESSAGE) String(ctx echo.Context) (string, error) {
	return ctx.Locale().Message(ctx, uint32(w))
}

func (w MESSAGE) Translate(
	ctx echo.Context,
) (string, error) {
	return ctx.Locale().Message(ctx, uint32(w))
}

// Resource
type RESOURCE uint32

func (w RESOURCE) Render(ctx echo.Context) (interface{}, error) {
	res, err := ctx.Locale().Resource(ctx, uint32(w))
	if err != nil {
		return nil, err
	}
	return template.HTML(res), nil
}

// Localize UTC datetime
type DATETIME string

func (w DATETIME) Render(ctx echo.Context) (interface{}, error) {
	tm, err := time.ParseInLocation(generic.DateTimeFormat, string(w), time.UTC)
	if err != nil {
		return nil, err
	}

	return ctx.Locale().FormatDateTime(tm), nil
}

// Time
type TIMESTAMP int64

func (w TIMESTAMP) Render(ctx echo.Context) (interface{}, error) {
	locale := ctx.Locale()
	tm := time.Unix(int64(w), 0).In(locale.Location())
	return locale.FormatDateTime(tm), nil
}

// Time duration
type DURATION int64

func (w DURATION) Render(ctx echo.Context) (interface{}, error) {
	// Seconds
	sec := w % 60
	w /= 60
	if w == 0 {
		return strconv.FormatInt(int64(sec), 10), nil
	}

	// Minutes
	min := w % 60
	w /= 60
	if w == 0 {
		return fmt.Sprintf("%d:%d", min, sec), nil
	}

	// Hours
	hour := w % 23
	w /= 24
	if w == 0 {
		return fmt.Sprintf("%d:%d:%d", hour, min, sec), nil
	}

	// Days
	return fmt.Sprintf("%d %d:%d:%d", w, hour, min, sec), nil
}

type HTML template.HTML

func (w HTML) Render(ctx echo.Context) (interface{}, error) {
	return template.HTML(w), nil
}

// Format by layout (using fmt.SprintF)
type Format struct {
	Layout interface{}   // Layout
	Args   []interface{} // Message parameters
}

func (w *Format) Render(ctx echo.Context) (interface{}, error) {
	msg, err := w.String(ctx)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (w *Format) String(ctx echo.Context) (string, error) {
	if w.Layout == nil {
		return "", nil
	}
	val, err := RenderWidget(ctx, w.Layout)
	if err != nil {
		return "", err
	}
	if layout, ok := generic.ConvertToString(val); ok {
		if len(w.Args) == 0 {
			return layout, nil
		} else {
			return fmt.Sprintf(layout, w.Args...), nil
		}
	}
	return "", nil
}

// Template with named params
// Layout can used params like this: {{name@param}}
type Template struct {
	Layout interface{}    // Layout
	Params generic.Params // Message arguments
}

func (w *Template) Render(ctx echo.Context) (interface{}, error) {
	if w.Layout == nil {
		return "", nil
	}
	val, err := RenderWidget(ctx, w.Layout)
	if err != nil {
		return nil, err
	}
	if layout, ok := generic.ConvertToString(val); ok {
		return w.renderParams(ctx, layout)
	}
	return "", nil
}

func (w *Template) renderParams(ctx echo.Context, code string) (interface{}, error) {
	if len(w.Params) == 0 {
		return code, nil
	}

	var firstErr error

	code2 := DocumentParamRe.ReplaceAllStringFunc(
		string(code),
		func(s string) string {
			matches := DocumentParamRe.FindStringSubmatch(s)
			if matches != nil {
				name := matches[1]
				value, ok := w.Params[name]
				if ok {
					if ww, valid := value.(Widget); valid {
						c, err := ww.Render(ctx)
						if err == nil {
							v, err := ConvertToHtml(c)
							if err != nil {
								return ""
							}
							return v
						}
						if firstErr == nil {
							firstErr = err
						}
					} else {
						val, success := generic.ConvertToString(value)
						if success {
							return val
						}
					}
				}
			}
			return s
		},
	)

	if firstErr != nil {
		return nil, firstErr
	}

	return template.HTML(code2), nil
}

// Any value with formatter
type Variant struct {
	echo.Formatter
	Value interface{}
}

func (w *Variant) Render(ctx echo.Context) (interface{}, error) {
	f := w.Formatter
	if f == nil {
		f = DefaultFormatter
	}

	val, err := f.Format(ctx, w.Value)
	if err != nil {
		return "", err
	}

	return val, nil
}

var DocumentParamRe = regexp.MustCompile(`(?i:{{\s*([\w\d.\-]+)@param\s*}})`)

func RenderParams(
	ctx echo.Context,
	params []interface{},
) ([]interface{}, error) {
	list := make([]interface{}, len(params))
	for i, param := range params {
		if param == nil {
			list[i] = ""
		} else {
			if widget, ok := param.(Widget); ok {
				item, err := widget.Render(ctx)
				if err != nil {
					return nil, err
				}
				list[i] = item
			} else {
				list[i] = param
			}
		}
	}
	return list, nil
}

func ConvertToHtml(v interface{}) (string, error) {
	switch val := v.(type) {
	case template.HTML:
		return string(val), nil
	default:
		res, _ := generic.ConvertToString(val)
		return html.EscapeString(res), nil
	}
}

func RenderWidget(
	ctx echo.Context,
	v interface{},
) (interface{}, error) {
	if w, ok := v.(Widget); ok {
		return w.Render(ctx)
	}

	return v, nil
}

func RenderMap(
	ctx echo.Context,
	widgets Map,
) (map[string]interface{}, error) {
	res := make(map[string]interface{}, len(widgets))
	for key, widget := range widgets {
		if widget != nil {
			item, err := RenderWidget(ctx, widget)
			if err != nil {
				return nil, err
			}
			if item != nil {
				res[key] = item
			}
		}
	}

	if len(res) == 0 {
		return nil, nil
	}

	return res, nil
}

func RenderList(
	ctx echo.Context,
	list List,
) ([]interface{}, error) {
	var res []interface{}
	for _, widget := range list {
		if widget != nil {
			item, err := RenderWidget(ctx, widget)
			if err != nil {
				return nil, err
			}
			if item != nil {
				res = append(res, item)
			}
		}
	}
	return res, nil
}

func RenderModel(
	ctx echo.Context,
	model echo.Model,
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

func RenderModels(
	ctx echo.Context,
	models echo.Models,
) (interface{}, error) {
	res := make([]interface{}, 0, len(models))
	for _, model := range models {
		if model != nil {
			item, err := RenderModel(ctx, model)
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

// Render all validation errors
func RenderValidationErrors(
	ctx echo.Context,
	errors echo.ValidationErrors,
) ([]string, error) {
	if len(errors) == 0 {
		return nil, nil
	}

	errs := make([]string, len(errors))
	for i, err := range errors {
		msg, err := err.Translate(ctx)
		if err != nil {
			return nil, err
		}
		errs[i] = msg
	}

	return errs, nil
}

// Render link without escape
func RenderLink(ctx echo.Context, v interface{}) (string, error) {
	switch val := v.(type) {
	case string:
		return val, nil
	case *url.URL:
		return ctx.Echo().UrlLinker.Render(ctx, val)
	case url.URL:
		return ctx.Echo().UrlLinker.Render(ctx, &val)
	default:
		return "", nil
	}
}

func RenderDataSet(
	ctx echo.Context,
	dataset echo.DataSet,
	selected map[string]bool,
) ([]interface{}, error) {
	if dataset == nil {
		return nil, nil
	}

	items, err := dataset.DataSet(ctx)
	if err != nil {
		return nil, err
	}

	length, err := items.Length(ctx)
	if err != nil {
		return nil, err
	}

	rows := make([]interface{}, 0, length)
	err = items.Enumerate(
		ctx,
		func(key, value string) error {
			row := make(map[string]interface{}, 4)
			row["Value"] = key
			row["Label"] = value
			if _, has := selected[key]; has {
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

func FormatMessage(
	ctx echo.Context,
	message interface{},
	args ...interface{},
) (string, error) {
	msg, err := RenderWidget(ctx, message)
	if err != nil {
		return "", err
	}

	if format, ok := msg.(string); ok {
		return fmt.Sprintf(format, args...), nil
	}

	return "", nil
}

func escapeUrl(u string) (template.URL, error) {
	uu, err := url.Parse(u)
	if err != nil {
		return "", err
	}
	uu.RawQuery = uu.Query().Encode()
	return template.URL(uu.String()), nil
}

// Render first of existing widgets
func Coalesce(
	ctx echo.Context,
	ws ...Widget,
) (interface{}, error) {
	for _, w := range ws {
		if w != nil {
			return w.Render(ctx)
		}
	}
	return nil, nil
}

type Mux interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	Any(path string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc) []*echo.Route
}

var (
	TimeInfinity  int64 = 0x7fffffffffffffff
	emptySelected       = make(map[string]bool)
)

func makeTimeout(target int64) time.Duration {
	return time.Unix(target, 0).Sub(time.Now())
}
