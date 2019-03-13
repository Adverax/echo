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
	"github.com/adverax/echo/generic"
)

type TableColumn struct {
	Label    interface{}  // Column label
	Hidden   bool         // Column is hidden
	Data     DataFunc     // Column data
	Expander ExpanderFund // Column expander
}

type TableAction struct {
	Action  interface{} // Custom action
	Tooltip interface{} // Action tooltip
	Confirm interface{} // Confirmation text
	Post    bool        // Use post request
	Hidden  bool        // Action is hidden
}

func (w *TableAction) Render(
	ctx echo.Context,
) (interface{}, error) {
	if w.Hidden {
		return nil, nil
	}

	res := make(map[string]interface{}, 16)
	if w.Action != nil {
		action, err := RenderLink(ctx, w.Action)
		if err != nil {
			return nil, err
		}
		res["Action"] = action
	}

	if w.Tooltip != nil {
		tooltip, err := RenderWidget(ctx, w.Tooltip)
		if err != nil {
			return nil, err
		}
		res["Tooltip"] = tooltip
	}

	if w.Post {
		res["Post"] = true
	}

	if w.Confirm != nil {
		confirm, err := RenderWidget(ctx, w.Confirm)
		if err != nil {
			return nil, err
		}
		res["Confirm"] = confirm
	}

	return res, nil
}

var DefaultTableActionView = &TableAction{}
var DefaultTableActionUpdate = &TableAction{}
var DefaultTableActionDelete = &TableAction{}

type TableActions struct {
	RowId interface{} // Row identifier
	Path  string      // Base path to family of actions (optional)
	Items Map         // Action map
}

func (w *TableActions) Render(
	ctx echo.Context,
) (interface{}, error) {
	path := w.Path
	if path == "" {
		path = ctx.Request().URL.Path
		path = ctx.Echo().UrlLinker.Collapse(ctx, path)
	}
	rowId, _ := generic.ConvertToString(w.RowId)

	actions := make(map[string]interface{}, len(w.Items))
	for key, act := range w.Items {
		switch act {
		case DefaultTableActionView:
			key = "View"
			act = WidgetFunc(func(ctx echo.Context) (interface{}, error) {
				return (&TableAction{
					Action:  path + "/view/" + rowId,
					Tooltip: MessageTableTooltipActionView,
				}).Render(ctx)
			})
		case DefaultTableActionUpdate:
			key = "Update"
			act = WidgetFunc(func(ctx echo.Context) (interface{}, error) {
				return (&TableAction{
					Action:  path + "/update/" + rowId,
					Tooltip: MessageTableTooltipActionUpdate,
				}).Render(ctx)
			})
		case DefaultTableActionDelete:
			key = "Delete"
			act = WidgetFunc(func(ctx echo.Context) (interface{}, error) {
				return (&TableAction{
					Action:  path + "/delete/" + rowId,
					Tooltip: MessageTableTooltipActionDelete,
					Confirm: MessageTableConfirmActionDelete,
					Post:    true,
				}).Render(ctx)
			})
		}

		action, err := RenderWidget(ctx, act)
		if err != nil {
			return nil, err
		}
		if action != nil {
			actions[key] = action
		}
	}

	return actions, nil
}

type TableColumns map[string]*TableColumn

// Widget for display simple html table.
type Table struct {
	Pager
	Columns     TableColumns // Columns declaration
	RowExpander ExpanderFund // Row expander
}

func (w *Table) Render(
	ctx echo.Context,
) (interface{}, error) {
	pager, err := w.Pager.execute(ctx)
	if err != nil {
		return nil, err
	}

	head, err := w.renderHead(ctx)
	if err != nil {
		return nil, err
	}

	body, err := w.renderBody(ctx, &pager.Info)
	if err != nil {
		return nil, err
	}

	res := map[string]interface{}{
		"Head":  head,
		"Body":  body,
		"Pager": pager.render(ctx),
	}

	return res, nil
}

func (w *Table) renderHead(
	ctx echo.Context,
) (interface{}, error) {
	cells := make(map[string]interface{}, len(w.Columns))

	for key, col := range w.Columns {
		if col == nil || col.Hidden {
			continue
		}

		cell := make(map[string]interface{}, 4)
		if col.Label == nil {
			cell["Label"] = ""
		} else {
			label, err := RenderWidget(ctx, col.Label)
			if err != nil {
				return nil, err
			}
			cell["Label"] = label
		}

		if col.Expander != nil {
			err := col.Expander(cell)
			if err != nil {
				return nil, err
			}
		}

		cells[key] = cell
	}

	return cells, nil
}

func (w *Table) renderBody(
	ctx echo.Context,
	info *PagerInfo,
) (interface{}, error) {
	rows := make([]interface{}, 0, info.Count)
	for row := 0; row < info.Count; row++ {
		err := w.Provider.Next(ctx)
		if err != nil {
			return nil, err
		}

		r, err := w.renderRow(ctx)
		if err != nil {
			return nil, err
		}
		if r != nil {
			rows = append(rows, r)
		}
	}

	return rows, nil
}

func (w *Table) renderRow(
	ctx echo.Context,
) (interface{}, error) {
	res := make(map[string]interface{}, 4)
	cols := make(map[string]interface{}, len(w.Columns))
	res["Cols"] = cols

	for col, c := range w.Columns {
		if c == nil || c.Hidden {
			continue
		}

		value, err := c.Data()
		if err != nil {
			return nil, err
		}

		val, err := RenderWidget(ctx, value)
		if err != nil {
			return nil, err
		}

		cols[col] = val
	}

	if w.RowExpander != nil {
		err := w.RowExpander(res)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}
