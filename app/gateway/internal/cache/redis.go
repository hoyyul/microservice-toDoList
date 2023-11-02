package cache

import (
	"micro-toDoList/pkg/util/jwts"
	"time"

	"github.com/go-redis/redis"
)

var prefix string = "user_logout_"

type RedisService struct {
	*redis.Client
}

// 接住返回的全局值
func NewRedisService() *RedisService {
	return &RedisService{ConnectRedisDB()}
}

// store token for logout user until expiration
func (s *RedisService) SaveToken(claim jwts.CustomClaim, token string) {
	expire := claim.ExpiresAt
	now := time.Now()
	duration := expire.Time.Sub(now)

	s.Set(prefix+token, "", duration)
}

// check if token belonged to a logout user
func (s *RedisService) CheckIfLogout(token string) bool {
	_, err := s.Get(prefix + token).Result()
	return err == nil
}
