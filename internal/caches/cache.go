package caches

import (
	"fmt"
	"github.com/mxmrykov/L0/internal/models"
	"sync"
)

type OrderCache struct {
	cache map[string]*models.Order
	mu    sync.Mutex
}

func NewCache() *OrderCache {
	return &OrderCache{
		cache: map[string]*models.Order{},
		mu:    sync.Mutex{},
	}
}

func (oc *OrderCache) CreateCache(or *models.Order) {

	oc.mu.Lock()
	oc.cache[or.OrderUid] = or
	oc.mu.Unlock()
	fmt.Printf("Cache written: %s\n", or.OrderUid)
}

func (oc *OrderCache) GetOrderByUid(uid string) *models.Order {
	return oc.cache[uid]
}

func (oc *OrderCache) GetOrders() map[string]*models.Order {
	return oc.cache
}
