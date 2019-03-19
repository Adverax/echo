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

package design

import (
	"html/template"
	"io"
	"path/filepath"
	"strings"

	"github.com/adverax/echo"
)

type Template interface {
	Execute(wr io.Writer, data interface{}) error
}

type Designer interface {
	// Parse templates (with relative paths)
	ParseFiles(files ...string) Template
	// Create new child designer.
	// Method create new Designer with related path.
	// Method extends set of funcs and views.
	Extends(funcs template.FuncMap, path string, views ...string) Designer
	// Get application
	Echo() *echo.Echo
}

type designer struct {
	echo    *echo.Echo
	layouts string                        // Path to the root layouts folder
	path    string                        // Path to the views folder
	funcs   template.FuncMap              // Map of shared funcs
	views   map[string]*template.Template // Map of shared views
}

func (d *designer) Echo() *echo.Echo {
	return d.echo
}

func (d *designer) Extends(
	funcs template.FuncMap,
	path string,
	views ...string,
) Designer {
	vs := make(map[string]*template.Template)
	fs := make(template.FuncMap)
	// Clone parent view list
	for k, v := range d.views {
		vs[k] = v
	}
	for k, v := range d.funcs {
		fs[k] = v
	}

	// Extends view list
	for _, view := range views {
		vs[view] = template.Must(template.New(view).ParseFiles(d.layouts + view))
	}
	for k, v := range funcs {
		fs[k] = v
	}

	// Create new designer
	return &designer{
		echo:    d.echo,
		layouts: d.layouts,
		path:    d.path + path,
		funcs:   fs,
		views:   vs,
	}
}

func (d *designer) ParseFiles(files ...string) Template {
	tpl := template.New("main").Funcs(d.funcs)
	var list []string

	for _, file := range files {
		if strings.HasPrefix(file, "@") {
			file = "/" + file[1:]
			name := getFileName(file)
			if item, has := d.views[name]; has {
				tpl = template.Must(tpl.AddParseTree(item.Name(), item.Tree))
			} else {
				tpl = template.Must(tpl.ParseFiles(d.layouts + file))
			}
		} else {
			list = append(list, d.path+file)
		}
	}

	return template.Must(tpl.ParseFiles(list...))
}

func NewDesigner(
	e *echo.Echo,
	funcs template.FuncMap,
	path string, // Path to the views folder
	layouts string, // Path to the layouts folder
	views ...string, // Loaded views
) Designer {
	if funcs == nil {
		funcs = make(template.FuncMap)
	}

	vs := make(map[string]*template.Template)

	for _, view := range views {
		name := getFileName(view)
		vs[name] = template.Must(template.ParseFiles(layouts + view))
	}

	return &designer{
		echo:    e,
		layouts: layouts,
		path:    path,
		funcs:   funcs,
		views:   vs,
	}
}

func getFileName(file string) string {
	if strings.HasPrefix(file, "/") {
		file = file[1:]
	}

	return strings.TrimSuffix(file, filepath.Ext(file))
}
