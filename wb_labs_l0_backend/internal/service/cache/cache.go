package cache

import (
	"sync"

	"wb_labs_l0_backend/internal/domain/order"
)

type Cache interface {
	Get(id string) (*order.Order, bool)
	Set(id string, o *order.Order)
	List() []*order.Order
	Len() int
}

type InMemoryCache struct {
	mu    sync.RWMutex
	store map[string]*order.Order
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		store: make(map[string]*order.Order),
	}
}

func (c *InMemoryCache) Get(id string) (*order.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	o, ok := c.store[id]
	return o, ok
}

func (c *InMemoryCache) Set(id string, o *order.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[id] = o
}

func (c *InMemoryCache) List() []*order.Order {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]*order.Order, 0, len(c.store))
	for _, v := range c.store {
		out = append(out, v)
	}
	return out
}

func (c *InMemoryCache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.store)
}
