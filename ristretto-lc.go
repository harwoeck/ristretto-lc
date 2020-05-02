package ristrettolc

import (
	"sync"
	"time"

	"github.com/dgraph-io/ristretto"
)

type LoadingCache interface {
	Cache() *ristretto.Cache
	Get(key interface{}) (interface{}, bool)
}

type CacheBuilderFunc func(key interface{}) (val interface{}, cost int64, ok bool)

func NewLoadingCache(cache *ristretto.Cache, builder CacheBuilderFunc) LoadingCache {
	return &loadingCache{
		cache:   cache,
		builder: builder,
	}
}

type loadingCache struct {
	cache   *ristretto.Cache
	builder CacheBuilderFunc
}

type cacheedItem struct {
	val         interface{}
	lastRefresh time.Time
	baton       sync.RWMutex
	refreshing  bool
}

func (lc *loadingCache) Cache() *ristretto.Cache {
	return lc.cache
}

func (lc *loadingCache) Get(key interface{}) (interface{}, bool) {
	val, ok := lc.cache.Get(key)
	if !ok {
		val, cost, ok := lc.builder(key)
		if !ok {
			return nil, false
		}

		lc.cache.Set(key, val, cost)
	}

	return val, true
}
