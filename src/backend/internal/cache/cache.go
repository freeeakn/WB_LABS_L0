package cache

import (
	"backend/internal/model"
	"sync"
)

type Cache struct {
	mu   sync.RWMutex
	data map[string]*model.Order
}

func New() *Cache {
	return &Cache{data: make(map[string]*model.Order)}
}

func (c *Cache) Set(o *model.Order) {
	c.mu.Lock()
	c.data[o.OrderUID] = o
	c.mu.Unlock()
}

func (c *Cache) Get(uid string) (*model.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	o, ok := c.data[uid]
	return o, ok
}

func (c *Cache) All() map[string]*model.Order {
	c.mu.RLock()
	defer c.mu.RUnlock()
	copy := make(map[string]*model.Order, len(c.data))
	for k, v := range c.data {
		copy[k] = v
	}
	return copy
}
