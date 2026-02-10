package memory

import (
	"sync"
	"time"

	"github.com/GEtBUsyliVn/url-shortener/cache/model"
)

type MemoryRepository struct {
	Store    map[string]*model.MemoryCache
	mu       *sync.RWMutex
	duration time.Duration
}

func NewMemoryStorage(duration time.Duration) *MemoryRepository {
	return &MemoryRepository{
		Store:    make(map[string]*model.MemoryCache),
		mu:       &sync.RWMutex{},
		duration: duration,
	}
}

func (r *MemoryRepository) Set(key string, value string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Store[key] = &model.MemoryCache{
		Url:       value,
		ExpiresAt: time.Now().Add(r.duration),
	}
}

func (r *MemoryRepository) Get(key string) string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	val, ok := r.Store[key]
	if !ok {
		return ""
	}
	return val.Url
}

func (r *MemoryRepository) Del(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.Store, key)
}

func (r *MemoryRepository) DeleteExpired(now time.Time) (deleted int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for k, v := range r.Store {
		if now.After(v.ExpiresAt) {
			delete(r.Store, k)
			deleted++
		}
	}
	return
}
