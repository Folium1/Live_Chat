package user

import (
	"chat/entity"
	"fmt"
)

const userTable = "users"

type UserService interface {
	CreateUser(newUser user) error
}

func New() UserService {
	return &user{}
}

// Creats new user
func (u *user) CreateUser(newUser user) error {
	db, err := entity.DbConnect(userTable)
	if err != nil {
		return err
	}
	defer db.Close()
	query := fmt.Sprintf("INSERT INTO users(name,mail,password,online) VALUES(%v,%v,%v,%v);", newUser.Name, newUser.Mail, newUser.Password, newUser.Online)
	_, err = db.Query(query)
	if err != nil {
		return err
	}
	return nil
}
