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

package generic

import (
	"reflect"
)

// Deep copy value
// Usage^
//  type data struct {
//  	a string
//  	b string
//  	c []string
//  }
//  old := &data{
//   "works1",
//   "works2",
//    []string{"a", "b"},
//  }
//  var new *data = &data{}
//  CloneValueTo(old, new)
//  old.c = append(old.c, "c")
//  fmt.Println(old)
// fmt.Println(new)
// So you can pass any type at run time as long as you're sure that source and
// destin are both of the same type, (and destin is a pointer to that type).
func CloneValueTo(dst interface{}, src interface{}) {
	x := reflect.ValueOf(src)
	if x.Kind() == reflect.Ptr {
		starX := x.Elem()
		y := reflect.New(starX.Type())
		starY := y.Elem()
		starY.Set(starX)
		reflect.ValueOf(dst).Elem().Set(y.Elem())
	} else {
		dst = x.Interface()
	}
}

func CloneValue(src interface{}) interface{} {
	x := reflect.ValueOf(src)
	if x.Kind() == reflect.Ptr {
		starX := x.Elem()
		y := reflect.New(starX.Type())
		starY := y.Elem()
		starY.Set(starX)
		return starY.Interface()
	} else {
		return x.Interface()
	}
}

func MakePointerTo(obj interface{}) interface{} {
	val := reflect.ValueOf(obj)
	vp := reflect.New(val.Type())
	vp.Elem().Set(val)
	return vp.Interface()
}
