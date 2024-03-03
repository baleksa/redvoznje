package main

import (
	"sync"
	"time"
)

type cache struct {
	store sync.Map
}

func newCache() *cache {
	return &cache{}
}

func (c *cache) set(key string, value interface{}) {
	c.store.Store(key, value)
	time.AfterFunc(durationUntilTomorrow(), func() {
		c.store.Delete(key)
	})
}

func (c *cache) get(key string) (interface{}, bool) {
	return c.store.Load(key)
}

func durationUntilTomorrow() time.Duration {
	now := time.Now()
	yyyy, mm, dd := now.Date()
	tomorrow := time.Date(yyyy, mm, dd+1, 0, 0, 0, 0, now.Location())
	return time.Until(tomorrow)
}
