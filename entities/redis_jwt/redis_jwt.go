package redis_jwt

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Username: "",
		Password: "",
		DB:       0,
	})
	ctx = context.Background()
)

type RdbService interface {
	SaveJWT(userData RedisData) error
	GetJWT(user_id string) string
	DeleteJWT(userId string) error

	SaveRefreshToken(userData RedisData) error
	GetRefresh(user_id string) string
	DeleteRefreshToken(userId string) error
}

func New() RdbService {
	return &RedisData{}
}

func (u *RedisData) SaveRefreshToken(userData RedisData) error {
	pipeline := rdb.Pipeline()
	pipeline.HSet(ctx, userData.Id, "RefreshToken", userData.Token)
	pipeline.Expire(ctx, userData.Id, (24*time.Hour)*30)
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

func (u *RedisData) GetRefresh(user_id string) string {
	rdbToken := rdb.HGet(ctx, user_id, "RefreshToken")
	return rdbToken.Val()
}

func (u *RedisData) DeleteRefreshToken(userId string) error {
	err := rdb.HDel(ctx, userId, "RefreshToken").Err()
	if err != nil {
		return err
	}
	return nil
}

func (u *RedisData) SaveJWT(userData RedisData) error {
	pipeline := rdb.Pipeline()
	pipeline.HSet(ctx, userData.Id, "jwt", userData.Token)
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

func (u *RedisData) GetJWT(user_id string) string{
	rdbToken := rdb.HGet(ctx, user_id, "jwt")
	return rdbToken.Val()
}

func (u *RedisData) DeleteJWT(userId string) error {
	err := rdb.HDel(ctx, userId, "jwt").Err()
	if err != nil {
		return err
	}
	return nil
}
