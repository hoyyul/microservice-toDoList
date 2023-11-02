package cache

import (
	"context"
	"micro-toDoList/global"
	"time"

	"github.com/go-redis/redis"
)

// 全局值
var RedisClient *redis.Client

// 初始化全局值
func InitRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     global.Config.Redis.Address,
		Password: global.Config.Redis.Password,
		DB:       0,
		PoolSize: global.Config.Redis.PoolSize,
	})
	_, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	_, err := rdb.Ping().Result()
	if err != nil {
		global.Logger.Errorf("Failed to connect redis %s", global.Config.Redis.Address)
	}

	RedisClient = rdb
}

// 返还全局值
func ConnectRedisDB() *redis.Client {
	return RedisClient
}

// 这是一个很好的和数据库/缓存交互的例子
