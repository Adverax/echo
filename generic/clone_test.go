package generic

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCloneValue(t *testing.T) {
	type data struct {
		a string
		b string
		c []string
	}

	src := data{
		"works1",
		"works2",
		[]string{"a", "b"},
	}

	dst := CloneValue(src)
	assert.Equal(t, src, dst)

	src.c = append(src.c, "c")
	assert.NotEqual(t, src, dst)
}
