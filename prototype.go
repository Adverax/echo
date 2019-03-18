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
	"github.com/adverax/echo/cache"
	"github.com/adverax/echo/cache/memory"
	"github.com/adverax/echo/data"
	"net/url"
	"time"

	"github.com/adverax/echo/generic"
)

type DefaultMessageFamily map[uint32]string

func (messages DefaultMessageFamily) Fetch(
	ctx stdContext.Context,
	id uint32,
) (string, error) {
	if msg, ok := messages[id]; ok {
		return msg, nil
	}

	return "", nil
}

type DefaultResourceFamily map[uint32]string

func (messages DefaultResourceFamily) Fetch(
	ctx stdContext.Context,
	id uint32,
) (string, error) {
	if msg, ok := messages[id]; ok {
		return msg, nil
	}

	return "", nil
}

type DefaultDataSetFamily map[uint32]DataSet

func (family DefaultDataSetFamily) Fetch(
	ctx stdContext.Context,
	id uint32,
) (DataSet, error) {
	if ds, ok := family[id]; ok {
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

type DefaultMessageManager struct {
	family MessageFamily
}

func (manager *DefaultMessageManager) Find(ctx stdContext.Context, id uint32, lang uint16, region uint16) (string, error) {
	return manager.family.Fetch(ctx, id)
}

type DefaultResourceManager struct {
	family ResourceFamily
}

func (manager *DefaultResourceManager) Find(ctx stdContext.Context, id uint32, lang uint16, region uint16) (string, error) {
	return manager.family.Fetch(ctx, id)
}

type DefaultDataSetManager struct {
	family DataSetFamily
}

func (manager *DefaultDataSetManager) Find(ctx stdContext.Context, doc uint32, lang uint16) (DataSet, error) {
	return manager.family.Fetch(ctx, doc)
}

func (manager *DefaultDataSetManager) FindAll(ctx stdContext.Context, doc uint32) (DataSets, error) {
	return nil, data.ErrNoMatch
}

var (
	DefaultMessages  = make(DefaultMessageFamily)
	DefaultResources = make(DefaultResourceFamily)
	DefaultDataSets  = make(DefaultDataSetFamily)

	Defaults = struct {
		Messages        MessageManager
		Resources       ResourceManager
		DataSets        DataSetManager
		UrlLinker       UrlLinker
		Cache           cache.Cache
		Locale          Locale
		MessageManager  MessageManager
		ResourceManager ResourceManager
		DataSetManager  DataSetManager
	}{
		UrlLinker: &DefaultUrlLinker{},
		Cache:     memory.New(memory.Options{}),
		Messages: &DefaultMessageManager{
			family: DefaultMessages,
		},
		Resources: &DefaultResourceManager{
			family: DefaultResources,
		},
		DataSets: &DefaultDataSetManager{
			family: DefaultDataSets,
		},
		Locale: &BaseLocale{
			DateFormat:     generic.DateFormat,
			TimeFormat:     generic.TimeFormat,
			DateTimeFormat: generic.DateTimeFormat,
			Lang:           1,
			TZone:          1,
			Loc:            time.UTC,
			Messages:       DefaultMessages,
			Resources:      DefaultResources,
			DataSets:       DefaultDataSets,
		},
	}
)
