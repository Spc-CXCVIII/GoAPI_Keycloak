package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var ch *cache.Cache

func InitCache() {
	c := cache.New(5*time.Minute, 10*time.Minute)
	ch = c
}

func AddData(key string, value interface{}) {
	ch.Set(key, value, cache.DefaultExpiration)
}

func GetData(key string) (interface{}, bool) {
	return ch.Get(key)
}

func DeleteData(key string) {
	ch.Delete(key)
}
