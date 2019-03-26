// Copyright 2019 Adverax. All Rights Reserved.
// This file is part of project
//
//      http://github.com/adverax/echo
//
// Licensed under the MIT (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://github.com/adverax/echo/blob/master/LICENSE
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package generic

func IsEmpty(value interface{}) bool {
	switch val := value.(type) {
	case uint:
		return val == 0
	case uint8:
		return val == 0
	case uint16:
		return val == 0
	case uint32:
		return val == 0
	case uint64:
		return val == 0
	case int:
		return val == 0
	case int8:
		return val == 0
	case int16:
		return val == 0
	case int32:
		return val == 0
	case int64:
		return val == 0
	case float32:
		return val == 0
	case float64:
		return val == 0
	case string:
		return val == ""
	case bool:
		return val == false
	case []string:
		return len(val) == 0
	default:
		return true
	}
}

func CoalesceUint(values ...uint) uint {
	for _, val := range values {
		if val != 0 {
			return val
		}
	}
	return 0
}

func CoalesceUint8(values ...uint8) uint8 {
	for _, val := range values {
		if val != 0 {
			return val
		}
	}
	return 0
}

func CoalesceUint16(values ...uint16) uint16 {
	for _, val := range values {
		if val != 0 {
			return val
		}
	}
	return 0
}

func CoalesceUint32(values ...uint32) uint32 {
	for _, val := range values {
		if val != 0 {
			return val
		}
	}
	return 0
}

func CoalesceUint64(values ...uint64) uint64 {
	for _, val := range values {
		if val != 0 {
			return val
		}
	}
	return 0
}

func CoalesceInt(values ...int) int {
	for _, val := range values {
		if val != 0 {
			return val
		}
	}
	return 0
}

func CoalesceInt8(values ...int8) int8 {
	for _, val := range values {
		if val != 0 {
			return val
		}
	}
	return 0
}

func CoalesceInt16(values ...int16) int16 {
	for _, val := range values {
		if val != 0 {
			return val
		}
	}
	return 0
}

func CoalesceInt32(values ...int32) int32 {
	for _, val := range values {
		if val != 0 {
			return val
		}
	}
	return 0
}

func CoalesceInt64(values ...int64) int64 {
	for _, val := range values {
		if val != 0 {
			return val
		}
	}
	return 0
}

func CoalesceFloat32(values ...float32) float32 {
	for _, val := range values {
		if val != 0 {
			return val
		}
	}
	return 0
}

func CoalesceFloat64(values ...float64) float64 {
	for _, val := range values {
		if val != 0 {
			return val
		}
	}
	return 0
}

func CoalesceString(values ...string) string {
	for _, val := range values {
		if val != "" {
			return val
		}
	}
	return ""
}
