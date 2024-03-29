package jwtcontroller

import (
	"chat/entities/jwt"
	"encoding/json"
)

func parseToDb(dtoData any, dbData *redis_jwt.RedisData) error {
	data, err := json.Marshal(dtoData)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &dbData)
	if err != nil {
		return err
	}
	return nil
}
