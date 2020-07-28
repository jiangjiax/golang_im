package db

import (
	"golang_im/config"

	"github.com/go-redis/redis"
)

var RedisCli *redis.Client

func Redis_init() {
	addr := config.DBConf.RedisIP
	password := config.DBConf.RedisPwd
	RedisCli = redis.NewClient(&redis.Options{
		Addr:     addr,
		DB:       0,
		Password: password,
	})

	_, err := RedisCli.Ping().Result()
	if err != nil {
		panic(err)
	}
}
