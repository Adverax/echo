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
)

// Widget for display list of key/value pairs.
type DetailView struct {
	Hidden    bool        // Details is hidden
	Items     Details     // Rows
	KeyColumn interface{} // Label for key column
	ValColumn interface{} // Label for value column
}

func (w *DetailView) Render(
	ctx echo.Context,
) (interface{}, error) {
	if w.Hidden {
		return nil, nil
	}

	body := make(map[string]interface{}, len(w.Items))

	for key, item := range w.Items {
		row, err := item.Render(ctx, 0)
		if err != nil {
			return nil, err
		}
		if row != nil {
			body[key] = row
		}
	}

	res := make(map[string]interface{}, 4)
	res["Body"] = body

	keyColumn := w.KeyColumn
	if keyColumn == nil {
		keyColumn = MessageDetailsColumnKey
	}
	key, err := RenderWidget(ctx, keyColumn)
	if err != nil {
		return nil, err
	}

	valColumn := w.ValColumn
	if valColumn == nil {
		valColumn = MessageDetailsColumnVal
	}
	value, err := RenderWidget(ctx, valColumn)
	if err != nil {
		return nil, err
	}

	res["Head"] = map[string]interface{}{
		"Key":   key,
		"Value": value,
	}

	return res, nil
}

type Details map[string]Detail

type Detail struct {
	Hidden bool        // Detail is hidden
	Label  interface{} // Detail title
	Value  interface{} // Value
	Type   string      // Column type (optional)
}

func (w *Detail) Render(
	ctx echo.Context,
	level int,
) (interface{}, error) {
	if w.Hidden {
		return nil, nil
	}

	res := make(map[string]interface{}, 16)

	if w.Label != nil {
		label, err := RenderWidget(ctx, w.Label)
		if err != nil {
			return nil, err
		}
		res["Label"] = label
	}

	if w.Value != nil {
		value, err := RenderWidget(ctx, w.Value)
		if err != nil {
			return nil, err
		}
		res["Value"] = value
	}

	if w.Type != "" {
		res["Type"] = w.Type
	}

	if level != 0 {
		res["Level"] = level
	}

	return res, nil
}
