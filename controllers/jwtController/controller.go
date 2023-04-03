package jwtcontroller

import (
	dto "chat/DTO/jwt"
	"chat/entities/jwt"
	"fmt"
)

type rdbController struct {
	redis redis_jwt.RedisDbService
}

type Service interface {
	SaveToken(data dto.RedisDto) error
	DeleteToken(data dto.RedisDto) error
	GetToken(data dto.RedisDto) (string, error)
}

func New() Service {
	redis := redis_jwt.New()
	return &rdbController{redis}
}

func (c *rdbController) SaveToken(data dto.RedisDto) error {
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

func (c *rdbController) DeleteToken(data dto.RedisDto) error {
	var dbData redis_jwt.RedisData
	parseToDb(data, &dbData)
	err := c.redis.DeleteToken(dbData)
	if err != nil {
		return err
	}
	return nil
}

func (c *rdbController) GetToken(data dto.RedisDto) (string, error) {
	var dbData redis_jwt.RedisData
	parseToDb(data, &dbData)
	rdbToken := c.redis.GetToken(dbData)
	if rdbToken == "" {
		return "", fmt.Errorf("no %v token found", data.Key)
	}
	return rdbToken, nil
}
