package RedisGo

import "github.com/go-redis/redis"

var client *redis.Client

func Init() {
	client = redis.NewClient(&redis.Options{
		Addr:     "0.0.0.0:32768",
		Password: "",
		DB:       0,
	})
}

func GetClient() *redis.Client {
	return client
}
