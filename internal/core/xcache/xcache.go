package xcache

import (
	"errors"
	"github.com/dgraph-io/ristretto"
	"time"
)

var IsCacheExists = false
var myCache *ristretto.Cache
var config = &ristretto.Config{
	NumCounters: 1e7,     // 10M
	MaxCost:     1 << 30, // 1GB
	BufferItems: 64,      // 64 is best, for fine tuning
}

func Get(key interface{}) (interface{}, error) {
	var err error
	if !IsCacheExists {
		myCache, err = ristretto.NewCache(config)
		if err != nil {
			return nil, err
		}
		IsCacheExists = true
	}

	v, ok := myCache.Get(key)
	if ok {
		return v, nil
	} else {
		return nil, errors.New("cache key not found")
	}
}

func Set(key, value interface{}) error {
	var err error
	if !IsCacheExists {
		myCache, err = ristretto.NewCache(config)
		if err != nil {
			return err
		}
		IsCacheExists = true
	}

	for i := 0; i < 10; i++ {
		isSet := myCache.Set(key, value, 0)
		if isSet {
			time.Sleep(10 * time.Nanosecond) // needed for making sure of a successful fetch
			_, isGet := myCache.Get(key)
			if isGet {
				return nil
			} else {
				continue
			}
		} else {
			continue
		}
	}
	return errors.New("set is failed")
}
