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
	"strings"

	"github.com/adverax/echo"
)

type Template interface {
	Execute(wr io.Writer, data interface{}) error
}

type Designer interface {
	// Parse templates (with relative paths or aliases, starts with "@")
	Compile(files ...string) Template
	// Create new child designer.
	// Method create new Designer with related path.
	// Method extends set of funcs and views.
	Extends(funcs template.FuncMap, path string, views ...string) Designer
	// Get application
	Echo() *echo.Echo
}

type designer struct {
	echo    *echo.Echo
	layouts string           // Path to the root layouts folder
	path    string           // Path to the views folder
	funcs   template.FuncMap // Map of shared funcs
	tpl     *template.Template
}

func (d *designer) Echo() *echo.Echo {
	return d.echo
}

func (d *designer) Extends(
	funcs template.FuncMap,
	path string,
	files ...string,
) Designer {
	fs := make(template.FuncMap)
	for k, v := range d.funcs {
		fs[k] = v
	}
	for k, v := range funcs {
		fs[k] = v
	}

	tpl, err := d.tpl.Clone()
	if err != nil {
		panic(err)
	}
	tpl.Funcs(fs)

	if len(files) != 0 {
		for i, file := range files {
			files[i] = d.layouts + file
		}

		template.Must(tpl.ParseFiles(files...))
	}

	return &designer{
		echo:    d.echo,
		layouts: d.layouts,
		path:    addTrailingSlash(d.path + path),
		funcs:   fs,
		tpl:     tpl,
	}
}

func (d *designer) Compile(files ...string) Template {
	for i, file := range files {
		files[i] = d.path + file
	}

	tpl, err := d.tpl.Clone()
	if err != nil {
		panic(err)
	}

	return template.Must(tpl.Funcs(d.funcs).ParseFiles(files...))
}

func NewDesigner(
	e *echo.Echo,
	funcs template.FuncMap,
	path string, // Path to the views folder
	layouts string, // Path to the layouts folder
	files ...string, // Loaded views
) Designer {
	path = addTrailingSlash(path)
	layouts = addTrailingSlash(layouts)

	if funcs == nil {
		funcs = make(template.FuncMap)
	}

	for i, file := range files {
		files[i] = layouts + file
	}

	tpl := template.Must(template.New("root").Funcs(funcs).ParseFiles(files...))

	return &designer{
		echo:    e,
		layouts: layouts,
		path:    path,
		funcs:   funcs,
		tpl:     tpl,
	}
}

func addTrailingSlash(path string) string {
	if strings.HasSuffix(path, "/") {
		return path
	}

	return path + "/"
}
