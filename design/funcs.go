package design

import (
	"fmt"
	"html/template"
	"math"
	"reflect"
	"sort"
	"strings"

	"github.com/adverax/echo/generic"
)

// Numeric functions

// Sum numbers with `add`
func add(v ...interface{}) int64 {
	var a int64 = 0
	for _, b := range v {
		val, _ := generic.ConvertToInt64(b)
		a += val
	}
	return a
}

// To subtract, use `sub`
func sub(a, b interface{}) int64 {
	aa, _ := generic.ConvertToInt64(a)
	bb, _ := generic.ConvertToInt64(b)
	return aa - bb
}

// Perform integer division with `div`
func div(a, b interface{}) int64 {
	aa, _ := generic.ConvertToInt64(a)
	bb, _ := generic.ConvertToInt64(b)
	return aa / bb
}

// Modulo with `mod`
func mod(a, b interface{}) int64 {
	aa, _ := generic.ConvertToInt64(a)
	bb, _ := generic.ConvertToInt64(b)
	return aa % bb
}

// Multiply with `mul`
func mul(a interface{}, v ...interface{}) int64 {
	val, _ := generic.ConvertToInt64(a)
	for _, b := range v {
		bb, _ := generic.ConvertToInt64(b)
		val = val * bb
	}
	return val
}

// Return the largest of a series of integers:
//     `max 1 2 3` will return  `3`.
func max(a interface{}, i ...interface{}) int64 {
	aa, _ := generic.ConvertToInt64(a)
	for _, b := range i {
		bb, _ := generic.ConvertToInt64(b)
		if bb > aa {
			aa = bb
		}
	}
	return aa
}

// Return the smallest of a series of integers.
//     `min 1 2 3` will return `1`.
func min(a interface{}, i ...interface{}) int64 {
	aa, _ := generic.ConvertToInt64(a)
	for _, b := range i {
		bb, _ := generic.ConvertToInt64(b)
		if bb < aa {
			aa = bb
		}
	}
	return aa
}

// Returns the greatest float value less than or equal to input value
//     `floor 123.9999`` will return `123.0`
func floor(a interface{}) float64 {
	aa, _ := generic.ConvertToFloat64(a)
	return math.Floor(aa)
}

// Returns the greatest float value greater than or equal to input value
//     `ceil 123.001` will return `124.0`
func ceil(a interface{}) float64 {
	aa, _ := generic.ConvertToFloat64(a)
	return math.Ceil(aa)
}

// Returns a float value with the remainder rounded to the given number to digits after the decimal point.
//     `round 123.555555` will return `123.556`
func round(a interface{}, p int, r_opt ...float64) float64 {
	roundOn := .5
	if len(r_opt) > 0 {
		roundOn = r_opt[0]
	}
	val, _ := generic.ConvertToFloat64(a)
	places, _ := generic.ConvertToFloat64(p)

	var round float64
	pow := math.Pow(10, places)
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	return round / pow
}

// Functions for handle lists.
// Simple `list` type that can contain arbitrary sequential lists
// of data. This is similar to arrays or slices, but lists are designed to be used
// as immutable data types.

// list creates new slice of items
//     `$myList := list 1 2 3 4 5` will return new list `1,2,3,4,5`
func list(v ...interface{}) []interface{} {
	return v
}

// Append a new item to an existing list, creating a new list.
//     `$new = append $myList 6`
// The above would set `$new` to `[1 2 3 4 5 6]`. `$myList` would remain unaltered.
func push(list interface{}, v interface{}) []interface{} {
	tp := reflect.TypeOf(list).Kind()
	switch tp {
	case reflect.Slice, reflect.Array:
		l2 := reflect.ValueOf(list)

		l := l2.Len()
		nl := make([]interface{}, l)
		for i := 0; i < l; i++ {
			nl[i] = l2.Index(i).Interface()
		}

		return append(nl, v)

	default:
		panic(fmt.Errorf("cannot push on type %s", tp))
	}
}

// Push an alement onto the front of a list, creating a new list.
//     `prepend $myList 0`
//The above would produce `[0 1 2 3 4 5]`. `$myList` would remain unaltered.
func prepend(list interface{}, v interface{}) []interface{} {
	//return append([]interface{}{v}, list...)

	tp := reflect.TypeOf(list).Kind()
	switch tp {
	case reflect.Slice, reflect.Array:
		l2 := reflect.ValueOf(list)

		l := l2.Len()
		nl := make([]interface{}, l)
		for i := 0; i < l; i++ {
			nl[i] = l2.Index(i).Interface()
		}

		return append([]interface{}{v}, nl...)

	default:
		panic(fmt.Errorf("cannot prepend on type %s", tp))
	}
}

// To get the last item on a list, use `last`:
//   `last $myList` returns `5`. This is roughly analogous to reversing a list and
// then calling `first`.
func last(list interface{}) interface{} {
	tp := reflect.TypeOf(list).Kind()
	switch tp {
	case reflect.Slice, reflect.Array:
		l2 := reflect.ValueOf(list)

		l := l2.Len()
		if l == 0 {
			return nil
		}

		return l2.Index(l - 1).Interface()
	default:
		panic(fmt.Errorf("cannot find last on type %s", tp))
	}
}

// To get the head item on a list, use `first`.
//     `first $myList` will returns `1`
func first(list interface{}) interface{} {
	tp := reflect.TypeOf(list).Kind()
	switch tp {
	case reflect.Slice, reflect.Array:
		l2 := reflect.ValueOf(list)

		l := l2.Len()
		if l == 0 {
			return nil
		}

		return l2.Index(0).Interface()
	default:
		panic(fmt.Errorf("cannot find first on type %s", tp))
	}
}

// To get the tail of the list (everything but the first item), use `rest`.
//    `rest $myList` will returns `[2 3 4 5]`
func rest(list interface{}) []interface{} {
	tp := reflect.TypeOf(list).Kind()
	switch tp {
	case reflect.Slice, reflect.Array:
		l2 := reflect.ValueOf(list)

		l := l2.Len()
		if l == 0 {
			return nil
		}

		nl := make([]interface{}, l-1)
		for i := 1; i < l; i++ {
			nl[i-1] = l2.Index(i).Interface()
		}

		return nl
	default:
		panic(fmt.Errorf("cannot find rest on type %s", tp))
	}
}

// This compliments `last` by returning all _but_ the last element.
//     `initial $myList` returns `[1 2 3 4]`.
func initial(list interface{}) []interface{} {
	tp := reflect.TypeOf(list).Kind()
	switch tp {
	case reflect.Slice, reflect.Array:
		l2 := reflect.ValueOf(list)

		l := l2.Len()
		if l == 0 {
			return nil
		}

		nl := make([]interface{}, l-1)
		for i := 0; i < l-1; i++ {
			nl[i] = l2.Index(i).Interface()
		}

		return nl
	default:
		panic(fmt.Errorf("cannot find initial on type %s", tp))
	}
}

// sort alpha sorts given list.
//     `sort 5 1 4 3 2` will returns `1,2,3,4,5`.
func sortAlpha(list interface{}) []string {
	k := reflect.Indirect(reflect.ValueOf(list)).Kind()
	switch k {
	case reflect.Slice, reflect.Array:
		a := strslice(list)
		s := sort.StringSlice(a)
		s.Sort()
		return s
	}
	val, _ := generic.ConvertToString(list)
	return []string{val}
}

// Produce a new list with the reversed elements of the given list.
//     `reverse $myList`
//The above would generate the list `[5 4 3 2 1]`.
func reverse(v interface{}) []interface{} {
	tp := reflect.TypeOf(v).Kind()
	switch tp {
	case reflect.Slice, reflect.Array:
		l2 := reflect.ValueOf(v)

		l := l2.Len()
		// We do not sort in place because the incoming array should not be altered.
		nl := make([]interface{}, l)
		for i := 0; i < l; i++ {
			nl[l-i-1] = l2.Index(i).Interface()
		}

		return nl
	default:
		panic(fmt.Errorf("cannot find reverse on type %s", tp))
	}
}

// Generate a list with all of the duplicates removed.
//     `list 1 1 1 2 | uniq`
// The above would produce `[1 2]`
func uniq(list interface{}) []interface{} {
	tp := reflect.TypeOf(list).Kind()
	switch tp {
	case reflect.Slice, reflect.Array:
		l2 := reflect.ValueOf(list)

		l := l2.Len()
		dest := []interface{}{}
		var item interface{}
		for i := 0; i < l; i++ {
			item = l2.Index(i).Interface()
			if !inList(dest, item) {
				dest = append(dest, item)
			}
		}

		return dest
	default:
		panic(fmt.Errorf("cannot find uniq on type %s", tp))
	}
}

func inList(haystack []interface{}, needle interface{}) bool {
	for _, h := range haystack {
		if reflect.DeepEqual(needle, h) {
			return true
		}
	}
	return false
}

// The `without` function filters items out of a list.
//     `without $myList 3`
// The above would produce `[1 2 4 5]`
// Without can take more than one filter:
//     `without $myList 1 3 5`
// That would produce `[2 4]`
func without(list interface{}, omit ...interface{}) []interface{} {
	tp := reflect.TypeOf(list).Kind()
	switch tp {
	case reflect.Slice, reflect.Array:
		l2 := reflect.ValueOf(list)

		l := l2.Len()
		res := []interface{}{}
		var item interface{}
		for i := 0; i < l; i++ {
			item = l2.Index(i).Interface()
			if !inList(omit, item) {
				res = append(res, item)
			}
		}

		return res
	default:
		panic(fmt.Errorf("cannot find without on type %s", tp))
	}
}

// Test to see if a list has a particular element.
//     `has 4 $myList`
// The above would return `true`, while `has "hello" $myList` would return false.
func has(needle interface{}, haystack interface{}) bool {
	tp := reflect.TypeOf(haystack).Kind()
	switch tp {
	case reflect.Slice, reflect.Array:
		l2 := reflect.ValueOf(haystack)
		var item interface{}
		l := l2.Len()
		for i := 0; i < l; i++ {
			item = l2.Index(i).Interface()
			if reflect.DeepEqual(needle, item) {
				return true
			}
		}

		return false
	default:
		panic(fmt.Errorf("cannot find has on type %s", tp))
	}
}

// To get partial elements of a list, use `slice list [n] [m]`. It is
// equivalent of `list[n:m]`.
//- `slice $myList` returns `[1 2 3 4 5]`. It is same as `myList[:]`.
//- `slice $myList 3` returns `[4 5]`. It is same as `myList[3:]`.
//- `slice $myList 1 3` returns `[2 3]`. It is same as `myList[1:3]`.
//- `slice $myList 0 3` returns `[1 2 3]`. It is same as `myList[:3]`.
func slice(list interface{}, indices ...interface{}) interface{} {
	tp := reflect.TypeOf(list).Kind()
	switch tp {
	case reflect.Slice, reflect.Array:
		l2 := reflect.ValueOf(list)

		l := l2.Len()
		if l == 0 {
			return nil
		}

		var start, end int
		if len(indices) > 0 {
			start, _ = generic.ConvertToInt(indices[0])
		}
		if len(indices) < 2 {
			end = l
		} else {
			end, _ = generic.ConvertToInt(indices[1])
		}

		return l2.Slice(start, end).Interface()
	default:
		panic(fmt.Errorf("list should be type of slice or array but %s", tp))
	}
}

// Include will append value into the slice, if it has no same value
//     `$myList := include $myList 5`
func include(v interface{}, value string) interface{} {
	items := strslice(v)
	if value == "" {
		return items
	}

	for _, item := range items {
		if item == value {
			return v
		}
	}

	return append(items, value)
}

// Exclude will remove value from list
//     `$myList := remove $myList 5`
func exclude(v interface{}, value string) interface{} {
	items := strslice(v)
	for i, item := range items {
		if item == value {
			return append(items[:i], items[i+1:]...)
		}
	}

	return v
}

// Join concat all string values.
//     `join 1 2 3 4 5` will return `12345`.
func join(sep string, v interface{}) string {
	if v == nil {
		return ""
	}

	return strings.Join(strslice(v), sep)
}

// Functions for handle dictionaries.
// The key to a dictionary MUST BE A STRING. However, the value can be any type.
// Dictionaries are not immutable. The `set` and `unset` functions will
// modify the contents of a dictionary.

// Clone dictionary with add extra capacity
func cloneDict(d map[string]interface{}, extra int) map[string]interface{} {
	dict := make(map[string]interface{}, len(d)+extra)
	for k, v := range d {
		dict[k] = v
	}
	return dict
}

// aliveDict makes new dictionary, if it is nil.
func aliveDict(d map[string]interface{}) interface{} {
	if d != nil {
		return d
	}

	return make(map[string]interface{}, 32)
}

// Expand by clone original and append into ONLY new items.
// The following expand a original dictionary with three items:
//     $myDict := expand $original "name1" "value1" "name2" "value2" "name3" "value 3"
func expand(d map[string]interface{}, v ...interface{}) map[string]interface{} {
	dict := cloneDict(d, len(v))
	lenv := len(v)
	for i := 0; i < lenv; i += 2 {
		key, _ := generic.ConvertToString(v[i])
		if i+1 >= lenv {
			dict[key] = ""
			continue
		}
		if _, has := dict[key]; !has {
			dict[key] = v[i+1]
		}
	}
	return dict
}

// Extends dictionary by clone original and append into it ALL items.
// The following extends a original dictionary with three items:
//     $myDict := extends $original "name1" "value1" "name2" "value2" "name3" "value 3"
func extends(d map[string]interface{}, v ...interface{}) map[string]interface{} {
	dict := cloneDict(d, len(v))
	lenv := len(v)
	for i := 0; i < lenv; i += 2 {
		key, _ := generic.ConvertToString(v[i])
		if i+1 >= lenv {
			dict[key] = ""
			continue
		}
		dict[key] = v[i+1]
	}
	return dict
}

// Use `set` to add a new key/value pair to a dictionary.
//     $_ := set $myDict "name4" "value4"
// Note that `set` _returns the dictionary_ (a requirement of Go template functions),
// so you may need to trap the value as done above with the `$_` assignment.
func set(d map[string]interface{}, key string, value interface{}) map[string]interface{} {
	d[key] = value
	return d
}

// Given a map and a key, delete the key from the map.
//     $_ := unset $myDict "name4"
// As with `set`, this returns the dictionary.
// Note that if the key is not found, this operation will simply return. No error
// will be generated.
func unset(d map[string]interface{}, key string) map[string]interface{} {
	delete(d, key)
	return d
}

// The `hasKey` function returns `true` if the given dict contains the given key.
//     hasKey $myDict "name1"
// If the key is not found, this returns `false`.
func hasKey(d map[string]interface{}, key string) bool {
	_, ok := d[key]
	return ok
}

// The `pluck` function makes it possible to give one key and multiple maps, and
// get a list of all of the matches:
//     pluck "name1" $myDict $myOtherDict
// The above will return a `list` containing every found value (`[value1 otherValue1]`).
// If the give key is _not found_ in a map, that map will not have an item in the
// list (and the length of the returned list will be less than the number of dicts
// in the call to `pluck`.
// If the key is _found_ but the value is an empty value, that value will be
// inserted.
// A common idiom in Sprig templates is to uses `pluck... | first` to get the first
// matching key out of a collection of dictionaries.
func pluck(key string, d ...map[string]interface{}) []interface{} {
	res := make([]interface{}, 0, len(d))
	for _, dict := range d {
		if val, ok := dict[key]; ok {
			res = append(res, val)
		}
	}
	return res
}

// The `keys` function will return a `list` of all of the keys in one or more `dict`
// types. Since a dictionary is _unordered_, the keys will not be in a predictable order.
// They can be sorted with `sortAlpha`.
//    keys $myDict | sortAlpha
// When supplying multiple dictionaries, the keys will be concatenated. Use the `uniq`
// function along with `sortAlpha` to get a unqiue, sorted list of keys.
//    keys $myDict $myOtherDict | uniq | sortAlpha
func keys(dicts ...map[string]interface{}) []string {
	k := make([]string, 0, len(dicts))
	for _, dict := range dicts {
		for key := range dict {
			k = append(k, key)
		}
	}
	return k
}

// The `pick` function selects just the given keys out of a dictionary, creating a
// new `dict`.
//     $new := pick $myDict "name1" "name2"
// The above returns `{name1: value1, name2: value2}`
func pick(dict map[string]interface{}, keys ...string) map[string]interface{} {
	res := make(map[string]interface{}, len(keys))
	for _, k := range keys {
		if v, ok := dict[k]; ok {
			res[k] = v
		}
	}
	return res
}

// The `omit` function is similar to `pick`, except it returns a new `dict` with all
// the keys that _do not_ match the given keys.
//     $new := omit $myDict "name1" "name3"
// The above returns `{name2: value2}`
func omit(dict map[string]interface{}, keys ...string) map[string]interface{} {
	res := make(map[string]interface{}, len(dict))

	omit := make(map[string]bool, len(keys))
	for _, k := range keys {
		omit[k] = true
	}

	for k, v := range dict {
		if _, ok := omit[k]; !ok {
			res[k] = v
		}
	}
	return res
}

// Creating dictionaries is done by calling the `dict` function and passing it a
// list of pairs.
// The following creates a dictionary with three items:
//     $myDict := dict "name1" "value1" "name2" "value2" "name3" "value 3"
func dict(v ...interface{}) map[string]interface{} {
	dict := make(map[string]interface{}, len(v)/2)
	lenv := len(v)
	for i := 0; i < lenv; i += 2 {
		key, _ := generic.ConvertToString(v[i])
		if i+1 >= lenv {
			dict[key] = ""
			continue
		}
		dict[key] = v[i+1]
	}
	return dict
}

// The `values` function is similar to `keys`, except it returns a new `list` with
// all the values of the source `dict`.
//    $vals := values $myDict
// The above returns `list["value1", "value2", "value 3"]`. Note that the `values`
// function gives no guarantees about the result ordering-
func values(dict map[string]interface{}) []interface{} {
	values := make([]interface{}, 0, len(dict))
	for _, value := range dict {
		values = append(values, value)
	}

	return values
}

// Translate key into associated value of dictionary
func translate(dict map[string]interface{}, key interface{}) interface{} {
	if k, ok := key.(string); ok {
		if v, ok := dict[k]; ok {
			return v
		}
	}

	return ""
}

// Produce the function map.
//
// Use this to pass the functions into the template engine:
//
// 	tpl := template.New("foo").Funcs(sprig.FuncMap()))
//
func FuncMap() template.FuncMap {
	m := make(map[string]interface{}, len(genericMap))
	for k, v := range genericMap {
		m[k] = v
	}
	return m
}

var genericMap = map[string]interface{}{
	// basic arithmetic.
	"add":     add,
	"sub":     sub,
	"div":     div,
	"mod":     mod,
	"mul":     mul,
	"biggest": max,
	"max":     max,
	"min":     min,
	"ceil":    ceil,
	"floor":   floor,
	"round":   round,

	// Lists:
	"list":      list,
	"include":   include,
	"exclude":   exclude,
	"append":    push,
	"push":      push,
	"prepend":   prepend,
	"first":     first,
	"rest":      rest,
	"last":      last,
	"initial":   initial,
	"reverse":   reverse,
	"uniq":      uniq,
	"without":   without,
	"has":       has,
	"slice":     slice,
	"join":      join,
	"sortAlpha": sortAlpha,

	// Dictionaries:
	"dict":      dict,
	"expand":    expand,
	"extends":   extends,
	"set":       set,
	"unset":     unset,
	"hasKey":    hasKey,
	"pluck":     pluck,
	"keys":      keys,
	"pick":      pick,
	"omit":      omit,
	"values":    values,
	"translate": translate,
	"DICT":      aliveDict,
}

func strslice(v interface{}) []string {
	if v == nil {
		var res []string
		return res
	}

	switch v := v.(type) {
	case string:
		return []string{v}
	case []string:
		return v
	case []interface{}:
		l := len(v)
		b := make([]string, l)
		for i := 0; i < l; i++ {
			b[i], _ = generic.ConvertToString(v[i])
		}
		return b
	default:
		val := reflect.ValueOf(v)
		switch val.Kind() {
		case reflect.Array, reflect.Slice:
			l := val.Len()
			b := make([]string, l)
			for i := 0; i < l; i++ {
				b[i], _ = generic.ConvertToString(val.Index(i).Interface())
			}
			return b
		default:
			vv, _ := generic.ConvertToString(v)
			return []string{vv}
		}
	}
}
