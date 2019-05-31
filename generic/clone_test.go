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

func TestCloneValueTo(t *testing.T) {
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

	var dst data
	CloneValueTo(&dst, src)
	assert.Equal(t, src, dst)

	src.a = "works3"
	assert.NotEqual(t, src, dst)
	src.a = "works1"

	src.c = append(src.c, "c")
	assert.NotEqual(t, src, dst)
}
