package generic

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestConvertTo(t *testing.T) {
	type color int

	type Test struct {
		src interface{}
		dst interface{}
		ok  bool
	}

	tests := map[string]Test{
		// ConvertToString
		"Convert string to string must be success": {
			src: "Hello",
			dst: "Hello",
			ok:  true,
		},
		"Convert int to string must be success": {
			src: int(777),
			dst: "777",
			ok:  true,
		},
		"Convert int8 to string must be success": {
			src: int8(77),
			dst: "77",
			ok:  true,
		},
		"Convert int16 to string must be success": {
			src: int16(777),
			dst: "777",
			ok:  true,
		},
		"Convert int32 to string must be success": {
			src: int32(777),
			dst: "777",
			ok:  true,
		},
		"Convert int64 to string must be success": {
			src: int64(777),
			dst: "777",
			ok:  true,
		},
		"Convert uint to string must be success": {
			src: uint(777),
			dst: "777",
			ok:  true,
		},
		"Convert uint8 to string must be success": {
			src: uint(77),
			dst: "77",
			ok:  true,
		},
		"Convert uint16 to string must be success": {
			src: uint16(777),
			dst: "777",
			ok:  true,
		},
		"Convert uint32 to string must be success": {
			src: uint32(777),
			dst: "777",
			ok:  true,
		},
		"Convert uint64 to string must be success": {
			src: uint64(777),
			dst: "777",
			ok:  true,
		},
		"Convert float32 to string must be success": {
			src: float32(777),
			dst: "7.77000000e+02",
			ok:  true,
		},
		"Convert float64 to string must be success": {
			src: float64(777),
			dst: "7.77000000e+02",
			ok:  true,
		},
		"Convert `true` to string must be success": {
			src: true,
			dst: "1",
			ok:  true,
		},
		"Convert `false` to string must be success": {
			src: false,
			dst: "0",
			ok:  true,
		},
		"Convert `color` to string must be success": {
			src: color(10),
			dst: "10",
			ok:  true,
		},

		// ConvertToInt
		"Convert invalid string to int must be failure": {
			src: "Hello",
			dst: int(0),
			ok:  false,
		},
		"Convert valid string to int must be success": {
			src: "777",
			dst: int(777),
			ok:  true,
		},
		"Convert int to int must be success": {
			src: int(777),
			dst: int(777),
			ok:  true,
		},
		"Convert int8 to int must be success": {
			src: int8(77),
			dst: int(77),
			ok:  true,
		},
		"Convert int16 to int must be success": {
			src: int16(777),
			dst: int(777),
			ok:  true,
		},
		"Convert int32 to int must be success": {
			src: int32(777),
			dst: int(777),
			ok:  true,
		},
		"Convert int64 to int must be success": {
			src: int64(777),
			dst: int(777),
			ok:  true,
		},
		"Convert uint to int must be success": {
			src: uint(777),
			dst: int(777),
			ok:  true,
		},
		"Convert uint8 to int must be success": {
			src: uint8(77),
			dst: int(77),
			ok:  true,
		},
		"Convert uint16 to int must be success": {
			src: uint16(777),
			dst: int(777),
			ok:  true,
		},
		"Convert uint32 to int must be success": {
			src: uint32(777),
			dst: int(777),
			ok:  true,
		},
		"Convert uint64 to int must be success": {
			src: uint64(777),
			dst: int(777),
			ok:  true,
		},
		"Convert float32 to int must be success": {
			src: float32(777),
			dst: int(777),
			ok:  true,
		},
		"Convert float64 to int must be success": {
			src: float64(777),
			dst: int(777),
			ok:  true,
		},
		"Convert `true` to int must be success": {
			src: true,
			dst: int(1),
			ok:  true,
		},
		"Convert `false` to int must be success": {
			src: false,
			dst: int(0),
			ok:  true,
		},
		"Convert `color` to int must be success": {
			src: color(10),
			dst: int(10),
			ok:  true,
		},

		// ConvertToInt8
		"Convert invalid string to int8 must be failure": {
			src: "Hello",
			dst: int8(0),
			ok:  false,
		},
		"Convert valid string to int8 must be success": {
			src: "77",
			dst: int8(77),
			ok:  true,
		},
		"Convert int to int8 must be success": {
			src: int(77),
			dst: int8(77),
			ok:  true,
		},
		"Convert int8 to int8 must be success": {
			src: int8(77),
			dst: int8(77),
			ok:  true,
		},
		"Convert int16 to int8 must be success": {
			src: int16(77),
			dst: int8(77),
			ok:  true,
		},
		"Convert int32 to int8 must be success": {
			src: int32(77),
			dst: int8(77),
			ok:  true,
		},
		"Convert int64 to int8 must be success": {
			src: int64(77),
			dst: int8(77),
			ok:  true,
		},
		"Convert uint to int8 must be success": {
			src: uint(77),
			dst: int8(77),
			ok:  true,
		},
		"Convert uint8 to int8 must be success": {
			src: uint8(77),
			dst: int64(77),
			ok:  true,
		},
		"Convert uint16 to int8 must be success": {
			src: uint16(77),
			dst: int8(77),
			ok:  true,
		},
		"Convert uint32 to int8 must be success": {
			src: uint32(77),
			dst: int8(77),
			ok:  true,
		},
		"Convert uint64 to int8 must be success": {
			src: uint64(77),
			dst: int8(77),
			ok:  true,
		},
		"Convert float32 to int8 must be success": {
			src: float32(77),
			dst: int8(77),
			ok:  true,
		},
		"Convert float64 to int8 must be success": {
			src: float64(77),
			dst: int8(77),
			ok:  true,
		},
		"Convert `true` to int8 must be success": {
			src: true,
			dst: int8(1),
			ok:  true,
		},
		"Convert `false` to int8 must be success": {
			src: false,
			dst: int8(0),
			ok:  true,
		},
		"Convert `color` to int8 must be success": {
			src: color(10),
			dst: int8(10),
			ok:  true,
		},

		// ConvertToInt16
		"Convert invalid string to int16 must be failure": {
			src: "Hello",
			dst: int16(0),
			ok:  false,
		},
		"Convert valid string to int16 must be success": {
			src: "777",
			dst: int16(777),
			ok:  true,
		},
		"Convert int to int16 must be success": {
			src: int(777),
			dst: int16(777),
			ok:  true,
		},
		"Convert int8 to int16 must be success": {
			src: int8(77),
			dst: int16(77),
			ok:  true,
		},
		"Convert int16 to int16 must be success": {
			src: int16(777),
			dst: int16(777),
			ok:  true,
		},
		"Convert int32 to int16 must be success": {
			src: int32(777),
			dst: int16(777),
			ok:  true,
		},
		"Convert int64 to int16 must be success": {
			src: int64(777),
			dst: int16(777),
			ok:  true,
		},
		"Convert uint to int16 must be success": {
			src: uint(777),
			dst: int16(777),
			ok:  true,
		},
		"Convert uint8 to int16 must be success": {
			src: uint8(77),
			dst: int16(77),
			ok:  true,
		},
		"Convert uint16 to int16 must be success": {
			src: uint16(777),
			dst: int16(777),
			ok:  true,
		},
		"Convert uint32 to int16 must be success": {
			src: uint32(777),
			dst: int16(777),
			ok:  true,
		},
		"Convert uint64 to int16 must be success": {
			src: uint64(777),
			dst: int16(777),
			ok:  true,
		},
		"Convert float32 to int16 must be success": {
			src: float32(777),
			dst: int16(777),
			ok:  true,
		},
		"Convert float64 to int16 must be success": {
			src: float64(777),
			dst: int16(777),
			ok:  true,
		},
		"Convert `true` to int16 must be success": {
			src: true,
			dst: int16(1),
			ok:  true,
		},
		"Convert `false` to int16 must be success": {
			src: false,
			dst: int16(0),
			ok:  true,
		},
		"Convert `color` to int16 must be success": {
			src: color(10),
			dst: int16(10),
			ok:  true,
		},

		// ConvertToInt32
		"Convert invalid string to int32 must be failure": {
			src: "Hello",
			dst: int32(0),
			ok:  false,
		},
		"Convert valid string to int32 must be success": {
			src: "777",
			dst: int32(777),
			ok:  true,
		},
		"Convert int to int32 must be success": {
			src: int(777),
			dst: int32(777),
			ok:  true,
		},
		"Convert int8 to int32 must be success": {
			src: int8(77),
			dst: int32(77),
			ok:  true,
		},
		"Convert int16 to int32 must be success": {
			src: int16(777),
			dst: int32(777),
			ok:  true,
		},
		"Convert int32 to int32 must be success": {
			src: int32(777),
			dst: int32(777),
			ok:  true,
		},
		"Convert int64 to int32 must be success": {
			src: int64(777),
			dst: int32(777),
			ok:  true,
		},
		"Convert uint to int32 must be success": {
			src: uint(777),
			dst: int32(777),
			ok:  true,
		},
		"Convert uint8 to int32 must be success": {
			src: uint8(77),
			dst: int32(77),
			ok:  true,
		},
		"Convert uint16 to int32 must be success": {
			src: uint16(777),
			dst: int32(777),
			ok:  true,
		},
		"Convert uint32 to int32 must be success": {
			src: uint32(777),
			dst: int32(777),
			ok:  true,
		},
		"Convert uint64 to int32 must be success": {
			src: uint64(777),
			dst: int32(777),
			ok:  true,
		},
		"Convert float32 to int32 must be success": {
			src: float32(777),
			dst: int32(777),
			ok:  true,
		},
		"Convert float64 to int32 must be success": {
			src: float64(777),
			dst: int32(777),
			ok:  true,
		},
		"Convert `true` to int32 must be success": {
			src: true,
			dst: int32(1),
			ok:  true,
		},
		"Convert `false` to int32 must be success": {
			src: false,
			dst: int32(0),
			ok:  true,
		},
		"Convert `color` to int32 must be success": {
			src: color(10),
			dst: int32(10),
			ok:  true,
		},

		// ConvertToInt64
		"Convert invalid string to int64 must be failure": {
			src: "Hello",
			dst: int64(0),
			ok:  false,
		},
		"Convert valid string to int64 must be success": {
			src: "777",
			dst: int64(777),
			ok:  true,
		},
		"Convert int to int64 must be success": {
			src: int(777),
			dst: int64(777),
			ok:  true,
		},
		"Convert int8 to int64 must be success": {
			src: int8(77),
			dst: int64(77),
			ok:  true,
		},
		"Convert int16 to int64 must be success": {
			src: int16(777),
			dst: int64(777),
			ok:  true,
		},
		"Convert int32 to int64 must be success": {
			src: int32(777),
			dst: int64(777),
			ok:  true,
		},
		"Convert int64 to int64 must be success": {
			src: int64(777),
			dst: int64(777),
			ok:  true,
		},
		"Convert uint to int64 must be success": {
			src: uint(777),
			dst: int64(777),
			ok:  true,
		},
		"Convert uint8 to int64 must be success": {
			src: uint8(77),
			dst: int64(77),
			ok:  true,
		},
		"Convert uint16 to int64 must be success": {
			src: uint16(777),
			dst: int64(777),
			ok:  true,
		},
		"Convert uint32 to int64 must be success": {
			src: uint32(777),
			dst: int64(777),
			ok:  true,
		},
		"Convert uint64 to int64 must be success": {
			src: uint64(777),
			dst: int64(777),
			ok:  true,
		},
		"Convert float32 to int64 must be success": {
			src: float32(777),
			dst: int64(777),
			ok:  true,
		},
		"Convert float64 to int64 must be success": {
			src: float64(777),
			dst: int64(777),
			ok:  true,
		},
		"Convert `true` to int64 must be success": {
			src: true,
			dst: int64(1),
			ok:  true,
		},
		"Convert `false` to int64 must be success": {
			src: false,
			dst: int64(0),
			ok:  true,
		},
		"Convert `color` to int64 must be success": {
			src: color(10),
			dst: int64(10),
			ok:  true,
		},

		// ConvertToUint
		"Convert invalid string to uint must be failure": {
			src: "Hello",
			dst: uint(0),
			ok:  false,
		},
		"Convert valid string to uint must be success": {
			src: "777",
			dst: uint(777),
			ok:  true,
		},
		"Convert int to uint must be success": {
			src: int(777),
			dst: uint(777),
			ok:  true,
		},
		"Convert int8 to uint must be success": {
			src: int8(77),
			dst: uint(77),
			ok:  true,
		},
		"Convert int16 to uint must be success": {
			src: int16(777),
			dst: uint(777),
			ok:  true,
		},
		"Convert int32 to uint must be success": {
			src: int32(777),
			dst: uint(777),
			ok:  true,
		},
		"Convert int64 to uint must be success": {
			src: int64(777),
			dst: uint(777),
			ok:  true,
		},
		"Convert uint to uint must be success": {
			src: uint(777),
			dst: uint(777),
			ok:  true,
		},
		"Convert uint8 to uint must be success": {
			src: uint8(77),
			dst: uint(77),
			ok:  true,
		},
		"Convert uint16 to uint must be success": {
			src: uint16(777),
			dst: uint(777),
			ok:  true,
		},
		"Convert uint32 to uint must be success": {
			src: uint32(777),
			dst: uint(777),
			ok:  true,
		},
		"Convert uint64 to uint must be success": {
			src: uint64(777),
			dst: uint(777),
			ok:  true,
		},
		"Convert float32 to uint must be success": {
			src: float32(777),
			dst: uint(777),
			ok:  true,
		},
		"Convert float64 to uint must be success": {
			src: float64(777),
			dst: uint(777),
			ok:  true,
		},
		"Convert `true` to uint must be success": {
			src: true,
			dst: uint(1),
			ok:  true,
		},
		"Convert `false` to uint must be success": {
			src: false,
			dst: uint(0),
			ok:  true,
		},
		"Convert `color` to uint must be success": {
			src: color(10),
			dst: uint(10),
			ok:  true,
		},

		// ConvertToUint8
		"Convert invalid string to uint8 must be failure": {
			src: "Hello",
			dst: uint8(0),
			ok:  false,
		},
		"Convert valid string to uint8 must be success": {
			src: "77",
			dst: uint8(77),
			ok:  true,
		},
		"Convert int to uint8 must be success": {
			src: int(77),
			dst: uint8(77),
			ok:  true,
		},
		"Convert int8 to uint8 must be success": {
			src: int8(77),
			dst: uint8(77),
			ok:  true,
		},
		"Convert int16 to uint8 must be success": {
			src: int16(77),
			dst: uint8(77),
			ok:  true,
		},
		"Convert int32 to uint8 must be success": {
			src: int32(77),
			dst: uint8(77),
			ok:  true,
		},
		"Convert int64 to uint8 must be success": {
			src: int64(77),
			dst: uint8(77),
			ok:  true,
		},
		"Convert uint to uint8 must be success": {
			src: uint(77),
			dst: uint8(77),
			ok:  true,
		},
		"Convert uint8 to uint8 must be success": {
			src: uint8(77),
			dst: uint8(77),
			ok:  true,
		},
		"Convert uint16 to uint8 must be success": {
			src: uint16(77),
			dst: uint8(77),
			ok:  true,
		},
		"Convert uint32 to uint8 must be success": {
			src: uint32(77),
			dst: uint8(77),
			ok:  true,
		},
		"Convert uint64 to uint8 must be success": {
			src: uint64(77),
			dst: uint8(77),
			ok:  true,
		},
		"Convert float32 to uint8 must be success": {
			src: float32(77),
			dst: uint8(77),
			ok:  true,
		},
		"Convert float64 to uint8 must be success": {
			src: float64(77),
			dst: uint8(77),
			ok:  true,
		},
		"Convert `true` to uint8 must be success": {
			src: true,
			dst: uint8(1),
			ok:  true,
		},
		"Convert `false` to uint8 must be success": {
			src: false,
			dst: uint8(0),
			ok:  true,
		},
		"Convert `color` to uint8 must be success": {
			src: color(10),
			dst: uint8(10),
			ok:  true,
		},

		// ConvertToUint16
		"Convert invalid string to uint16 must be failure": {
			src: "Hello",
			dst: uint16(0),
			ok:  false,
		},
		"Convert valid string to uint16 must be success": {
			src: "777",
			dst: uint16(777),
			ok:  true,
		},
		"Convert int to uint16 must be success": {
			src: int(777),
			dst: uint16(777),
			ok:  true,
		},
		"Convert int8 to uint16 must be success": {
			src: int8(77),
			dst: uint16(77),
			ok:  true,
		},
		"Convert int16 to uint16 must be success": {
			src: int16(777),
			dst: uint16(777),
			ok:  true,
		},
		"Convert int32 to uint16 must be success": {
			src: int32(777),
			dst: uint16(777),
			ok:  true,
		},
		"Convert int64 to uint16 must be success": {
			src: int64(777),
			dst: uint16(777),
			ok:  true,
		},
		"Convert uint to uint16 must be success": {
			src: uint(777),
			dst: uint16(777),
			ok:  true,
		},
		"Convert uint8 to uint16 must be success": {
			src: uint8(77),
			dst: uint16(77),
			ok:  true,
		},
		"Convert uint16 to uint16 must be success": {
			src: uint16(777),
			dst: uint16(777),
			ok:  true,
		},
		"Convert uint32 to uint16 must be success": {
			src: uint32(777),
			dst: uint16(777),
			ok:  true,
		},
		"Convert uint64 to uint16 must be success": {
			src: uint64(777),
			dst: uint16(777),
			ok:  true,
		},
		"Convert float32 to uint16 must be success": {
			src: float32(777),
			dst: uint16(777),
			ok:  true,
		},
		"Convert float64 to uint16 must be success": {
			src: float64(777),
			dst: uint16(777),
			ok:  true,
		},
		"Convert `true` to uint16 must be success": {
			src: true,
			dst: uint16(1),
			ok:  true,
		},
		"Convert `false` to uint16 must be success": {
			src: false,
			dst: uint16(0),
			ok:  true,
		},
		"Convert `color` to uint16 must be success": {
			src: color(10),
			dst: uint16(10),
			ok:  true,
		},

		// ConvertToUint32
		"Convert invalid string to uint32 must be failure": {
			src: "Hello",
			dst: uint32(0),
			ok:  false,
		},
		"Convert valid string to uint32 must be success": {
			src: "777",
			dst: uint32(777),
			ok:  true,
		},
		"Convert int to uint32 must be success": {
			src: int(777),
			dst: uint32(777),
			ok:  true,
		},
		"Convert int8 to uint32 must be success": {
			src: int8(77),
			dst: uint32(77),
			ok:  true,
		},
		"Convert int16 to uint32 must be success": {
			src: int16(777),
			dst: uint32(777),
			ok:  true,
		},
		"Convert int32 to uint32 must be success": {
			src: int32(777),
			dst: uint32(777),
			ok:  true,
		},
		"Convert int64 to uint32 must be success": {
			src: int64(777),
			dst: uint32(777),
			ok:  true,
		},
		"Convert uint to uint32 must be success": {
			src: uint(777),
			dst: uint32(777),
			ok:  true,
		},
		"Convert uint8 to uint32 must be success": {
			src: uint8(77),
			dst: uint32(77),
			ok:  true,
		},
		"Convert uint16 to uint32 must be success": {
			src: uint16(777),
			dst: uint32(777),
			ok:  true,
		},
		"Convert uint32 to uint32 must be success": {
			src: uint32(777),
			dst: uint32(777),
			ok:  true,
		},
		"Convert uint64 to uint32 must be success": {
			src: uint64(777),
			dst: uint32(777),
			ok:  true,
		},
		"Convert float32 to uint32 must be success": {
			src: float32(777),
			dst: uint32(777),
			ok:  true,
		},
		"Convert float64 to uint32 must be success": {
			src: float64(777),
			dst: uint32(777),
			ok:  true,
		},
		"Convert `true` to uint32 must be success": {
			src: true,
			dst: uint32(1),
			ok:  true,
		},
		"Convert `false` to uint32 must be success": {
			src: false,
			dst: uint32(0),
			ok:  true,
		},
		"Convert `color` to uint32 must be success": {
			src: color(10),
			dst: uint32(10),
			ok:  true,
		},

		// ConvertToUint64
		"Convert invalid string to uint64 must be failure": {
			src: "Hello",
			dst: uint64(0),
			ok:  false,
		},
		"Convert valid string to uint64 must be success": {
			src: "777",
			dst: uint64(777),
			ok:  true,
		},
		"Convert int to uint64 must be success": {
			src: int(777),
			dst: uint64(777),
			ok:  true,
		},
		"Convert int8 to uint64 must be success": {
			src: int8(77),
			dst: uint64(77),
			ok:  true,
		},
		"Convert int16 to uint64 must be success": {
			src: int16(777),
			dst: uint64(777),
			ok:  true,
		},
		"Convert int32 to uint64 must be success": {
			src: int32(777),
			dst: uint64(777),
			ok:  true,
		},
		"Convert int64 to uint64 must be success": {
			src: int64(777),
			dst: uint64(777),
			ok:  true,
		},
		"Convert uint to uint64 must be success": {
			src: uint(777),
			dst: uint64(777),
			ok:  true,
		},
		"Convert uint8 to uint64 must be success": {
			src: uint8(77),
			dst: uint64(77),
			ok:  true,
		},
		"Convert uint16 to uint64 must be success": {
			src: uint16(777),
			dst: uint64(777),
			ok:  true,
		},
		"Convert uint32 to uint64 must be success": {
			src: uint32(777),
			dst: uint64(777),
			ok:  true,
		},
		"Convert uint64 to uint64 must be success": {
			src: uint64(777),
			dst: uint64(777),
			ok:  true,
		},
		"Convert float32 to uint64 must be success": {
			src: float32(777),
			dst: uint64(777),
			ok:  true,
		},
		"Convert float64 to uint64 must be success": {
			src: float64(777),
			dst: uint64(777),
			ok:  true,
		},
		"Convert `true` to uint64 must be success": {
			src: true,
			dst: uint64(1),
			ok:  true,
		},
		"Convert `false` to uint64 must be success": {
			src: false,
			dst: uint64(0),
			ok:  true,
		},
		"Convert `color` to uint64 must be success": {
			src: color(10),
			dst: uint64(10),
			ok:  true,
		},

		// ConvertToFloat32
		"Convert invalid string to float32 must be failure": {
			src: "Hello",
			dst: float32(0),
			ok:  false,
		},
		"Convert valid string to float32 must be success": {
			src: "777",
			dst: float32(777),
			ok:  true,
		},
		"Convert int to float32 must be success": {
			src: int(777),
			dst: float32(777),
			ok:  true,
		},
		"Convert int8 to float32 must be success": {
			src: int8(77),
			dst: float32(77),
			ok:  true,
		},
		"Convert int16 to float32 must be success": {
			src: int16(777),
			dst: float32(777),
			ok:  true,
		},
		"Convert int32 to float32 must be success": {
			src: int32(777),
			dst: float32(777),
			ok:  true,
		},
		"Convert int64 to float32 must be success": {
			src: int64(777),
			dst: float32(777),
			ok:  true,
		},
		"Convert uint to float32 must be success": {
			src: uint(777),
			dst: float32(777),
			ok:  true,
		},
		"Convert uint8 to float32 must be success": {
			src: uint8(77),
			dst: float32(77),
			ok:  true,
		},
		"Convert uint16 to float32 must be success": {
			src: uint16(777),
			dst: float32(777),
			ok:  true,
		},
		"Convert uint32 to float32 must be success": {
			src: uint32(777),
			dst: float32(777),
			ok:  true,
		},
		"Convert uint64 to float32 must be success": {
			src: uint64(777),
			dst: float32(777),
			ok:  true,
		},
		"Convert float32 to float32 must be success": {
			src: float32(777),
			dst: float32(777),
			ok:  true,
		},
		"Convert float64 to float32 must be success": {
			src: float64(777),
			dst: float32(777),
			ok:  true,
		},
		"Convert `true` to float32 must be success": {
			src: true,
			dst: float32(1),
			ok:  true,
		},
		"Convert `false` to float32 must be success": {
			src: false,
			dst: float32(0),
			ok:  true,
		},
		"Convert `color` to float32 must be success": {
			src: color(10),
			dst: float32(10),
			ok:  true,
		},

		// ConvertToFloat64
		"Convert invalid string to float64 must be failure": {
			src: "Hello",
			dst: float64(0),
			ok:  false,
		},
		"Convert valid string to float64 must be success": {
			src: "777",
			dst: float64(777),
			ok:  true,
		},
		"Convert int to float64 must be success": {
			src: int(777),
			dst: float64(777),
			ok:  true,
		},
		"Convert int8 to float64 must be success": {
			src: int8(77),
			dst: float64(77),
			ok:  true,
		},
		"Convert int16 to float64 must be success": {
			src: int16(777),
			dst: float64(777),
			ok:  true,
		},
		"Convert int32 to float64 must be success": {
			src: int32(777),
			dst: float64(777),
			ok:  true,
		},
		"Convert int64 to float64 must be success": {
			src: int64(777),
			dst: float64(777),
			ok:  true,
		},
		"Convert uint to float64 must be success": {
			src: uint(777),
			dst: float64(777),
			ok:  true,
		},
		"Convert uint8 to float64 must be success": {
			src: uint8(77),
			dst: float64(77),
			ok:  true,
		},
		"Convert uint16 to float64 must be success": {
			src: uint16(777),
			dst: float64(777),
			ok:  true,
		},
		"Convert uint32 to float64 must be success": {
			src: uint32(777),
			dst: float64(777),
			ok:  true,
		},
		"Convert uint64 to float64 must be success": {
			src: uint64(777),
			dst: float64(777),
			ok:  true,
		},
		"Convert float32 to float64 must be success": {
			src: float32(777),
			dst: float64(777),
			ok:  true,
		},
		"Convert float64 to float64 must be success": {
			src: float64(777),
			dst: float64(777),
			ok:  true,
		},
		"Convert `true` to float64 must be success": {
			src: true,
			dst: float64(1),
			ok:  true,
		},
		"Convert `false` to float64 must be success": {
			src: false,
			dst: float64(0),
			ok:  true,
		},
		"Convert `color` to float64 must be success": {
			src: color(10),
			dst: float64(10),
			ok:  true,
		},

		// ConvertToBoolean
		"Convert string to bool must be success": {
			src: "777",
			dst: true,
			ok:  true,
		},
		"Convert int to bool must be success": {
			src: int(777),
			dst: true,
			ok:  true,
		},
		"Convert int8 to bool must be success": {
			src: int8(77),
			dst: true,
			ok:  true,
		},
		"Convert int16 to bool must be success": {
			src: int16(777),
			dst: true,
			ok:  true,
		},
		"Convert int32 to bool must be success": {
			src: int32(777),
			dst: true,
			ok:  true,
		},
		"Convert int64 to bool must be success": {
			src: int64(777),
			dst: true,
			ok:  true,
		},
		"Convert uint to bool must be success": {
			src: uint(777),
			dst: true,
			ok:  true,
		},
		"Convert uint8 to bool must be success": {
			src: uint8(77),
			dst: true,
			ok:  true,
		},
		"Convert uint16 to bool must be success": {
			src: uint16(777),
			dst: true,
			ok:  true,
		},
		"Convert uint32 to bool must be success": {
			src: uint32(777),
			dst: true,
			ok:  true,
		},
		"Convert uint64 to bool must be success": {
			src: uint64(777),
			dst: true,
			ok:  true,
		},
		"Convert float32 to bool must be success": {
			src: float32(777),
			dst: true,
			ok:  true,
		},
		"Convert float64 to bool must be success": {
			src: float64(777),
			dst: true,
			ok:  true,
		},
		"Convert `true` to bool must be success": {
			src: true,
			dst: true,
			ok:  true,
		},
		"Convert `false` to bool must be success": {
			src: false,
			dst: false,
			ok:  true,
		},
		"Convert `color` to bool must be success": {
			src: color(10),
			dst: true,
			ok:  true,
		},
	}

	fs := map[string]interface{}{
		"int":     ConvertToInt,
		"int8":    ConvertToInt8,
		"int16":   ConvertToInt16,
		"int32":   ConvertToInt32,
		"int64":   ConvertToInt64,
		"uint":    ConvertToUint,
		"uint8":   ConvertToUint8,
		"uint16":  ConvertToUint16,
		"uint32":  ConvertToUint32,
		"uint64":  ConvertToUint64,
		"float32": ConvertToFloat32,
		"float64": ConvertToFloat64,
		"bool":    ConvertToBoolean,
		"string":  ConvertToString,
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			fname := reflect.TypeOf(test.dst).Name()
			f := reflect.ValueOf(fs[fname])
			in := []reflect.Value{reflect.ValueOf(test.src)}
			out := f.Call(in)
			dst, ok := out[0].Interface(), out[1].Bool()
			if assert.Equal(t, test.ok, ok) {
				assert.Equal(t, test.dst, dst)
			}
		})
	}
}
