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

package widget

import (
	"github.com/adverax/echo"
)

const (
	AlertInfo    AlertType = "info"
	AlertSuccess AlertType = "siccess"
	AlertWarning AlertType = "warning"
	AlertDanger  AlertType = "danger"
)

type AlertType string

// Widget for display highlighted message.
type Alert struct {
	Hidden  bool        // Message is hidden
	Type    AlertType   // Type of message
	Message interface{} // Content of message
}

func (w *Alert) Render(
	ctx echo.Context,
) (interface{}, error) {
	if w.Message == nil || w.Hidden {
		return nil, nil
	}

	message, err := echo.RenderWidget(ctx, w.Message)
	if err != nil || message == nil {
		return nil, err
	}

	res := make(map[string]interface{}, 4)
	res["Type"] = w.Type
	res["Message"] = w.Message

	return res, nil
}
