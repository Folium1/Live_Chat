package rediscontroller

import (
	dto "chat/DTO/redis_jwt"
	"chat/entities/redis_jwt"
	"errors"
)

type RdbController struct {
	redis redis_jwt.RedisDbService
}

type Service interface {
	SaveToken(data dto.RedisDto) error
	DeleteJWT(data dto.RedisDto) error
	GetToken(data dto.RedisDto) (string, error)
}

func New(redis redis_jwt.RedisDbService) Service {
	return &RdbController{redis}
}

func (c *RdbController) SaveToken(data dto.RedisDto) error {
	var dbData redis_jwt.RedisData
	err := parseToDb(data, &dbData)
	if err != nil {
		return err
	}
	err = c.redis.SaveToken(dbData)
	if err != nil {
		return err
	}
	return nil
}

func (c *RdbController) DeleteJWT(data dto.RedisDto) error {
	var dbData redis_jwt.RedisData
	parseToDb(data, &dbData)
	err := c.redis.DeleteToken(dbData)
	if err != nil {
		return err
	}
	return nil
}

func (c *RdbController) GetToken(data dto.RedisDto) (string, error) {
	rdbToken := c.redis.GetToken(redis_jwt.RedisData(data))
	if rdbToken == "" {
		return "", errors.New("no refresh token found")
	}
	return rdbToken, nil
}
