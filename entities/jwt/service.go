package redis_jwt

import (
	"context"
	"fmt"

	entity "chat/entities"
)

var (
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

func (r *RedisData) SaveToken(userData RedisData) error {
	rds, err := entity.RedisConnect()
	if err != nil {
		return err
	}
	pipeline := rds.Pipeline()
	pipeline.HSet(ctx, fmt.Sprintf("auth:%v:%v", userData.Id, userData.Key), userData.Id, userData.Token)
	pipeline.Expire(ctx, fmt.Sprintf("auth:%v:%v", userData.Id, userData.Key), userData.Expiration)
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

func (r *RedisData) GetToken(data RedisData) string {
	rds, err := entity.RedisConnect()
	if err != nil {
		return ""
	}
	RedisDbToken := rds.HGet(ctx, fmt.Sprintf("auth:%v:%v", data.Id, data.Key), data.Id)
	return RedisDbToken.Val()
}

func (r *RedisData) DeleteToken(data RedisData) error {
	rds, err := entity.RedisConnect()
	if err != nil {
		return err
	}
	err = rds.HDel(ctx, fmt.Sprintf("auth:%v:%v", data.Id, data.Key), data.Id).Err()
	if err != nil {
		return err
	}
	return nil
}
