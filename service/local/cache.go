package local

import (
	"errors"
	"strconv"
	"time"
)

type Cache struct {
	KV map[string]map[string]string
}

var errKeyNotExists = errors.New("cache key is not exists")
var errCacheExpired = errors.New("cache was expired")

// 获取缓存
func (cache *Cache) get(key string) (value string, err error) {
	// 缓存key不存在
	if cache.KV[key] == nil {
		return "", errKeyNotExists
	}
	// key过期
	expired, err := strconv.Atoi(cache.KV[key]["expired"])
	if err != nil {
		return "", errors.New("key expired is not int")
	}
	if int(time.Now().Unix()) > expired {
		return "", errCacheExpired
	}
	return cache.KV[key]["data"], nil
}
