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
	"testing"

	"github.com/stretchr/testify/assert"
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
	ctx := e.NewContext(nil, nil)
	// Check enumeration
	var keys, values string
	err := ds.Enumerate(
		ctx,
		func(key, value string) error {
			keys += key
			values += value
			return nil
		},
	)
	assert.NoError(t, err)

	assert.Equal(t, "abc", keys)
	assert.Equal(t, "123", values)

	_, err = ds.Decode(ctx, "a")
	assert.NoError(t, err)

	_, err = ds.Decode(ctx, "e")
	assert.Error(t, err)
}

func TestParseDataSet(t *testing.T) {
	ds := ParseDataSet(`
#! MAP SORTED DELIMITER ::
1::London
2::New York
3::Paris
`,
	)

	e := New()
	ctx := e.NewContext(nil, nil)
	check := func(key, value string) {
		_, err := ds.Decode(ctx, key)
		assert.NoError(t, err)
	}

	check("1", "London")
	check("2", "New York")
	check("3", "Paris")

	v1, err := ds.Decode(ctx, "1")
	assert.NoError(t, err)
	assert.Equal(t, "London", v1)

	v2, err := ds.Encode(ctx, "London")
	assert.NoError(t, err)
	assert.Equal(t, "1", v2)
}
