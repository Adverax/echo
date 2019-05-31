package arbiter

import "sync"

type latch struct {
	sync.Mutex
	usage int
}

type localArbiter struct {
	sync.Mutex
	latches map[string]*latch
}

func (arbiter *localArbiter) Lock(key string) {
	arbiter.Mutex.Lock()
	aLatch, ok := arbiter.latches[key]
	if !ok {
		aLatch := new(latch)
		arbiter.latches[key] = aLatch
	}
	aLatch.usage++
	arbiter.Mutex.Unlock()
	aLatch.Lock()
}

func (arbiter *localArbiter) Unlock(key string) {
	arbiter.Mutex.Lock()
	if latch, ok := arbiter.latches[key]; ok {
		latch.usage--
		if latch.usage == 0 {
			delete(arbiter.latches, key)
		}
		arbiter.Mutex.Unlock()
		latch.Unlock()
	} else {
		arbiter.Mutex.Unlock()
	}
}

func NewLocal() Arbiter {
	return &localArbiter{
		latches: make(map[string]*latch, 1024),
	}
}
