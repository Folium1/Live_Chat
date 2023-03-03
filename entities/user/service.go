package user

import (
	entity "chat/entities"
	"fmt"
	"log"
)

const userTable = "users"

type UserService interface {
	CreateUser(newUser User) error
	GetUser(user User) (User, error)
}

func New() UserService {
	return &User{}
}

// Creats new User
func (u *User) CreateUser(newUser User) error {
	db, err := entity.DbConnect()
	if err != nil {
		log.Panicf("Couldn't connect to db, err: %v", err)
	}
	defer db.Close()
	query := fmt.Sprintf("INSERT INTO chat.users(name,email,password) VALUES('%v','%v','%v');", newUser.Name, newUser.Mail, newUser.Password)
	_, err = db.Query(query)
	if err != nil {
		log.Panicf("Couldn't create new user, err: %v", err)
		return err
	}
	return nil
}

// Gets user by password and mail
func (u *User) GetUser(user User) (User, error) {
	db, err := entity.DbConnect()
	if err != nil {
		log.Panicf("Couldn't connect to db, err: %v", err)
	}
	defer db.Close()
	query := fmt.Sprintf("SELECT id,password FROM chat.users WHERE email = '%v';", user.Mail)
	q, err := db.Query(query)
	if err != nil {
		err = fmt.Errorf("Couldn't get user, err: %v", err)
		log.Printf(err.Error())
		return User{}, err
	}
	for q.Next() {
		err := q.Scan(&user.Id, &user.Password)
		if err != nil {
			log.Printf("Couldn't parse user's id, err: %v", err)
		}
	}
	return user, nil

}
