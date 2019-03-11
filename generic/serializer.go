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
	"encoding/json"
)

type Serializer interface {
	Serialize(src interface{}) (string, error)
	Unserialize(src string, dst interface{}) error
}

type serializer struct{}

func (engine *serializer) Serialize(src interface{}) (string, error) {
	data, err := json.Marshal(src)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (engine *serializer) Unserialize(src string, dst interface{}) error {
	return json.Unmarshal([]byte(src), dst)
}

var ser = new(serializer)

func NewSerializer() Serializer {
	return ser
}
