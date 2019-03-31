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

type Breadcrumb struct {
	Label  interface{} // Label for action
	Action interface{} // Action (optional)
}

// Widget for display path in the tree navigation.
type Breadcrumbs []*Breadcrumb

func (w Breadcrumbs) Render(
	ctx echo.Context,
) (interface{}, error) {
	res := make([]interface{}, 0, len(w))

	for _, breadcrumb := range w {
		if breadcrumb.Label == nil {
			continue
		}

		label, err := echo.RenderWidget(ctx, breadcrumb.Label)
		if err != nil {
			return nil, err
		}
		if label == nil {
			continue
		}

		item := make(map[string]interface{}, 4)
		item["Label"] = label

		if breadcrumb.Action != nil {
			action, err := RenderLink(ctx, breadcrumb.Action)
			if err != nil {
				return nil, err
			}
			item["Action"] = action
		}

		res = append(res, item)
	}

	return res, nil
}
