package models

import (
	"culture/client"
	"github.com/go-redis/redis/v7"
	"log"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop"
)

// DB is a connection to your database to be used
// throughout your application.
var DB *pop.Connection

var REDIS *redis.Client

func init() {
	var err error
	env := envy.Get("GO_ENV", "development")
	DB, err = pop.Connect(env)
	if err != nil {
		log.Fatal(err)
	}
	pop.Debug = env == "development"
	REDIS, err = client.Connect(env)
	if err != nil {
		log.Fatal(err)
	}
}
