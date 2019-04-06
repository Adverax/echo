package memory

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	c := New(Options{})

	// Get/Set
	var vs, vs2 string
	err := c.Set("city1", "London", time.Hour)
	require.NoError(t, err)
	err = c.Set("city2", "Paris", time.Hour)
	require.NoError(t, err)
	err = c.Get("city1", &vs)
	require.NoError(t, err)
	assert.Equal(t, "London", vs)
	err = c.GetMulti(map[string]interface{}{
		"city2": &vs,
		"city1": &vs2,
	})
	assert.Equal(t, "Paris", vs)
	assert.Equal(t, "London", vs2)

	// IsExists/Clear
	has, err := c.IsExists("city1")
	require.NoError(t, err)
	assert.True(t, has)
	err = c.Clear()
	require.NoError(t, err)
	has, err = c.IsExists("city1")
	require.NoError(t, err)
	assert.False(t, has)

	// Increment/Decrement
	var vi int
	err = c.Set("visitors", int(0), time.Hour)
	require.NoError(t, err)

	err = c.Increase("visitors")
	require.NoError(t, err)
	err = c.Get("visitors", &vi)
	require.NoError(t, err)
	assert.Equal(t, int(1), vi)

	err = c.Increase("visitors")
	require.NoError(t, err)
	err = c.Get("visitors", &vi)
	require.NoError(t, err)
	assert.Equal(t, int(2), vi)

	err = c.Decrease("visitors")
	require.NoError(t, err)
	err = c.Get("visitors", &vi)
	require.NoError(t, err)
	assert.Equal(t, int(1), vi)
}
