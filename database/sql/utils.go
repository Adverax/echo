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

package sql

import (
	"context"
	"time"
)

type DbId int

const PrimaryDatabase DbId = 1

func scatter(n int, fn func(i int) error) error {
	errors := make(chan error, n)

	var i int
	for i = 0; i < n; i++ {
		go func(i int) { errors <- fn(i) }(i)
	}

	var err, innerErr error
	for i = 0; i < cap(errors); i++ {
		if innerErr = <-errors; innerErr != nil {
			err = innerErr
		}
	}

	return err
}

// Extract database context from context
func FromContext(ctx context.Context, key DbId) Scope {
	val := ctx.Value(key)
	if c, valid := val.(Scope); valid {
		return c
	}
	return nil
}

// Append scope into context
func ToContext(
	ctx context.Context,
	scope Scope,
) context.Context {
	return context.WithValue(ctx, scope.DbId(), scope)
}

func Heartbeart(
	ctx context.Context,
	db DB,
	interval time.Duration,
) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(interval):
			if err := db.Ping(); err != nil {
				return err
			}
		}
	}
}

// Exclusive open database for escape any concurrency.
func (dsc DSC) OpenForTest(
	ctx context.Context,
) DB {
	dsn := dsc.Primary()
	dsn.Database += "_test"
	db, err := dsn.Open(
		dsc.Driver,
		OpenExclusive(0x7ffffff, nil),
	)
	if err != nil {
		panic(err)
	}

	return db
}
