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
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
)

func TestPager_execute(t *testing.T) {
	type Test struct {
		src *Pager
		dst string
	}

	tests := map[string]Test{
		"Empty set": {
			src: &Pager{
				Capacity: 10,
				BtnCount: 4,
				Url: &url.URL{
					Host: "google.com",
					Path: "/base",
				},
				Provider: &dataProviderMock{
					all: []*dataMock{},
				},
			},
			dst: `{"info":{"CurPage":1,"LastPage":1,"FirstVisible":0,"LastVisible":0,"Total":0,"Count":0},"message":"No data for display"}`,
		},
		"Single page": {
			src: &Pager{
				Capacity: 10,
				BtnCount: 4,
				Url: &url.URL{
					Host: "google.com",
					Path: "/base",
				},
				Provider: &dataProviderMock{
					all: []*dataMock{
						{Key: "Key1"},
						{Key: "Key2"},
					},
				},
			},
			dst: `{"info":{"CurPage":1,"LastPage":1,"FirstVisible":1,"LastVisible":2,"Total":2,"Count":2},"message":"Shows rows from 1 to 2 of 2"}`,
		},
		"First page": {
			src: &Pager{
				Capacity: 2,
				BtnCount: 4,
				Url: &url.URL{
					Host:     "google.com",
					Path:     "/base",
					RawQuery: "pg=1",
				},
				Provider: &dataProviderMock{
					all: []*dataMock{
						{Key: "Key1"},
						{Key: "Key2"},
						{Key: "Key3"},
						{Key: "Key4"},
						{Key: "Key5"},
						{Key: "Key6"},
						{Key: "Key7"},
						{Key: "Key8"},
						{Key: "Key9"},
						{Key: "Key10"},
					},
				},
			},
			dst: `{"info":{"CurPage":1,"LastPage":5,"FirstVisible":1,"LastVisible":2,"Total":10,"Count":2},"buttons":{"Band":[{"Action":"//google.com/base?pg=1","Active":true,"Label":"1"},{"Action":"//google.com/base?pg=2","Label":"2"},{"Action":"//google.com/base?pg=3","Label":"3"},{"Action":"//google.com/base?pg=4","Label":"4"}],"Next":{"Action":"//google.com/base?pg=2","Label":"Next"},"Prev":{"Action":"//google.com/base?pg=0","Disabled":true,"Label":"Prev"}},"message":"Shows rows from 1 to 2 of 10"}`,
		},
		"Last page": {
			src: &Pager{
				Capacity: 2,
				BtnCount: 4,
				Url: &url.URL{
					Host:     "google.com",
					Path:     "/base",
					RawQuery: "pg=5",
				},
				Provider: &dataProviderMock{
					all: []*dataMock{
						{Key: "Key1"},
						{Key: "Key2"},
						{Key: "Key3"},
						{Key: "Key4"},
						{Key: "Key5"},
						{Key: "Key6"},
						{Key: "Key7"},
						{Key: "Key8"},
						{Key: "Key9"},
						{Key: "Key10"},
					},
				},
			},
			dst: `{"info":{"CurPage":5,"LastPage":5,"FirstVisible":9,"LastVisible":10,"Total":10,"Count":2},"buttons":{"Band":[{"Action":"//google.com/base?pg=2","Label":"2"},{"Action":"//google.com/base?pg=3","Label":"3"},{"Action":"//google.com/base?pg=4","Label":"4"},{"Action":"//google.com/base?pg=5","Active":true,"Label":"5"}],"Next":{"Action":"//google.com/base?pg=6","Disabled":true,"Label":"Next"},"Prev":{"Action":"//google.com/base?pg=4","Label":"Prev"}},"message":"Shows rows from 9 to 10 of 10"}`,
		},
		"Middle page": {
			src: &Pager{
				Capacity: 2,
				BtnCount: 4,
				Url: &url.URL{
					Host:     "google.com",
					Path:     "/base",
					RawQuery: "pg=3",
				},
				Provider: &dataProviderMock{
					all: []*dataMock{
						{Key: "Key1"},
						{Key: "Key2"},
						{Key: "Key3"},
						{Key: "Key4"},
						{Key: "Key5"},
						{Key: "Key6"},
						{Key: "Key7"},
						{Key: "Key8"},
						{Key: "Key9"},
						{Key: "Key10"},
					},
				},
			},
			dst: `{"info":{"CurPage":3,"LastPage":5,"FirstVisible":5,"LastVisible":6,"Total":10,"Count":2},"buttons":{"Band":[{"Action":"//google.com/base?pg=2","Label":"2"},{"Action":"//google.com/base?pg=3","Active":true,"Label":"3"},{"Action":"//google.com/base?pg=4","Label":"4"},{"Action":"//google.com/base?pg=5","Label":"5"}],"Next":{"Action":"//google.com/base?pg=4","Label":"Next"},"Prev":{"Action":"//google.com/base?pg=2","Label":"Prev"}},"message":"Shows rows from 5 to 6 of 10"}`,
		},
	}

	e := echo.New()
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			c := e.NewContext(nil, nil)
			tree, err := test.src.execute(c)
			require.NoError(t, err)
			dst, err := json.Marshal(tree)
			require.NoError(t, err)
			require.Equal(t, test.dst, string(dst))
		})
	}
}
