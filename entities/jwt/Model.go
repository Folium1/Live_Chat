package redis_jwt

import "time"

type RedisData struct {
	Id    string `json:"user_id"`
	Key   string `json:"key"`
	Token string `json:"jwt"`
	Expiration time.Duration `json:"exp"`
}
