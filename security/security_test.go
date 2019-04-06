package security

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEngine_Filter(t *testing.T) {
	s := m.Filter("10,27&3 5hj", NUMBER)
	assert.Equal(t, "102735", s)
}

func TestEngine_CreateIdentifier(t *testing.T) {
	id := m.CreateIdentifier(NUMBER, 4)
	assert.Equal(t, 16, len(id))
}

func TestEngine_DecodeIdentifier(t *testing.T) {
	id := m.DecodeIdentifier("1234567812", 3)
	assert.Equal(t, "0012-3456-7812", id)
}

func TestEngine_EncodeIdentifier(t *testing.T) {
	id := m.EncodeIdentifier("0012-3456-7812", NUMBER, 3)
	assert.Equal(t, "001234567812", id)
}

func TestEngine_DecodeGuid(t *testing.T) {
	id := m.DecodeGuid(0x1111222233334444)
	assert.Equal(t, "1111-2222-3333-4444", id)
}

func TestEngine_EncodeGuid(t *testing.T) {
	id := m.EncodeGuid("1111-2222-3333-4444")
	var exp uint64 = 0x1111222233334444
	assert.Equal(t, exp, id)
}
