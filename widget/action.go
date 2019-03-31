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

const (
	ActionTypeButton ActionType = "button"
	ActionTypeSubmit ActionType = "submit"
	ActionTypeReset  ActionType = "reset"
)

type ActionType string

// Widget for display button or hyperlink.
type Action struct {
	Label    interface{} // Visible label
	Action   interface{} // Action
	Confirm  interface{} // Confirmation text
	Tooltip  interface{} // Widget tooltip
	Type     ActionType  // Type of action
	Post     bool        // Use post request
	Hidden   bool        // Action is hidden
	Disabled bool        // Action is disabled
	Name     string      // Name of action
	Value    interface{} // Value of action
}

func (w *Action) Render(
	ctx echo.Context,
) (interface{}, error) {
	if w.Hidden {
		return nil, nil
	}

	res := make(map[string]interface{}, 16)

	if w.Label != nil {
		label, err := echo.RenderWidget(ctx, w.Label)
		if err != nil {
			return nil, err
		}
		res["Label"] = label
	}

	if w.Action != nil {
		action, err := RenderLink(ctx, w.Action)
		if err != nil {
			return nil, err
		}
		res["Action"] = action
	}

	if w.Confirm != nil {
		confirm, err := echo.RenderWidget(ctx, w.Confirm)
		if err != nil {
			return nil, err
		}
		res["Confirm"] = confirm
	}

	if w.Tooltip != nil {
		tooltip, err := echo.RenderWidget(ctx, w.Tooltip)
		if err != nil {
			return nil, err
		}
		res["Tooltip"] = tooltip
	}

	if w.Value != nil {
		value, err := echo.RenderWidget(ctx, w.Value)
		if err != nil {
			return nil, err
		}
		res["Value"] = value
	}

	if w.Name != "" {
		res["Name"] = w.Name
	}

	if w.Disabled {
		res["Disabled"] = true
	}

	if w.Post {
		res["Post"] = true
	}

	tp := w.Type
	if tp == "" {
		tp = ActionTypeSubmit
	}
	res["Type"] = string(tp)

	return res, nil
}
