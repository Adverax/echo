package cache

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDummyCache(t *testing.T) {
	// Stub for test
	c := &DummyCache{}
	assert.NoError(t, c.Set("a", "b", time.Hour))
	var val interface{}
	assert.NoError(t, c.Get("a", &val))
	assert.Equal(t, nil, val)
}
