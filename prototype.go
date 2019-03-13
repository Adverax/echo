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

package echo

import (
	stdContext "context"
	"fmt"
	"github.com/adverax/echo/data"
	"github.com/adverax/echo/generic"
	"net/http"
	"net/url"
	"time"

	"github.com/adverax/echo/cache/memory"
)

type DefaultMessageManager map[uint32]string

func (messages DefaultMessageManager) Fetch(
	ctx stdContext.Context,
	id uint32,
) (string, error) {
	if msg, ok := messages[id]; ok {
		return msg, nil
	}

	return "", nil
}

type DefaultResourceManager map[uint32]string

func (messages DefaultResourceManager) Fetch(
	ctx stdContext.Context,
	id uint32,
) (string, error) {
	if msg, ok := messages[id]; ok {
		return msg, nil
	}

	return "", nil
}

type DefaultDataSetManager map[uint32]data.Set

func (datasets DefaultDataSetManager) Fetch(
	ctx stdContext.Context,
	id uint32,
) (data.Set, error) {
	if ds, ok := datasets[id]; ok {
		return ds, nil
	}

	return nil, nil
}

type DefaultUrlLinker struct{}

func (linker *DefaultUrlLinker) Render(ctx Context, u *url.URL) (string, error) {
	return u.String(), nil
}

func (linker *DefaultUrlLinker) Expand(ctx Context, url string) string {
	return url
}

func (linker *DefaultUrlLinker) Collapse(ctx Context, url string) string {
	return url
}

type DefaultSessionManager struct{}

func (manager *DefaultSessionManager) Load(
	ctc stdContext.Context,
	request *http.Request,
) (Session, error) {
	return nil, fmt.Errorf("abstact method for session, loading")
}

var (
	DefaultMessages  = &DefaultMessageManager{}
	DefaultResources = &DefaultResourceManager{}
	DefaultDataSets  = &DefaultDataSetManager{}
	DefaultLinker    = &DefaultUrlLinker{}
	DefaultCache     = memory.New(memory.Options{})
	DefaultSessions  = &DefaultSessionManager{}
	DefaultLocale    = &BaseLocale{
		DateFormat:     generic.DateFormat,
		TimeFormat:     generic.TimeFormat,
		DateTimeFormat: generic.DateTimeFormat,
		Lang:           1,
		TZone:          1,
		Loc:            time.UTC,
		Messages:       DefaultMessages,
		Resources:      DefaultResources,
		DataSets:       DefaultDataSets,
	}
)
