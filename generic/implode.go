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

import (
	"strconv"
)

func ImplodeInt(list []int, delim string) string {
	if list == nil {
		return ""
	}
	var res = ""
	var first = true
	for _, id := range list {
		if first {
			res = res + strconv.FormatInt(int64(id), 10)
			first = false
		} else {
			res = res + delim + strconv.FormatInt(int64(id), 10)
		}
	}
	return res
}

func ImplodeInt8(list []int8, delim string) string {
	if list == nil {
		return ""
	}
	var res = ""
	var first = true
	for _, id := range list {
		if first {
			res = res + strconv.FormatInt(int64(id), 10)
			first = false
		} else {
			res = res + delim + strconv.FormatInt(int64(id), 10)
		}
	}
	return res
}

func ImplodeInt16(list []int16, delim string) string {
	if list == nil {
		return ""
	}
	var res = ""
	var first = true
	for _, id := range list {
		if first {
			res = res + strconv.FormatInt(int64(id), 10)
			first = false
		} else {
			res = res + delim + strconv.FormatInt(int64(id), 10)
		}
	}
	return res
}

func ImplodeInt32(list []int32, delim string) string {
	if list == nil {
		return ""
	}
	var res = ""
	var first = true
	for _, id := range list {
		if first {
			res = res + strconv.FormatInt(int64(id), 10)
			first = false
		} else {
			res = res + delim + strconv.FormatInt(int64(id), 10)
		}
	}
	return res
}

func ImplodeInt64(list []int64, delim string) string {
	if list == nil {
		return ""
	}
	var res = ""
	var first = true
	for _, id := range list {
		if first {
			res = res + strconv.FormatInt(id, 10)
			first = false
		} else {
			res = res + delim + strconv.FormatInt(id, 10)
		}
	}
	return res
}

func ImplodeUint(list []uint, delim string) string {
	if list == nil {
		return ""
	}
	var res = ""
	var first = true
	for _, id := range list {
		if first {
			res = res + strconv.FormatUint(uint64(id), 10)
			first = false
		} else {
			res = res + delim + strconv.FormatUint(uint64(id), 10)
		}
	}
	return res
}

func ImplodeUint8(list []uint8, delim string) string {
	if list == nil {
		return ""
	}
	var res = ""
	var first = true
	for _, id := range list {
		if first {
			res = res + strconv.FormatUint(uint64(id), 10)
			first = false
		} else {
			res = res + delim + strconv.FormatUint(uint64(id), 10)
		}
	}
	return res
}

func ImplodeUint16(list []uint16, delim string) string {
	if list == nil {
		return ""
	}
	var res = ""
	var first = true
	for _, id := range list {
		if first {
			res = res + strconv.FormatUint(uint64(id), 10)
			first = false
		} else {
			res = res + delim + strconv.FormatUint(uint64(id), 10)
		}
	}
	return res
}

func ImplodeUint32(list []uint32, delim string) string {
	if list == nil {
		return ""
	}
	var res = ""
	var first = true
	for _, id := range list {
		if first {
			res = res + strconv.FormatUint(uint64(id), 10)
			first = false
		} else {
			res = res + delim + strconv.FormatUint(uint64(id), 10)
		}
	}
	return res
}

func ImplodeUint64(list []uint64, delim string) string {
	if list == nil {
		return ""
	}
	var res = ""
	var first = true
	for _, id := range list {
		if first {
			res = res + strconv.FormatUint(id, 10)
			first = false
		} else {
			res = res + delim + strconv.FormatUint(id, 10)
		}
	}
	return res
}
