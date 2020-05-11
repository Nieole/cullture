package cache

import (
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/cache/v7"
	"github.com/gobuffalo/envy"
)

//CACHE CACHE
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
			return json.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return json.Unmarshal(b, v)
		},
	}
}

//Once Once
func Once(key string, value interface{}, load func() (interface{}, error), expiration time.Duration) error {
	return CACHE.Once(&cache.Item{
		Key:        key,
		Object:     value,
		Func:       load,
		Expiration: expiration,
	})
}

//Clean Clean
func Clean(key string) error {
	return CACHE.Delete(key)
}
