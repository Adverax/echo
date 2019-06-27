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

// Widget for display band of items.
// Widget supports pager for navigate between pages.
type Band struct {
	Pager           // Pager info
	Data   DataFunc // Custom data
	Hidden bool     // Band is hidden and not can't be render
}

func (w *Band) Render(ctx echo.Context) (interface{}, error) {
	if w.Hidden {
		return nil, nil
	}

	pager, err := w.Pager.execute(ctx)
	if err != nil {
		return nil, err
	}

	items, err := w.renderItems(ctx, pager.Info)
	if err != nil {
		return nil, err
	}

	res := make(map[string]interface{}, 4)
	res["Items"] = items
	res["Pager"] = pager.render(ctx)

	return res, nil
}

func (w *Band) renderItems(
	ctx echo.Context,
	info PagerInfo,
) ([]interface{}, error) {
	res := make([]interface{}, 0, info.Count)

	for row := 0; row < info.Count; row++ {
		if w.Provider != nil {
			err := w.Provider.Next(ctx)
			if err != nil {
				return nil, err
			}
		}

		item, err := w.Data()
		if err != nil {
			return nil, err
		}
		if item == nil {
			continue
		}

		r, err := echo.RenderWidget(ctx, item)
		if err != nil {
			return nil, err
		}
		if r != nil {
			res = append(res, r)
		}
	}

	return res, nil
}
