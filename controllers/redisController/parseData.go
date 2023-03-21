package rediscontroller

import (
	"chat/entities/redis_jwt"
	"encoding/json"
)

func parseToDb(dtoData any, dbData *redis_jwt.UserJwt) error {
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
