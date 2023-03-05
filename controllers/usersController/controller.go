package usersController

import (
	"log"

	dto "chat/DTO/userdto"
	parser "chat/controllers"
	"chat/entities/user"
)

type userController struct {
	db user.UserService
}

func New(msg user.UserService) UserController {
	return &userController{msg}
}

type UserController interface {
	CreateUser(newUser dto.CreateUserDTO) error
	GetUser(userData dto.GetUserDTO) (dto.GetUserDTO, error)
}

func (c *userController) CreateUser(newUser dto.CreateUserDTO) error {
	var dbUser user.User
	err := parser.ParseToDb(newUser, &dbUser)
	if err != nil {
		log.Printf("couldn't parse dto data to db struct, err: %v", err)
	}
	err = c.db.CreateUser(dbUser)
	if err != nil {
		return err
	}
	return nil
}

func (c *userController) GetUser(userData dto.GetUserDTO) (dto.GetUserDTO, error) {
	var dbUser user.User
	err := parser.ParseToDb(userData, &dbUser)
	if err != nil {
		log.Printf("couldn't parse dto data to db struct, err: %v", err)
	}
	userFromDb, err := c.db.GetUser(dbUser)
	if err != nil {
		log.Printf("couldn't get user, err: %v", err)
		return dto.GetUserDTO{}, err
	}
	var userDTO dto.GetUserDTO
	err = parser.ParseToDTO(userFromDb, &userDTO)
	if err != nil {
		log.Printf("couldn't parse db data to dto struct, err: %v", err)
		return dto.GetUserDTO{}, err
	}
	return userDTO, nil
}
