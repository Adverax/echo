package generic

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertAssign(t *testing.T) {
	{
		var dst int64
		err := ConvertAssign(&dst, int(10))
		assert.NoError(t, err)
		assert.Equal(t, int64(10), dst)
	}
	{
		var dst string
		err := ConvertAssign(&dst, int(10))
		assert.NoError(t, err)
		assert.Equal(t, "10", dst)
	}
	{
		type color int
		var dst color
		err := ConvertAssign(&dst, int(10))
		assert.NoError(t, err)
		assert.Equal(t, color(10), dst)
	}
}
