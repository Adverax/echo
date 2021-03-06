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

// Widget for display primary navigation bar.
// Example:
//   menu := &widget.NavBar{
//		Brand: "My brand",
//		Items: widget.Map{
//			"Home": &widget.NavBarItem{
//		         Label:  "Home",
//		         Action: "/",
//	        },
//			"Auth": &widget.NavBarItem{
//				Label:  "Login",
//				Action: "/login",
//			},
//			"Language": &widget.NavBarDropDown{
//				Label: &widget.Sprintf{
//					Layout: "Language (%s)",
//					Params: []interface{}{curLang},
//				},
//				Items: widget.List{
//                  &widget.NavBarItem{
//				        Label:  "English",
//				        Action: "/lang/english",
//			        },
//                  &widget.NavBarItem{
//				        Label:  "Russian",
//				        Action: "/lang/russian",
//			        },
//				},
//			},
//		},
//	 }
type NavBar struct {
	Brand interface{} // Brand
	Items Map         // Nav bar items
}

func (w *NavBar) Render(ctx echo.Context) (interface{}, error) {
	res := make(map[string]interface{}, 16)
	if w.Brand != nil {
		brand, err := echo.RenderWidget(ctx, w.Brand)
		if err != nil {
			return nil, err
		}
		res["Brand"] = brand
	}

	items, err := RenderMap(ctx, w.Items)
	if err != nil {
		return nil, err
	}
	res["Items"] = items

	return res, nil
}

// Simple menu item of navigation bar.
// Must be nested for NavBarNav or NavBarDropDown.
type NavBarItem struct {
	Label  interface{} // Label of element
	Action interface{} // Url of element
	Hidden bool        // Element is hidden and can't be render
	Active bool        // Element is active (highlighted)
}

func (w *NavBarItem) Render(ctx echo.Context) (interface{}, error) {
	if w.Hidden {
		return nil, nil
	}

	res := make(map[string]interface{}, 16)
	if w.Action == nil {
		res["Type"] = "separator"
		return res, nil
	}
	if w.Label == nil {
		return nil, nil
	}

	if w.Active {
		res["Active"] = true
	}

	label, err := echo.RenderWidget(ctx, w.Label)
	if err != nil {
		return nil, err
	}
	res["Label"] = label

	action, err := RenderLink(ctx, w.Action)
	if err != nil {
		return nil, err
	}
	res["Action"] = action
	res["Type"] = "item"

	return res, nil
}

// Drop down submenu of navigation bar.
// Must be nested for NavBarNav
type NavBarDropDown struct {
	Label      interface{} // Title of submenu
	Items      List        // Items of submenu (NavBarItem)
	Hidden     bool        // Submenu is hidden and can't be render
	AllowEmpty bool        // Allow empty submenu (otherwise hide)
}

func (w *NavBarDropDown) Render(ctx echo.Context) (interface{}, error) {
	if w.Hidden || w.Label == nil {
		return nil, nil
	}

	label, err := echo.RenderWidget(ctx, w.Label)
	if err != nil {
		return nil, err
	}

	items, err := RenderList(ctx, w.Items)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 && !w.AllowEmpty {
		return nil, nil
	}
	res := make(map[string]interface{}, 16)
	res["Items"] = items
	res["Label"] = label
	res["Type"] = "menu"

	return res, nil
}

// Simple text, nested in navigatin bar.
// Must be nested for NavBarNav (not navbar link)
type NavBarText struct {
	Body   interface{} // Text
	Hidden bool        // Text is hidden and can't be render
}

func (w *NavBarText) Render(ctx echo.Context) (interface{}, error) {
	if w.Hidden || w.Body == nil {
		return nil, nil
	}

	res := make(map[string]interface{}, 16)
	items, err := echo.RenderWidget(ctx, w.Body)
	if err != nil {
		return nil, err
	}
	res["Items"] = items
	res["Type"] = "text"

	return res, nil
}

// Hyperlink, nested in navigation bar.
// Must be nested for NavBarText (not in navbar link)
type NavBarLink struct {
	Label   interface{} // Hyperlink label
	Action  interface{} // Hyperlink action
	Tooltip interface{} // Hyperlink tooltip
	Post    bool        // Hyperlink used post
	Hidden  bool        // Hyperlink is hidden and can't be render
}

func (w *NavBarLink) Render(ctx echo.Context) (interface{}, error) {
	if w.Hidden || w.Label == nil {
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

	if w.Label != nil {
		label, err := echo.RenderWidget(ctx, w.Label)
		if err != nil {
			return nil, err
		}
		res["Label"] = label
	}

	if w.Post {
		res["Post"] = true
	}

	if w.Tooltip != nil {
		tooltip, err := echo.RenderWidget(ctx, w.Tooltip)
		if err != nil {
			return nil, err
		}
		res["Tooltip"] = tooltip
	}

	res["type"] = "link"

	return res, nil
}

// Inline form, nested in navigation bar
// Must be nested for NavBarNav (not navbar link)
type NavBarForm struct {
	Form // Form elements
}

func (w *NavBarForm) Render(ctx echo.Context) (interface{}, error) {
	if w.Hidden {
		return nil, nil
	}

	res, err := w.Form.Render(ctx)
	if err != nil {
		return nil, err
	}

	result := res.(map[string]interface{})
	result["type"] = "form"

	return result, nil
}
