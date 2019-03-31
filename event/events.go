package event

import (
	"context"
	"sync"
)

type (
	Subscriber func(ctx context.Context, event interface{}) error

	subscribers []Subscriber

	Messenger interface {
		Trigger(ctx context.Context, name string, event interface{}) error
	}

	Publisher interface {
		Messenger
		On(name string, subscriber Subscriber)
		Off(name string, subscriber *interface{})
	}

	publisher struct {
		subscribers map[string]subscribers
	}

	publisherSafe struct {
		publisher
		sync.RWMutex
	}
)

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

func (pub *publisherSafe) On(name string, action Subscriber) {
	pub.Lock()
	defer pub.Unlock()
	pub.publisher.On(name, action)
}

func (pub *publisherSafe) Off(name string, action *interface{}) {
	pub.Lock()
	defer pub.Unlock()
	pub.publisher.Off(name, action)
}

func (pub *publisherSafe) Trigger(ctx context.Context, name string, event interface{}) error {
	pub.RLock()
	defer pub.RUnlock()
	return pub.publisher.Trigger(ctx, name, event)
}

func New() Publisher {
	return &publisher{
		subscribers: make(map[string]subscribers, 4),
	}
}

func NewSafe() Publisher {
	return &publisherSafe{
		publisher: publisher{
			subscribers: make(map[string]subscribers, 4),
		},
	}
}
