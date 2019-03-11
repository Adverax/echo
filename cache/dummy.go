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

type DummyCache struct{}

func (cache *DummyCache) Get(key string, val interface{}) error {
	return nil
}

func (cache *DummyCache) GetMulti(keys map[string]interface{}) error {
	return nil
}

func (cache *DummyCache) Set(key string, val interface{}, timeout time.Duration) error {
	return nil
}

func (cache *DummyCache) IsExist(key string) (bool, error) {
	return false, nil
}

func (cache *DummyCache) Delete(key string) error {
	return nil
}

func (cache *DummyCache) Increase(key string) error {
	return nil
}

func (cache *DummyCache) Decrease(key string) error {
	return nil
}

func (cache *DummyCache) Clear() error {
	return nil
}

func (cache *DummyCache) Stop() {
	// nothing
}
