package controllers

import (
	"chat/entities/user"
	"encoding/json"
)

func ParseUserData(DTOdata any, dbData user.User) (user.User, error) {
	data, err := json.Marshal(DTOdata)
	if err != nil {
		return user.User{}, err
	}
	err = json.Unmarshal(data, &dbData)
	if err != nil {
		return user.User{}, err
	}
	return dbData, nil
}
