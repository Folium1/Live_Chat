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
	SaveJWT(userData UserJwt) error
	CompareJWT(userData UserJwt) bool
	DeleteJWT(userId string) error
}

func New() RdbService {
	return &UserJwt{}
}

func (u *UserJwt) SaveJWT(userData UserJwt) error {
	pipline := rdb.Pipeline()
	pipline.HSet(ctx, userData.Id, "jwt", userData.Token)
	pipline.Expire(ctx, userData.Id, 15*time.Minute)
	result, err := pipline.Exec(ctx)
	if err != nil {
		return err
	}
	for _, command := range result {
		if err = command.Err(); err != nil {
			pipline.Discard()
			return err
		}
	}
	return nil
}

func (u *UserJwt) CompareJWT(userData UserJwt) bool {
	rdbToken := rdb.HGet(ctx, userData.Id, "jwt")
	if userData.Token == rdbToken.Val() {
		return true
	}
	return false
}

func (u *UserJwt) DeleteJWT(userId string) error {
	err := rdb.HDel(ctx, userId, "jwt").Err()
	if err != nil {
		return err
	}
	return nil
}
