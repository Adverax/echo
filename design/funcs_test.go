package design

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestEmpty(t *testing.T) {
	tpl := `{{if empty 1}}1{{else}}0{{end}}`
	if err := runt(tpl, "0"); err != nil {
		t.Error(err)
	}

	tpl = `{{if empty 0}}1{{else}}0{{end}}`
	if err := runt(tpl, "1"); err != nil {
		t.Error(err)
	}
	tpl = `{{if empty ""}}1{{else}}0{{end}}`
	if err := runt(tpl, "1"); err != nil {
		t.Error(err)
	}
	tpl = `{{if empty 0.0}}1{{else}}0{{end}}`
	if err := runt(tpl, "1"); err != nil {
		t.Error(err)
	}
	tpl = `{{if empty false}}1{{else}}0{{end}}`
	if err := runt(tpl, "1"); err != nil {
		t.Error(err)
	}

	dict := map[string]interface{}{"top": map[string]interface{}{}}
	tpl = `{{if empty .top.NoSuchThing}}1{{else}}0{{end}}`
	if err := runtv(tpl, "1", dict); err != nil {
		t.Error(err)
	}
	tpl = `{{if empty .bottom.NoSuchThing}}1{{else}}0{{end}}`
	if err := runtv(tpl, "1", dict); err != nil {
		t.Error(err)
	}
}

func TestCoalesce(t *testing.T) {
	tests := map[string]string{
		`{{ coalesce 1 }}`:                            "1",
		`{{ coalesce "" 0 nil 2 }}`:                   "2",
		`{{ $two := 2 }}{{ coalesce "" 0 nil $two }}`: "2",
		`{{ $two := 2 }}{{ coalesce "" $two 0 0 0 }}`: "2",
		`{{ $two := 2 }}{{ coalesce "" $two 3 4 5 }}`: "2",
		`{{ coalesce }}`:                              "",
	}
	for tpl, expect := range tests {
		assert.NoError(t, runt(tpl, expect))
	}

	dict := map[string]interface{}{"top": map[string]interface{}{}}
	tpl := `{{ coalesce .top.NoSuchThing .bottom .bottom.dollar "airplane"}}`
	if err := runtv(tpl, "airplane", dict); err != nil {
		t.Error(err)
	}
}

func TestList(t *testing.T) {
	tpl := `{{$t := list 1 "a" "foo"}}{{index $t 2}}{{index $t 0 }}{{index $t 1}}`
	if err := runt(tpl, "foo1a"); err != nil {
		t.Error(err)
	}
}

func TestPush(t *testing.T) {
	// Named `append` in the function map
	tests := map[string]string{
		`{{ $t := list 1 2 3  }}{{ append $t 4 | len }}`:                    "4",
		`{{ $t := list 1 2 3 4  }}{{ append $t 5 | join "-" }}`:             "1-2-3-4-5",
		`{{ $t := list "foo" "bar" "baz"}}{{ append $t "qux" | join "-" }}`: "foo-bar-baz-qux",
	}
	for tpl, expect := range tests {
		assert.NoError(t, runt(tpl, expect))
	}
}

func TestPrepend(t *testing.T) {
	tests := map[string]string{
		`{{ $t := list 1 2 3  }}{{ prepend $t 0 | len }}`:                    "4",
		`{{ $t := list 1 2 3 4  }}{{ prepend $t 0 | join "-" }}`:             "0-1-2-3-4",
		`{{ $t := list "foo" "bar" "baz"}}{{ prepend $t "qux" | join "-" }}`: "qux-foo-bar-baz",
	}
	for tpl, expect := range tests {
		assert.NoError(t, runt(tpl, expect))
	}
}

func TestFirst(t *testing.T) {
	tests := map[string]string{
		`{{ list 1 2 3 | first }}`:       "1",
		`{{ list | first }}`:             "",
		`{{ list "foo" "bar" | first }}`: "foo",
	}
	for tpl, expect := range tests {
		assert.NoError(t, runt(tpl, expect))
	}
}
func TestLast(t *testing.T) {
	tests := map[string]string{
		`{{ list 1 2 3 | last }}`:      "3",
		`{{ list | last }}`:            "",
		`{{ list "foo" "bar"| last }}`: "bar",
	}
	for tpl, expect := range tests {
		assert.NoError(t, runt(tpl, expect))
	}
}

func TestInitial(t *testing.T) {
	tests := map[string]string{
		`{{ list 1 2 3 | initial | len }}`:       "2",
		`{{ list 1 2 3 | initial | last }}`:      "2",
		`{{ list 1 2 3 | initial | first }}`:     "1",
		`{{ list | initial }}`:                   "[]",
		`{{ list "foo" "bar" "baz" | initial }}`: "[foo bar]",
	}
	for tpl, expect := range tests {
		assert.NoError(t, runt(tpl, expect))
	}
}

func TestRest(t *testing.T) {
	tests := map[string]string{
		`{{ list 1 2 3 | rest | len }}`:       "2",
		`{{ list 1 2 3 | rest | last }}`:      "3",
		`{{ list 1 2 3 | rest | first }}`:     "2",
		`{{ list | rest }}`:                   "[]",
		`{{ list "foo" "bar" "baz" | rest }}`: "[bar baz]",
	}
	for tpl, expect := range tests {
		assert.NoError(t, runt(tpl, expect))
	}
}

func TestReverse(t *testing.T) {
	tests := map[string]string{
		`{{ list 1 2 3 | reverse | first }}`:        "3",
		`{{ list 1 2 3 | reverse | rest | first }}`: "2",
		`{{ list 1 2 3 | reverse | last }}`:         "1",
		`{{ list 1 2 3 4 | reverse }}`:              "[4 3 2 1]",
		`{{ list 1 | reverse }}`:                    "[1]",
		`{{ list | reverse }}`:                      "[]",
		`{{ list "foo" "bar" "baz" | reverse }}`:    "[baz bar foo]",
	}
	for tpl, expect := range tests {
		assert.NoError(t, runt(tpl, expect))
	}
}

func TestUniq(t *testing.T) {
	tests := map[string]string{
		`{{ list 1 2 3 4 | uniq }}`:                   `[1 2 3 4]`,
		`{{ list "a" "b" "c" "d" | uniq }}`:           `[a b c d]`,
		`{{ list 1 1 1 1 2 2 2 2 | uniq }}`:           `[1 2]`,
		`{{ list "foo" 1 1 1 1 "foo" "foo" | uniq }}`: `[foo 1]`,
		`{{ list | uniq }}`:                           `[]`,
		`{{ list "foo" "foo" "bar" | uniq }}`:         "[foo bar]",
	}
	for tpl, expect := range tests {
		assert.NoError(t, runt(tpl, expect))
	}
}

func TestWithout(t *testing.T) {
	tests := map[string]string{
		`{{ without (list 1 2 3 4) 1 }}`:               `[2 3 4]`,
		`{{ without (list "a" "b" "c" "d") "a" }}`:     `[b c d]`,
		`{{ without (list 1 1 1 1 2) 1 }}`:             `[2]`,
		`{{ without (list) 1 }}`:                       `[]`,
		`{{ without (list 1 2 3) }}`:                   `[1 2 3]`,
		`{{ without list }}`:                           `[]`,
		`{{ without (list "foo" "bar" "baz") "foo" }}`: "[bar baz]",
	}
	for tpl, expect := range tests {
		assert.NoError(t, runt(tpl, expect))
	}
}

func TestHas(t *testing.T) {
	tests := map[string]string{
		`{{ list 1 2 3 | has 1 }}`:                 `true`,
		`{{ list 1 2 3 | has 4 }}`:                 `false`,
		`{{ list "foo" "bar" "baz" | has "bar" }}`: `true`,
	}
	for tpl, expect := range tests {
		assert.NoError(t, runt(tpl, expect))
	}
}

func TestSlice(t *testing.T) {
	tests := map[string]string{
		`{{ slice (list 1 2 3) }}`:                 "[1 2 3]",
		`{{ slice (list 1 2 3) 0 1 }}`:             "[1]",
		`{{ slice (list 1 2 3) 1 3 }}`:             "[2 3]",
		`{{ slice (list 1 2 3) 1 }}`:               "[2 3]",
		`{{ slice (list "foo" "bar" "baz") 1 2 }}`: "[bar]",
	}
	for tpl, expect := range tests {
		assert.NoError(t, runt(tpl, expect))
	}
}

func TestDict(t *testing.T) {
	tpl := `{{$d := dict 1 2 "three" "four" 5}}{{range $k, $v := $d}}{{$k}}{{$v}}{{end}}`
	out, err := runRaw(tpl, nil)
	if err != nil {
		t.Error(err)
	}
	if len(out) != 12 {
		t.Errorf("Expected length 12, got %d", len(out))
	}
	// dict does not guarantee ordering because it is backed by a map.
	if !strings.Contains(out, "12") {
		t.Error("Expected grouping 12")
	}
	if !strings.Contains(out, "threefour") {
		t.Error("Expected grouping threefour")
	}
	if !strings.Contains(out, "5") {
		t.Error("Expected 5")
	}
	tpl = `{{$t := dict "I" "shot" "the" "albatross"}}{{$t.the}} {{$t.I}}`
	if err := runt(tpl, "albatross shot"); err != nil {
		t.Error(err)
	}
}

func TestExpand(t *testing.T) {
	tpl := `{{$d := expand (dict 1 2) 1 10 "three" "four" 5}}{{range $k, $v := $d}}{{$k}}{{$v}}{{end}}`
	out, err := runRaw(tpl, nil)
	if err != nil {
		t.Error(err)
	}
	if len(out) != 12 {
		t.Errorf("Expected length 12, got %d", len(out))
	}
	// dict does not guarantee ordering because it is backed by a map.
	if !strings.Contains(out, "12") {
		t.Error("Expected grouping 12")
	}
	if !strings.Contains(out, "threefour") {
		t.Error("Expected grouping threefour")
	}
	if !strings.Contains(out, "5") {
		t.Error("Expected 5")
	}
	tpl = `{{$t := dict "I" "shot" "the" "albatross"}}{{$t.the}} {{$t.I}}`
	if err := runt(tpl, "albatross shot"); err != nil {
		t.Error(err)
	}
}

func TestExtends(t *testing.T) {
	tpl := `{{$d := extends (dict 1 2) 1 9 "three" "four" 5}}{{range $k, $v := $d}}{{$k}}{{$v}}{{end}}`
	out, err := runRaw(tpl, nil)
	if err != nil {
		t.Error(err)
	}
	if len(out) != 12 {
		t.Errorf("Expected length 12, got %d", len(out))
	}
	// dict does not guarantee ordering because it is backed by a map.
	if !strings.Contains(out, "19") {
		t.Error("Expected grouping 19")
	}
	if !strings.Contains(out, "threefour") {
		t.Error("Expected grouping threefour")
	}
	if !strings.Contains(out, "5") {
		t.Error("Expected 5")
	}
	tpl = `{{$t := dict "I" "shot" "the" "albatross"}}{{$t.the}} {{$t.I}}`
	if err := runt(tpl, "albatross shot"); err != nil {
		t.Error(err)
	}
}

func TestUnset(t *testing.T) {
	tpl := `{{- $d := dict "one" 1 "two" 222222 -}}
	{{- $_ := unset $d "two" -}}
	{{- range $k, $v := $d}}{{$k}}{{$v}}{{- end -}}
	`

	expect := "one1"
	if err := runt(tpl, expect); err != nil {
		t.Error(err)
	}
}
func TestHasKey(t *testing.T) {
	tpl := `{{- $d := dict "one" 1 "two" 222222 -}}
	{{- if hasKey $d "one" -}}1{{- end -}}
	`

	expect := "1"
	if err := runt(tpl, expect); err != nil {
		t.Error(err)
	}
}

func TestPluck(t *testing.T) {
	tpl := `
	{{- $d := dict "one" 1 "two" 222222 -}}
	{{- $d2 := dict "one" 1 "two" 33333 -}}
	{{- $d3 := dict "one" 1 -}}
	{{- $d4 := dict "one" 1 "two" 4444 -}}
	{{- pluck "two" $d $d2 $d3 $d4 -}}
	`

	expect := "[222222 33333 4444]"
	if err := runt(tpl, expect); err != nil {
		t.Error(err)
	}
}

func TestKeys(t *testing.T) {
	tests := map[string]string{
		`{{ dict "foo" 1 "bar" 2 | keys | sortAlpha }}`: "[bar foo]",
		`{{ dict | keys }}`:                             "[]",
		`{{ keys (dict "foo" 1) (dict "bar" 2) (dict "bar" 3) | uniq | sortAlpha }}`: "[bar foo]",
	}
	for tpl, expect := range tests {
		if err := runt(tpl, expect); err != nil {
			t.Error(err)
		}
	}
}

func TestPick(t *testing.T) {
	tests := map[string]string{
		`{{- $d := dict "one" 1 "two" 222222 }}{{ pick $d "two" | len -}}`:               "1",
		`{{- $d := dict "one" 1 "two" 222222 }}{{ pick $d "two" -}}`:                     "map[two:222222]",
		`{{- $d := dict "one" 1 "two" 222222 }}{{ pick $d "one" "two" | len -}}`:         "2",
		`{{- $d := dict "one" 1 "two" 222222 }}{{ pick $d "one" "two" "three" | len -}}`: "2",
		`{{- $d := dict }}{{ pick $d "two" | len -}}`:                                    "0",
	}
	for tpl, expect := range tests {
		if err := runt(tpl, expect); err != nil {
			t.Error(err)
		}
	}
}
func TestOmit(t *testing.T) {
	tests := map[string]string{
		`{{- $d := dict "one" 1 "two" 222222 }}{{ omit $d "one" | len -}}`:         "1",
		`{{- $d := dict "one" 1 "two" 222222 }}{{ omit $d "one" -}}`:               "map[two:222222]",
		`{{- $d := dict "one" 1 "two" 222222 }}{{ omit $d "one" "two" | len -}}`:   "0",
		`{{- $d := dict "one" 1 "two" 222222 }}{{ omit $d "two" "three" | len -}}`: "1",
		`{{- $d := dict }}{{ omit $d "two" | len -}}`:                              "0",
	}
	for tpl, expect := range tests {
		if err := runt(tpl, expect); err != nil {
			t.Error(err)
		}
	}
}

func TestSet(t *testing.T) {
	tpl := `{{- $d := dict "one" 1 "two" 222222 -}}
	{{- $_ := set $d "two" 2 -}}
	{{- $_ := set $d "three" 3 -}}
	{{- if hasKey $d "one" -}}{{$d.one}}{{- end -}}
	{{- if hasKey $d "two" -}}{{$d.two}}{{- end -}}
	{{- if hasKey $d "three" -}}{{$d.three}}{{- end -}}
	`

	expect := "123"
	if err := runt(tpl, expect); err != nil {
		t.Error(err)
	}
}

func TestValues(t *testing.T) {
	tests := map[string]string{
		`{{- $d := dict "a" 1 "b" 2 }}{{ values $d | sortAlpha | join "," }}`:       "1,2",
		`{{- $d := dict "a" "first" "b" 2 }}{{ values $d | sortAlpha | join "," }}`: "2,first",
	}

	for tpl, expect := range tests {
		if err := runt(tpl, expect); err != nil {
			t.Error(err)
		}
	}
}
