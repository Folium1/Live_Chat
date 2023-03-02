package usersController

import (
	dto "chat/DTO/userdto"
	parser "chat/controllers"
	"chat/entities/user"
	"log"
)

type userController struct {
	db user.UserService
}

func New(msg user.UserService) UserController {
	return &userController{msg}
}

type UserController interface {
	CreateUser(newUser dto.CreateUserDTO) error
	GetUser(userData dto.GetUserDTO) (int, error)
}

func (c *userController) CreateUser(newUser dto.CreateUserDTO) error {
	dbUser, err := parser.ParseUserData(newUser, user.User{})
	if err != nil {
		log.Printf("couldn't parse dto data to db struct, err: %v", err)
	}
	err = c.db.CreateUser(dbUser)
	if err != nil {
		return err
	}
	return nil
}

func (c *userController) GetUser(userData dto.GetUserDTO) (int, error) {
	dbUser, err := parser.ParseUserData(userData, user.User{})
	if err != nil {
		log.Printf("couldn't parse dto data to db struct, err: %v", err)
	}
	id, err := c.db.GetUser(dbUser)
	if err != nil {
		log.Printf("couldn't get user, err: %v", err)
		return 0, err
	}
	return id, nil
}
