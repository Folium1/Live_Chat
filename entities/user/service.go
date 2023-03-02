package user

import (
	entity "chat/entities"
	"fmt"
	"log"
)

const userTable = "users"

type UserService interface {
	CreateUser(newUser User) error
	GetUser(user User) (int, error)
}

func New() UserService {
	return &User{}
}

// Creats new User
func (u *User) CreateUser(newUser User) error {
	db, err := entity.DbConnect(userTable)
	if err != nil {
		log.Panicf("Couldn't connect to db, err: %v", err)
	}
	defer db.Close()
	query := fmt.Sprintf("INSERT INTO users(name,mail,password,online) VALUES(%v,%v,%v);", newUser.Name, newUser.Mail, newUser.Password)
	_, err = db.Query(query)
	if err != nil {
		return err
	}
	return nil
}

// Gets user by password and mail
func (u *User) GetUser(user User) (int, error) {
	db, err := entity.DbConnect(userTable)
	if err != nil {
		log.Panicf("Couldn't connect to db, err: %v", err)
	}
	defer db.Close()
	query := fmt.Sprintf("SELECT id WHERE password = %v, email = %v;", user.Password, user.Mail)
	q, err := db.Query(query)
	if err != nil {
		err = fmt.Errorf("Couldn't get user, err: %v", err)
		log.Printf(err.Error())
		return -1, err
	}

	var id int
	for q.Next() {
		err := q.Scan(&id)
		if err != nil {
			log.Printf("Couldn't parse user's id, err: %v", err)
		}
	}
	return id, nil

}
