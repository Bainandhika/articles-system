package connections

import (
	"articles-system/app/configs"
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func InitRedis()  (*redis.Client, error) {
	domainFunc := "[ connections.InitRedis ]"
	redisConfig := configs.Config.Redis

	redisConn := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
		Username: redisConfig.Username,
		Password: redisConfig.Password,
	})

	err := redisConn.Ping(context.Background()).Err()
	if err != nil {
		return nil, fmt.Errorf("%s error redis ping: %v", domainFunc, err)
	}

	return redisConn, nil
}
