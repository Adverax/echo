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

package data

import (
	"context"
	"sort"
	"strconv"
	"strings"
)

// Pair enumerator
type PairEnumerator interface {
	Enumerate(ctx context.Context, action PairEnumeratorFunc) error
}

type PairEnumeratorFunc func(key, value string) error

// Abstract data set
// Worked with literal representation keys and values.
type Set interface {
	PairEnumerator
	Has(key string) bool
	Length() int
}

// Map of Set by language code.
type Sets map[uint16]Set // Set by language

// Set provider
type SetProvider interface {
	DataSet(ctx context.Context) (Set, error)
}

// Simple pair of key and value
type pair struct {
	key string
	val string
}

type index []pair

func (index index) Len() int {
	return len(index)
}

func (index index) Less(i, j int) bool {
	return index[i].val < index[j].val
}

func (index index) Swap(i, j int) {
	index[i], index[j] = index[j], index[i]
}

// Set implementation
type set struct {
	items map[string]string
	index index
}

func (ds *set) Has(key string) bool {
	_, has := ds.items[key]
	return has
}

func (ds *set) Enumerate(
	ctx context.Context,
	action PairEnumeratorFunc,
) error {
	for _, pair := range ds.index {
		if pair.val != "" {
			err := action(pair.key, pair.val)
			if err != nil {
				return nil
			}
		}
	}

	return nil
}

func (ds *set) Length() int {
	return len(ds.index)
}

// Create new Set from map
func NewSet(
	items map[string]string, // Map of items
	sorted bool, // If you want sort items by alphabetically.
) Set {
	index := make(index, 0, len(items))

	for key, val := range items {
		index = append(index, pair{key, val})
	}

	if sorted {
		sort.Sort(index)
	}

	return &set{
		items: items,
		index: index,
	}
}

// Create new Set from list.
func NewSetFromList(
	items []string, // Slice of items
	sorted bool, // If you want sort items by alphabetically.
) Set {
	m := make(map[string]string, len(items))
	for index, val := range items {
		key := strconv.FormatInt(int64(index)+1, 10)
		m[key] = val
	}

	return NewSet(m, sorted)
}

var EmptySet = NewSet(make(map[string]string), false)

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
//  * MAP - dictionary required
//  * LIST - list required
//  * SORTED - items must be sorted
//  * DELIMITER - separator for maps (between key and value). Default - ":".
// Example: #! MAP SORTED DELIMITER ::
//
// LIST FORMAT
// The list has no features. One line - one element.
// The numbering of elements begins with one.
//
// DICTIONARY FORMAT
// Each element has a format: key: value or just a value.
// The numbering of elements starts from one (used only if the key is omitted). The next item is max + 1.
// Blank lines (or "_" lines) have code, but are not displayed.
// Empty characters at the beginning and end of the line are ignored.
func ParseSet(source string) Set {
	source = strings.TrimSpace(source)
	if source == "" {
		return EmptySet
	}

	lines := strings.Split(source, "\n")
	head := lines[0]
	if strings.HasPrefix(head, "#!") {
		lines = parseSetSkipSpace(lines[1:])
		list, sorted, delimiter := parseSetHeader(head[2:])
		if list {
			return parseSetList(lines, sorted)
		} else {
			return parseSetMap(lines, sorted, delimiter)
		}
	}

	lines = parseSetSkipSpace(lines)
	return parseSetList(lines, false)
}

func parseSetSkipSpace(lines []string) []string {
	for i, s := range lines {
		cols := strings.SplitN(s, "#", 2)
		lines[i] = strings.TrimSpace(cols[0])
	}
	return lines
}

func parseSetHeader(
	source string,
) (
	list bool,
	sorted bool,
	delimiter string,
) {
	sorted = false
	list = true
	delimiter = ":"
	columns := strings.SplitN(source, "#", 2)
	items := strings.Split(columns[0], " ")
	var index int
	for ; index < len(items); index++ {
		item := strings.ToLower(strings.TrimSpace(items[index]))
		switch item {
		case "map":
			list = false
		case "list":
			list = true
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
	return list, sorted, delimiter
}

func parseSetList(
	lines []string,
	sorted bool,
) Set {
	if len(lines) == 0 {
		return EmptySet
	}
	return NewSetFromList(lines, sorted)
}

func parseSetMap(
	lines []string,
	sorted bool,
	delimiter string,
) Set {
	if len(lines) == 0 {
		return EmptySet
	}

	enum := make(map[string]string, len(lines))
	var key int64 = 1
	var val string
	for _, line := range lines {
		line = strings.TrimSpace(line)
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
				// Skip item
				continue
			}
			key = k
			val = pair[1]
		}
		enum[strconv.FormatInt(key, 10)] = strings.TrimSpace(val)
		key++
	}

	return NewSet(enum, sorted)
}
