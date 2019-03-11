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

package cache

import "time"

type Cache interface {
	// Get value by key.
	Get(key string, dst interface{}) error
	// GetMulti is a batch version of Get.
	GetMulti(dict map[string]interface{}) error
	// Set value with key and expire time.
	Set(key string, val interface{}, timeout time.Duration) error
	// Check if value exists or not.
	IsExist(key string) (bool, error)
	// Delete cached value by key.
	Delete(key string) error
	// Increase cached int value by key, as a counter.
	Increase(key string) error
	// Decrease cached int value by key, as a counter.
	Decrease(key string) error
	// Clear cache.
	Clear() error
	// Stops the background worker.
	Stop()
}
