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

package event

import (
	"context"
)

type Subscriber func(ctx context.Context, event interface{}) error

type subscribers []Subscriber

type Messenger interface {
	Trigger(ctx context.Context, name string, event interface{}) error
}

type Publisher interface {
	Messenger
	On(name string, subscriber Subscriber)
	Off(name string, subscriber *interface{})
}

type publisher struct {
	subscribers map[string]subscribers
}

func (pub *publisher) On(name string, action Subscriber) {
	if actions, found := pub.subscribers[name]; found {
		pub.subscribers[name] = append(actions, action)
	} else {
		pub.subscribers[name] = subscribers{action}
	}
}

func (pub *publisher) Off(name string, action *interface{}) {
	if list, found := pub.subscribers[name]; found {
		for i, act := range list {
			if interface{}(act) == action {
				pub.subscribers[name] = append(list[:i], list[i+1:]...)
				break
			}
		}
	}
}

func (pub *publisher) Trigger(ctx context.Context, name string, event interface{}) error {
	if actions, found := pub.subscribers[name]; found {
		for _, action := range actions {
			err := action(ctx, event)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func New() Publisher {
	return &publisher{
		subscribers: make(map[string]subscribers, 4),
	}
}
