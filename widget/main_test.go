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
	stdContext "context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/adverax/echo"
	"github.com/adverax/echo/data"
	"github.com/adverax/echo/generic"
)

type dataMock struct {
	Key   string
	Value string
}

type dataProviderMock struct {
	all   []*dataMock
	rows  []*dataMock
	row   *dataMock
	index int
}

func (provider *dataProviderMock) Count(ctx stdContext.Context) (int, error) {
	return len(provider.rows), nil
}

func (provider *dataProviderMock) Total(ctx stdContext.Context) (int, error) {
	return len(provider.all), nil
}

func (provider *dataProviderMock) Import(ctx stdContext.Context, pagination *data.Pagination) error {
	provider.index = 0
	last := pagination.Offset + pagination.Limit
	if last < int64(len(provider.all)) {
		provider.rows = provider.all[pagination.Offset:last]
	} else {
		provider.rows = provider.all[pagination.Offset:]
	}
	return nil
}

func (provider *dataProviderMock) Next(ctx stdContext.Context) error {
	index := provider.index + 1
	if index > len(provider.rows) {
		return errors.New("range check error")
	}

	provider.row = provider.rows[provider.index]
	provider.index = index
	return nil
}

func TestRenderWidget(t *testing.T) {
	type Test struct {
		src interface{}
		dst string
	}

	provider := &dataProviderMock{
		all: []*dataMock{
			{
				Key:   "Key1",
				Value: "Value1",
			},
			{
				Key:   "Key2",
				Value: "Value2",
			},
		},
	}

	tests := map[string]Test{
		"Not widget": {
			src: "Text message",
			dst: `"Text message"`,
		},
		"Text": {
			src: TEXT("Text message"),
			dst: `"Text message"`,
		},
		"Signed": {
			src: INT(123),
			dst: `"123"`,
		},
		"Unsigned": {
			src: INT(123),
			dst: `"123"`,
		},
		"Decimal": {
			src: FLOAT64(123.5),
			dst: `"123.5"`,
		},
		"MESSAGE": {
			src: MessagePagerNext,
			dst: `"Next"`,
		},
		"Format": {
			src: &MessageFmt{
				Layout: MessageListRecords,
				Params: []interface{}{1, 2, 3},
			},
			dst: `"Shows rows from 1 to 2 of 3"`,
		},
		"Template": {
			src: &Document{
				Layout: "Hello, {{name}}",
				Params: generic.Params{
					"name": "Bob",
				},
			},
			dst: `"Hello, Bob"`,
		},
		"Variant": {
			src: &Variant{
				Formatter: StringFormatter,
				Value:     123,
			},
			dst: `"123"`,
		},
		"Map": {
			src: Map{
				"First":  "first",
				"Second": "second",
			},
			dst: `{"First":"first","Second":"second"}`,
		},
		"List": {
			src: List{
				"first",
				"second",
			},
			dst: `["first","second"]`,
		},

		// Action
		"Action: simple": {
			src: &Action{
				Label:  "label",
				Action: "action",
			},
			dst: `{"Action":"action","Label":"label","Type":"submit"}`,
		},
		"Action: complex": {
			src: &Action{
				Label:    "label",
				Action:   "action",
				Confirm:  "confirm",
				Tooltip:  "tooltip",
				Type:     ActionTypeSubmit,
				Post:     true,
				Disabled: true,
				Name:     "apply",
				Value:    "ok",
			},
			dst: `{"Action":"action","Confirm":"confirm","Disabled":true,"Label":"label","Name":"apply","Post":true,"Tooltip":"tooltip","Type":"submit","Value":"ok"}`,
		},
		"Action: hidden": {
			src: &Action{
				Label:  "label",
				Action: "action",
				Hidden: true,
			},
			dst: `null`,
		},

		"Alert": {
			src: &Alert{
				Type:    AlertSuccess,
				Message: "message",
			},
			dst: `{"Message":"message","Type":"siccess"}`,
		},

		"Band": {
			src: &Band{
				Pager: Pager{
					Capacity: 2,
					BtnCount: 10,
					Url: &url.URL{
						Host: "google.com",
						Path: "/base",
					},
					Provider: provider,
				},
				Data: func() (interface{}, error) {
					return fmt.Sprintf("Key=%s, Value=%s", provider.row.Key, provider.row.Value), nil
				},
			},
			dst: `{"Items":["Key=Key1, Value=Value1","Key=Key2, Value=Value2"],"Pager":{"Count":2,"CurPage":1,"FirstVisible":1,"LastPage":1,"LastVisible":2,"Message":"Shows rows from 1 to 2 of 2","Total":2}}`,
		},

		"Breadcrumbs": {
			src: &Breadcrumbs{
				{
					Label:  "First",
					Action: "action",
				},
				{
					Label: "Second",
				},
			},
			dst: `[{"Action":"action","Label":"First"},{"Label":"Second"}]`,
		},

		"DetailView": {
			src: &DetailView{
				Items: Details{
					"First": {
						Label: TEXT("Label1"),
						Value: 101,
					},
					"Second": {
						Label: "Label2",
						Value: 102,
						Type:  "special",
					},
					"Third": {
						Label:  "Label3",
						Value:  103,
						Hidden: true,
					},
				},
			},
			dst: `{"Body":{"First":{"Label":"Label1","Value":101},"Second":{"Label":"Label2","Type":"special","Value":102}},"Head":{"Key":"Key","Value":"Value"}}`,
		},

		"Table": {
			src: &Table{
				Pager: Pager{
					Provider: provider,
					Url: &url.URL{
						Host: "google.com",
						Path: "/base",
					},
				},
				Columns: TableColumns{
					"Key": {
						Label: "Key",
						Data: func() (interface{}, error) {
							return provider.row.Key, nil
						},
					},
					"Value": {
						Label: "Value",
						Data: func() (interface{}, error) {
							return provider.row.Value, nil
						},
					},
					"Hidden": {
						Label:  "Hidden",
						Hidden: true,
						Data: func() (interface{}, error) {
							return "hidden", nil
						},
					},
				},
			},
			dst: `{"Body":[{"Cols":{"Key":"Key1","Value":"Value1"}},{"Cols":{"Key":"Key2","Value":"Value2"}}],"Head":{"Key":{"Label":"Key"},"Value":{"Label":"Value"}},"Pager":{"Count":2,"CurPage":1,"FirstVisible":1,"LastPage":1,"LastVisible":2,"Message":"Shows rows from 1 to 2 of 2","Total":2}}`,
		},

		"Form": {
			src: &Form{
				Name:   "login-form",
				Action: "google.com",
				Model: echo.Model{
					"Submit": &FormSubmit{
						Label: "Submit",
					},
				},
			},
			dst: `{"Action":"google.com","Method":"POST","Model":{"Submit":{"Label":"Submit"}},"Name":"login-form"}`,
		},

		"MultiForm": {
			src: &MultiForm{
				Name:   "login-form",
				Action: "google.com",
				Models: echo.Models{
					echo.Model{
						"Submit": &FormSubmit{
							Label: "Submit",
						},
					},
				},
			},
			dst: `{"Action":"google.com","Method":"POST","Model":[{"Submit":{"Label":"Submit"}}],"Name":"login-form"}`,
		},

		"NavBar": {
			src: &NavBar{
				Brand: "Brand",
				Items: Map{
					"First": "first",
				},
			},
			dst: `{"Brand":"Brand","Items":{"First":"first"}}`,
		},

		"NavBarItem: simple": {
			src: &NavBarItem{
				Label:  "Label",
				Action: "google.com",
			},
			dst: `{"Action":"google.com","Label":"Label","Type":"item"}`,
		},
		"NavBarItem: separator": {
			src: &NavBarItem{},
			dst: `{"Type":"separator"}`,
		},
		"NavBarItem: active": {
			src: &NavBarItem{
				Label:  "Label",
				Action: "google.com",
				Active: true,
			},
			dst: `{"Action":"google.com","Active":true,"Label":"Label","Type":"item"}`,
		},
		"NavBarItem: hidden": {
			src: &NavBarItem{
				Label:  "Label",
				Action: "google.com",
				Hidden: true,
			},
			dst: `null`,
		},

		"NabBarDropDown: normal": {
			src: &NavBarDropDown{
				Label: "Label",
				Items: List{
					&NavBarItem{
						Label:  "Item",
						Action: "google.com",
					},
				},
			},
			dst: `{"Items":[{"Action":"google.com","Label":"Item","Type":"item"}],"Label":"Label","Type":"menu"}`,
		},
		"NabBarDropDown: hidden": {
			src: &NavBarDropDown{
				Label:  "Label",
				Hidden: true,
				Items: List{
					&NavBarItem{
						Label:  "Item",
						Action: "google.com",
					},
				},
			},
			dst: `null`,
		},

		"NavBarText: normal": {
			src: &NavBarText{
				Body: "Text",
			},
			dst: `{"Items":"Text","Type":"text"}`,
		},
		"NavBarText: hidden": {
			src: &NavBarText{
				Body:   "Text",
				Hidden: true,
			},
			dst: `null`,
		},

		"NavBarForm: normal": {
			src: &NavBarForm{
				Form: Form{
					Model: echo.Model{
						"First": "first",
					},
				},
			},
			dst: `{"Method":"POST","Model":{"First":"first"},"type":"form"}`,
		},
		"NavBarForm: hidden": {
			src: &NavBarForm{
				Form: Form{
					Model: echo.Model{
						"First": "first",
					},
					Hidden: true,
				},
			},
			dst: `null`,
		},
	}

	e := echo.New()
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			c := e.NewContext(nil, nil)
			tree, err := echo.RenderWidget(c, test.src)
			require.NoError(t, err)
			dst, err := json.Marshal(tree)
			require.NoError(t, err)
			require.Equal(t, test.dst, string(dst))
		})
	}
}
