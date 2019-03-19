package echo

/*
import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRenderer(t *testing.T) {
	r := NewRenderer(
		"_fixture/views",
		map[string]string{
			"main": "/main.tmpl",
		},
		nil,
	)

	const viewPath = "_fixture/views"
	tpl := r.ParseFiles(
		"@main",
		viewPath+"/content.tmpl",
		viewPath+"/library.tmpl",
	)
	var buf bytes.Buffer
	err := tpl.Execute(&buf, "Jack")
	assert.NoError(t, err)
	assert.Equal(t, "Header\n\nHello, Jack.\n\nElement.\n\n\nFooter", buf.String())
}
*/
