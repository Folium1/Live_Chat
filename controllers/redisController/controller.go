package rediscontroller

import (
	dto "chat/DTO/redis_jwt"
	"chat/entities/redis_jwt"
	"errors"
)

type RdbController struct {
	rdb redis_jwt.RdbService
}

type Service interface {
	SaveJWT(userData dto.RdbDTO) error
	DeleteJWT(userId string) error
	GetJWT(user_id string) (string, error)

	SaveRefreshToken(userData dto.RdbDTO) error
	GetRefreshToken(user_Id string) (string, error)
	DeleteRefreshToken(userId string) error
}

func New(redis redis_jwt.RdbService) Service {
	return &RdbController{redis}
}

func (c *RdbController) SaveJWT(userData dto.RdbDTO) error {
	var dbData redis_jwt.RedisData
	err := parseToDb(userData, &dbData)
	if err != nil {
		return err
	}
	err = c.rdb.SaveJWT(dbData)
	if err != nil {
		return err
	}
	return nil
}

func (c *RdbController) DeleteJWT(userId string) error {
	err := c.DeleteJWT(userId)
	if err != nil {
		return err
	}
	return nil
}

func (c *RdbController) GetJWT(user_id string) (string, error) {
	rdbToken := c.rdb.GetJWT(user_id)
	if rdbToken == "" {
		return "", errors.New("no refresh token found")
	}
	return rdbToken, nil
}

func (c *RdbController) SaveRefreshToken(userData dto.RdbDTO) error {
	var dbData redis_jwt.RedisData
	err := parseToDb(userData, &dbData)
	if err != nil {
		return err
	}
	err = c.rdb.SaveRefreshToken(dbData)
	if err != nil {
		return err
	}
	return nil
}

func (c *RdbController) GetRefreshToken(user_Id string) (string, error) {
	token := c.rdb.GetRefresh(user_Id)
	if token == "" {
		return "", errors.New("no refresh token found")
	}
	return token, nil
}

func (c *RdbController) DeleteRefreshToken(userId string) error {
	err := c.rdb.DeleteRefreshToken(userId)
	if err != nil {
		return err
	}
	return nil
}
