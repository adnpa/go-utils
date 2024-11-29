package local

import lru "github.com/hashicorp/golang-lru"

//https://github.com/hashicorp/golang-lru

type LruCache struct {
	*lru.Cache
}

func NewLruCache(size int) *LruCache {
	cache, _ := lru.New(size)
	return &LruCache{cache}
}

func (l *LruCache) Get(key string) (interface{}, bool) {
	v, ok := l.Cache.Get(key)
	return v, ok
}

func (l *LruCache) Set(key string, value interface{}) {

	l.Cache.Add(key, value)
}

func (l *LruCache) Del(key string) {
	l.Cache.Remove(key)
}
