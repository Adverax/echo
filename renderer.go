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
	"fmt"
	"html/template"
	"strings"
)

// Template engine.
// Root template must be without {{DEFINE}} wrapper.
// Example:
// ...
// {{template "content" .}}
// ...
//
// Usage example:
// func setup(e *echo.Echo) {
// 	 r := echo.NewRenderer(
// 		 "/path/to/layouts",
// 		 map[string]string{
// 			 "main": "/main.tmpl",
// 		 },
// 		 nil,
// 	 )
//	 tpl1 := r.ParseFiles("@main", "/path/to/views/view-name.tmpl")
//   ...
//   err := tpl1.Execute(w, data)
// }
type Renderer interface {
	// Parse template with relative path
	// If name starts with char "@" - from layout anchor
	// Else starts from views anchor
	ParseFiles(files ...string) *template.Template
}

type renderer struct {
	registry map[string]*template.Template
	funcs    template.FuncMap // Base set of functions
}

func (renderer *renderer) ParseFiles(files ...string) *template.Template {
	tpl := template.New("main").Funcs(renderer.funcs)
	var list []string
	for _, file := range files {
		if strings.HasPrefix(file, "@") {
			name := file[1:]
			item, has := renderer.registry[name]
			if !has {
				panic(fmt.Errorf("unknown layout file %q", name))
			}
			tpl = template.Must(tpl.AddParseTree(name, item.Tree))
		} else {
			list = append(list, file)
		}
	}

	if len(list) == 0 {
		return tpl
	}

	return template.Must(tpl.ParseFiles(list...))
}

// Create new template renderer
func NewRenderer(
	layouts string, // Path to the folder with layouts
	files map[string]string, // Files
	funcs template.FuncMap, // Used functions
) Renderer {
	registry := make(map[string]*template.Template, 0)
	for key, val := range files {
		path := layouts + val
		registry[key] = template.Must(template.ParseFiles(path))
	}

	return &renderer{
		registry: registry,
		funcs:    funcs,
	}
}
