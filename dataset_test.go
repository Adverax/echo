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
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewDataSet(t *testing.T) {
	ds := NewDataSet(
		map[string]string{
			"a": "1",
			"b": "2",
			"c": "3",
		},
		true,
	)

	e := New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Check enumeration
	var keys, values string
	ds.Enumerate(c, func(key, value string) error {
		keys += key
		values += value
		return nil
	})

	assert.Equal(t, "abc", keys)
	assert.Equal(t, "123", values)

	assert.True(t, ds.Has("a"))
	assert.False(t, ds.Has("e"))
}

func TestParseDataSet(t *testing.T) {
	ds := ParseDataSet(`
#! MAP SORTED DELIMITER ::
1::London
2::New York
3::Paris
`,
	)

	check := func(key, value string) {
		assert.True(t, ds.Has(key))
	}

	check("1", "London")
	check("2", "New York")
	check("3", "Paris")
}
