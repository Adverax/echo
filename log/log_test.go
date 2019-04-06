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

package log

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	var buf bytes.Buffer
	l := NewEx(
		&buf,
		&buf,
		&buf,
		&buf,
		0,
		"#",
		false,
	)
	l.Trace("a")
	l.Info("b")
	l.Warning("c")
	l.Error("d")
	assert.Equal(t, "#\na\n#\nb\n#\nc\n#\nd\n", buf.String())
}
