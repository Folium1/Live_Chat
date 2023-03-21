package rediscontroller

import (
	dto "chat/DTO/redis_jwt"
	"chat/entities/redis_jwt"
)

type RdbController struct {
	rdb redis_jwt.RdbService
}

type Service interface {
}

func New(redis redis_jwt.RdbService) Service {
	return &RdbController{redis}
}

func (c *RdbController) SaveJWT(userData dto.RdbDTO) error {
	var dbData redis_jwt.UserJwt
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

func (c *RdbController) CompareJWT(userData dto.RdbDTO) (error, bool) {
	var dbData redis_jwt.UserJwt
	err := parseToDb(userData, &dbData)
	if err != nil {
		return err, false
	}
	if c.rdb.CompareJWT(dbData) {
		return nil, true
	}
	return nil, false
}
