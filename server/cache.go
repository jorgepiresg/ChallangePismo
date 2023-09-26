package server

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

func (s *server) startCache() *redis.Client {

	cache := redis.NewClient(&redis.Options{
		Addr: s.config.Cache.Addr,
	})

	err := cache.Ping(context.Background()).Err()
	if err != nil {
		log.Fatal("cache ping: ", err.Error())
	}

	log.Println("cache started")
	return cache
}
