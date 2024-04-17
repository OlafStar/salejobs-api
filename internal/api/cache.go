package api

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
)

type CacheID string

const (
	CacheIDAdvBase CacheID = "adv_page_%d_limit_%d" 
	CacheIDAdvc    CacheID = "advC"
)

type allCache struct {
	cache *cache.Cache
}

const (
	defaultExpiration = 5 * time.Minute
	purgeTime         = 10 * time.Minute
)

func NewCache() *allCache {
	Cache := cache.New(defaultExpiration, purgeTime)
	return &allCache{
		cache: Cache,
	}
}

func (c *allCache) read(id CacheID) (item []byte, ok bool) {
	cacheItem, ok := c.cache.Get(string(id))
	if !ok {
		return nil, false
	}

	switch v := cacheItem.(type) {
	case []byte:
		return v, true
	default:
		res, err := json.Marshal(v)
		if err != nil {
			return nil, false
		}
		return res, true
	}
}

func (c *allCache) update(id CacheID, item any) {
	c.cache.Set(string(id), item, cache.DefaultExpiration)
}

func (c *allCache) clear(id CacheID) {
	c.cache.Delete(string(id))
}

func (c *allCache) clearByPattern(pattern string) {
	items := c.cache.Items()
	for key := range items {
		if strings.Contains(key, pattern) {
			c.cache.Delete(key)
		}
	}
}
