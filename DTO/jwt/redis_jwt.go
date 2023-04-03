package redisjwt

import "time"

type RedisDto struct {
	Id         string        `json:"user_id"`
	Key        string        `json:"key"`
	Token      string        `json:"jwt"`
	Expiration time.Duration `json:"exp"`
}
