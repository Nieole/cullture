package cache

import (
	"github.com/go-redis/cache/v7"
	"github.com/gobuffalo/envy"
	"github.com/vmihailenco/msgpack/v4"
	"log"
	"time"
)

var CACHE *cache.Codec

func init() {
	env := envy.Get("GO_ENV", "development")
	ring, err := Connect(env)
	if err != nil {
		log.Fatal(err)
	}
	CACHE = &cache.Codec{
		Redis: ring,
		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return msgpack.Unmarshal(b, v)
		},
	}
}

func Once(key string, value interface{}, load func() (interface{}, error), expiration time.Duration) error {
	err := CACHE.Once(&cache.Item{
		Key:        key,
		Object:     value,
		Func:       load,
		Expiration: expiration,
	})
	if err != nil {
		return err
	}
	return nil
}

func Clean(key string) error {
	return CACHE.Delete(key)
}
