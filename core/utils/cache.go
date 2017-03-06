// Package cache implements a simple goroutine safe cache with expiration time
// acceptes any value but it only gives back valid values ( not expired once ).
// Provides only 3 methods, Put, Get and Stop ( stop must be called when done with
// the timed cache to avoid memory leakage ) and a constructor New.

package tcache

import (
	"sync"
	"time"
)

type tCache struct {
	values map[string]*cacheObject
	mtx    *sync.RWMutex
	tick   *time.Ticker
	done   chan bool
	exp    time.Duration
}

type cacheObject struct {
	Value interface{}
	Valid time.Time
}

// New creates a new cache with minutes, which rappresent a interval at which
// old values are deleted and exp ( in minutes as well ) which sets the expiration
// time for values ( cannot be change once the cache has been istanciated ).
func New(minutes time.Duration, exp time.Duration) *tCache {
	cache := &tCache{
		make(map[string]*cacheObject{}),
		&sync.RWMutex{},
		time.NewTicker(time.Minute * minutes),
		make(chan bool, 1),
		exp,
	}
	cache.purger()

	return cache
}

// Put adds a valus to the cache
func (c *Cache) Put(key string, v interface{}) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.values[key] = append(c.values[key], &cacheObject{v, time.Now()})
}

// Get give you back the value assumining it hasn't be purged yet
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	elem, ok := c.values[key]
	if !ok {
		return nil, ok
	}

	return elem.Value, ok
}

// Stop Must be called otherwise the cache will leak
func (c *Cache) Stop() {
	c.done <- true
}

func (c *Cache) purger() {
	go func() {
		for {
			select {
			case stop := <-c.done:
				c.tick.Stop()
				return

			case now := <-c.tick.C:
				c.cleaner(now)
			}
		}
	}()
}

func (c *Cache) cleaner(now time.Time) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	for k, value := range c.values {
		if now.After(value.Valid.Add(c.exp)) {
			c.mtx.RUnlock()

			c.mtx.Lock()
			delete(c.values, k)
			c.mtx.Unlock()

			c.mtx.RLock()
		}
	}
}
