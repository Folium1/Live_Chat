package usersController

import (
	dto "chat/DTO/userdto"
	"chat/entities/user"
	"log"
)

type userController struct {
	db user.UserService
}

func New(user user.UserService) UserController {
	return &userController{user}
}

type UserController interface {
	CreateUser(newUser dto.CreateUserDTO) error
	GetUser(userData dto.GetUserDTO) (dto.GetUserDTO, error)
	GetUserById(id int) (dto.ChatUserDTO, error)
}

func (c *userController) CreateUser(newUser dto.CreateUserDTO) error {
	var dbUser user.User
	err := parseUserToDb(newUser, &dbUser)
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
	err := parseUserToDb(userData, &dbUser)
	if err != nil {
		log.Printf("couldn't parse dto data to db struct, err: %v", err)
	}
	userFromDb, err := c.db.UserByMail(dbUser)
	if err != nil {
		log.Printf("couldn't get user, err: %v", err)
		return dto.GetUserDTO{}, err
	}
	var userDTO dto.GetUserDTO
	err = parseUserToDTO(userFromDb, &userDTO)
	if err != nil {
		log.Printf("couldn't parse db data to dto struct, err: %v", err)
		return dto.GetUserDTO{}, err
	}
	return userDTO, nil
}

func (c *userController) GetUserById(id int) (dto.ChatUserDTO, error) {
	user, err := c.db.GetUserById(id)
	if err != nil {

		return dto.ChatUserDTO{}, err
	}
	var userDTO dto.ChatUserDTO
	err = parseUserToDTO(user, &userDTO)
	if err != nil {
		log.Printf("Couldn't parse db data to dto, err: %v", err)
		return dto.ChatUserDTO{}, err
	}
	return userDTO, nil
}
