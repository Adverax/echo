package design

import (
	"bytes"
	"fmt"
	"github.com/adverax/echo"
	"github.com/stretchr/testify/assert"
	"html/template"
	"testing"
)

// runt runs a template and checks that the output exactly matches the expected string.
func runt(tpl, expect string) error {
	return runtv(tpl, expect, map[string]string{})
}

// runtv takes a template, and expected return, and values for substitution.
//
// It runs the template and verifies that the output is an exact match.
func runtv(tpl, expect string, vars interface{}) error {
	fmap := FuncMap()
	t := template.Must(template.New("test").Funcs(fmap).Parse(tpl))
	var b bytes.Buffer
	err := t.Execute(&b, vars)
	if err != nil {
		return err
	}
	if expect != b.String() {
		return fmt.Errorf("expected '%s', got '%s'", expect, b.String())
	}
	return nil
}

// runRaw runs a template with the given variables and returns the result.
func runRaw(tpl string, vars interface{}) (string, error) {
	fmap := FuncMap()
	t := template.Must(template.New("test").Funcs(fmap).Parse(tpl))
	var b bytes.Buffer
	err := t.Execute(&b, vars)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

func TestDesigner(t *testing.T) {
	d := NewDesigner(
		echo.New(),
		nil,
		"../_fixture/views",
		"../_fixture/views",
		"main.tmpl",
	)

	tpl := d.Compile("content.tmpl", "library.tmpl")

	var buf bytes.Buffer
	err := tpl.Execute(&buf, "Jack")
	assert.NoError(t, err)
	assert.Equal(t, "Header\n\nHello, Jack.\n\nElement.\n\n\nFooter", buf.String())
}
