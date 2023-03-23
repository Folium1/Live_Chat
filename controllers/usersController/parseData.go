package usersController

import (
	dto "chat/DTO/userdto"
	"chat/entities/user"
	"encoding/json"
)

func parseUserToDb(DTOdata any, dbData *user.User) error {
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

func parseUserToDTO[T dto.UserDTO | dto.CreateUserDTO | dto.GetUserDTO | dto.ChatUserDTO](dbUser user.User, userDTO *T) error {
	data, err := json.Marshal(dbUser)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, userDTO)
	if err != nil {
		return err
	}
	return nil
}

