package controllers

import (
	dto "chat/DTO/userdto"
	"chat/entities/user"
	"encoding/json"
)

func ParseToDb(DTOdata any, dbData *user.User) error {
	data, err := json.Marshal(DTOdata)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &dbData)
	if err != nil {
		return err
	}
	return nil
}

func ParseToDTO[T dto.UserDTO | dto.CreateUserDTO | dto.GetUserDTO](dbUser user.User, userDTO *T) error {
	data, err := json.Marshal(dbUser)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, userDTO)
	return nil
}
