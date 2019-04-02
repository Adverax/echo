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

package data

import (
	"context"
	"errors"
	"time"
)

const (
	SortAsc  Sorting = true
	SortDesc Sorting = false
)

type Sorting bool

type Sort map[string]Sorting

type Pagination struct {
	Offset int64 `json:"offset"`
	Limit  int64 `json:"limit"`
}

// Data provider for receive external data
type Provider interface {
	// Get displayed records count
	Count(ctx context.Context) (int, error)
	// Get total records count
	Total(ctx context.Context) (int, error)
	// Import records
	Import(ctx context.Context, pagination *Pagination) error
	// Go to next row
	Next(ctx context.Context) error
}

// ArrayProvider is base for CustomArrayProvider
type ArrayProvider struct {
	Quantity int // Items count
	Index    int // Current index (starts from 1)
	Loader   func(ctx context.Context) (int, error)
	loaded   bool
}

func (provider *ArrayProvider) Count(ctx context.Context) (int, error) {
	err := provider.Import(ctx, nil)
	if err != nil {
		return 0, err
	}

	return provider.Quantity, nil
}

func (provider *ArrayProvider) Total(ctx context.Context) (int, error) {
	err := provider.Import(ctx, nil)
	if err != nil {
		return 0, err
	}

	return provider.Quantity, nil
}

func (provider *ArrayProvider) Import(ctx context.Context, pagination *Pagination) error {
	if provider.loaded || provider.Loader == nil {
		return nil
	}

	quantity, err := provider.Loader(ctx)
	if err != nil {
		return err
	}
	provider.Quantity = quantity
	provider.loaded = true

	return nil
}

func (provider *ArrayProvider) Next(ctx context.Context) error {
	index := provider.Index + 1
	if index > provider.Quantity {
		return errors.New("range check error")
	}

	provider.Index = index
	return nil
}

var (
	ErrNoMatch         = errors.New("no match")
	ErrRangeCheckError = errors.New("range check error")
)

// Current time
var Now = func() time.Time {
	return time.Now()
}
