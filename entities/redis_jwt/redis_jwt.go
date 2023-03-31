package redis_jwt

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	RedisDb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("RedisAddr"),
		Username: os.Getenv("RedisUsername"),
		Password: os.Getenv("RedisPassword"),
	})
	ctx = context.Background()
)

type RedisDbService interface {
	SaveToken(userData RedisData) error
	GetToken(data RedisData) string
	DeleteToken(data RedisData) error
}

func New() RedisDbService {
	return &RedisData{}
}

func (u *RedisData) SaveToken(userData RedisData) error {
	pipeline := RedisDb.Pipeline()
	pipeline.HSet(ctx, userData.Id, userData.Key, userData.Token)
	pipeline.Expire(ctx, userData.Id, 15*time.Minute)
	result, err := pipeline.Exec(ctx)
	if err != nil {
		return err
	}
	for _, command := range result {
		if err = command.Err(); err != nil {
			pipeline.Discard()
			return err
		}
	}
	return nil
}

func (u *RedisData) DeleteRefreshToken(userId string) error {
	err := RedisDb.HDel(ctx, userId, "RefreshToken").Err()
	if err != nil {
		return err
	}
	return nil
}

func (u *RedisData) GetToken(data RedisData) string {
	RedisDbToken := RedisDb.HGet(ctx, data.Id, data.Key)
	return RedisDbToken.Val()
}

func (u *RedisData) DeleteToken(data RedisData) error {
	err := RedisDb.HDel(ctx, data.Id, data.Key).Err()
	if err != nil {
		return err
	}
	return nil
}
