package usersController

import (
	dto "chat/DTO/userdto"
	"chat/entity/user"
	"encoding/json"
)

type userController struct {
	db user.UserService
}

func New(msg user.UserService) UserController {
	return &userController{msg}
}

type UserController interface {
	CreateUser(newUser dto.CreateUserDTO) error
}

func (u *userController) CreateUser(newUser dto.CreateUserDTO) error {
	var dbUser user.User
	data, err := json.Marshal(newUser)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &dbUser)
	if err != nil {
		return err
	}
	err = u.db.CreateUser(dbUser)
	if err != nil {
		return err
	}
	return nil
}
