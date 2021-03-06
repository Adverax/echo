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
	stdContext "context"
	"github.com/adverax/echo/data"
	"github.com/adverax/echo/generic"
	"sort"
	"strconv"
	"strings"
)

type DATASET uint32

type DataSetManager interface {
	Find(ctx stdContext.Context, id uint32, language uint16) (DataSet, error)
	FindAll(ctx stdContext.Context, id uint32) (DataSets, error)
}

type DataSetConsumer func(key string, value string) error

// DataSet enumerator
type DataSetEnumerator interface {
	// Get items count
	Length(ctx Context) (int, error)
	// Enumerate all items
	Enumerate(ctx Context, consumer DataSetConsumer) error
}

// Abstract data set
// Works with literal representation keys and values.
type DataSet interface {
	Codec
	DataSetEnumerator
	DataSetProvider
}

// Map of DataSet by language code.
type DataSets map[uint16]DataSet // DataSet by language

func (datasets DataSets) Empty(ctx Context) (interface{}, error) {
	ds, err := datasets.DataSet(ctx)
	if err != nil {
		return nil, err
	}

	return ds.Empty(ctx)
}

func (datasets DataSets) DataSet(ctx Context) (DataSet, error) {
	lang := ctx.Locale().Language()
	if ds, ok := datasets[lang]; ok {
		return ds, nil
	}

	return nil, data.ErrNoMatch
}

func (datasets DataSets) Enumerate(ctx Context, action DataSetConsumer) error {
	ds, err := datasets.DataSet(ctx)
	if err != nil {
		return err
	}

	return ds.Enumerate(ctx, action)
}

func (datasets DataSets) Decode(ctx Context, value interface{}) (string, error) {
	ds, err := datasets.DataSet(ctx)
	if err != nil {
		return "", err
	}

	return ds.Decode(ctx, value)
}

func (datasets DataSets) Encode(ctx Context, value string) (interface{}, error) {
	ds, err := datasets.DataSet(ctx)
	if err != nil {
		return nil, err
	}

	return ds.Encode(ctx, value)
}

func (datasets DataSets) Length(ctx Context) (int, error) {
	ds, err := datasets.DataSet(ctx)
	if err != nil {
		return 0, err
	}

	return ds.Length(ctx)
}

// DataSet provider
type DataSetProvider interface {
	DataSet(ctx Context) (DataSet, error)
}

// Simple pair of key and value
type Pair struct {
	Key string
	Val string
}

type Pairs []Pair

func (ps Pairs) Len() int {
	return len(ps)
}

func (ps Pairs) Less(i, j int) bool {
	return ps[i].Val < ps[j].Val
}

func (ps Pairs) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

// RawPair supports transformation values from interfaces into the strings.
type RawPair struct {
	Key string
	Val interface{}
}

type RawPairs []RawPair

func (ps RawPairs) Compile(ctx Context) (Pairs, error) {
	res := make(Pairs, len(ps))
	for i, p := range ps {
		val, err := RenderString(ctx, p.Val)
		if err != nil {
			return nil, err
		}
		res[i].Val = val
		res[i].Key = p.Key
	}
	return res, nil
}

// DataSet implementation
type dataSet struct {
	encoders map[string]string
	decoders map[string]string
	index    Pairs
}

func (ds *dataSet) Empty(ctx Context) (interface{}, error) {
	return "", nil
}

func (ds *dataSet) Encode(ctx Context, value string) (interface{}, error) {
	if val, ok := ds.encoders[value]; ok {
		return val, nil
	}
	return "", data.ErrNoMatch
}

func (ds *dataSet) Decode(ctx Context, value interface{}) (string, error) {
	val, ok := generic.ConvertToString(value)
	if !ok {
		return "", data.ErrNoMatch
	}
	if v, ok := ds.decoders[val]; ok {
		return v, nil
	}
	return "", data.ErrNoMatch
}

func (ds *dataSet) Enumerate(
	ctx Context,
	action DataSetConsumer,
) error {
	for _, pair := range ds.index {
		if pair.Val != "" {
			err := action(pair.Key, pair.Val)
			if err != nil {
				return nil
			}
		}
	}

	return nil
}

func (ds *dataSet) Length(ctx Context) (int, error) {
	return len(ds.index), nil
}

func (ds *dataSet) DataSet(ctx Context) (DataSet, error) {
	return ds, nil
}

// Create new DataSet from map
func NewDataSet(
	items map[string]string, // Map of items
	sorted bool, // If you want sort items by alphabetically.
) DataSet {
	if items == nil {
		items = make(map[string]string)
	}

	encoders := make(map[string]string, len(items))
	index := make(Pairs, 0, len(items))

	for key, val := range items {
		encoders[val] = key
		index = append(index, Pair{key, val})
	}

	if sorted {
		sort.Sort(index)
	}

	return &dataSet{
		decoders: items,
		encoders: encoders,
		index:    index,
	}
}

// Create new DataSet from list.
func NewDataSetFromList(
	items []string, // Slice of items
	sorted bool, // If you want sort items by alphabetically.
) DataSet {
	m := make(map[string]string, len(items))
	for index, val := range items {
		key := strconv.FormatInt(int64(index)+1, 10)
		m[key] = val
	}

	return NewDataSet(m, sorted)
}

// Create new DataSet from sequence of key/val pairs.
func NewDataSetFromSequence(
	items Pairs,
	sorted bool,
) DataSet {
	if len(items) == 0 {
		return EmptySet
	}

	encoders := make(map[string]string, len(items))
	decoders := make(map[string]string, len(items))

	for _, item := range items {
		encoders[item.Val] = item.Key
		decoders[item.Key] = item.Val
	}

	if sorted {
		index := make(Pairs, len(items))
		copy(index, items)
		sort.Sort(index)
		items = index
	}

	return &dataSet{
		decoders: decoders,
		encoders: encoders,
		index:    items,
	}
}

var EmptySet = NewDataSet(make(map[string]string), false)

// Create DataSource from literal representation
// DATA SET FORMAT:
// The source may have a headline or not. In the absence of a headline, a list is implied.
// Each element is located on a separate line.
// Each line can end with a comment, that starts from "#".
//
// HEADLINE
// The title must begin with "#!".
// The header consists of attributes (with or without values), separated by spaces.
// Valid attributes:
//  * SORTED - items must be sorted
//  * DELIMITER - separator for maps (between key and value). Default - ":".
// Example: #! MAP SORTED DELIMITER ::
//
// ITEM FORMAT
// Each element has a format: key: value or just a value.
// The numbering of elements starts from one (used only if the key is omitted). The next item is max + 1.
// Blank lines (or "_" lines) have code, but are not displayed.
// Empty characters at the beginning and end of the line are ignored.
func ParseDataSet(source string) DataSet {
	source = strings.TrimSpace(source)
	if source == "" {
		return EmptySet
	}

	lines := strings.Split(source, "\n")
	head := lines[0]
	if strings.HasPrefix(head, "#!") {
		lines = parseDataSetSkipSpace(lines[1:])
		sorted, delimiter := parseDataSetHeader(head[2:])
		return parseDataSetSequence(lines, sorted, delimiter)
	}

	lines = parseDataSetSkipSpace(lines)
	return parseDataSetSequence(lines, false, ":")
}

func parseDataSetSkipSpace(lines []string) []string {
	for i, s := range lines {
		cols := strings.SplitN(s, "#", 2)
		lines[i] = strings.TrimSpace(cols[0])
	}
	return lines
}

func parseDataSetHeader(
	source string,
) (
	sorted bool,
	delimiter string,
) {
	sorted = false
	delimiter = ":"
	columns := strings.SplitN(source, "#", 2)
	items := strings.Split(columns[0], " ")
	var index int
	for ; index < len(items); index++ {
		item := strings.ToLower(strings.TrimSpace(items[index]))
		switch item {
		case "sorted":
			sorted = true
		case "delimiter":
			index++
			if index < len(items) {
				item = items[index]
				delimiter = item
			}
		}
	}
	return sorted, delimiter
}

func parseDataSetSequence(
	lines []string,
	sorted bool,
	delimiter string,
) DataSet {
	if len(lines) == 0 {
		return EmptySet
	}

	enum := make([]Pair, 0, len(lines))
	parseDataSetItems(
		lines,
		delimiter,
		func(key, val string) {
			enum = append(enum, Pair{
				Key: key,
				Val: val,
			})
		},
	)

	return NewDataSetFromSequence(enum, sorted)
}

func parseDataSetItems(
	lines []string,
	delimiter string,
	consume func(key, val string),
) {
	var key int64 = 1
	var val string
	for _, line := range lines {
		line = trimLine(line)
		pair := strings.SplitN(line, delimiter, 2)
		if len(pair) == 1 {
			val = pair[0]
			if val == "" || val == "_" {
				key++
				continue
			}
		} else {
			k, err := strconv.ParseInt(pair[0], 10, 32)
			if err != nil {
				// Is not numeric value.
				consume(pair[0], pair[1])
				continue
			}
			key = k
			val = pair[1]
		}
		consume(
			strconv.FormatInt(key, 10),
			strings.TrimSpace(val),
		)
		key++
	}
}

func IsPrimitiveDataSet(
	dataset DataSet,
) bool {
	_, ok := dataset.(*dataSet)
	return ok
}

// DataSetKeys returns list of all keys of dataset
func DataSetKeys(ctx Context, ds DataSet) ([]string, error) {
	length, err := ds.Length(ctx)
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0, length)
	err = ds.Enumerate(
		ctx,
		func(key string, value string) error {
			keys = append(keys, key)
			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	return keys, nil
}

// Trim line by deleting comments and white chars.
func trimLine(line string) string {
	pair := strings.SplitN(line, "#", 2)
	if len(pair) > 0 {
		return strings.TrimSpace(pair[0])
	}

	return strings.TrimSpace(line)
}
