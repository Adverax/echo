package design

import (
	"html/template"
	"math"

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
//     max 1 2 3 will return  3.
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
//     min 1 2 3` will return 1.
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
//     floor 123.9999` will return `123.0
func floor(a interface{}) float64 {
	aa, _ := generic.ConvertToFloat64(a)
	return math.Floor(aa)
}

// Returns the greatest float value greater than or equal to input value
//     ceil 123.001` will return `124.0
func ceil(a interface{}) float64 {
	aa, _ := generic.ConvertToFloat64(a)
	return math.Ceil(aa)
}

// Returns a float value with the remainder rounded to the given number to digits after the decimal point.
//     round 123.555555` will return `123.556
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

// Functions for handle dictionaries.
// The key to a dictionary MUST BE A STRING. However, the value can be any type.
// Dictionaries are not immutable. The `set` and `unset` functions will
// modify the contents of a dictionary.

// Clone dictionary with add extra capacity
func clone(d map[string]interface{}, extra int) map[string]interface{} {
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
	dict := clone(d, len(v))
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
	dict := clone(d, len(v))
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
