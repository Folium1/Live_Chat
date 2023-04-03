package user

import (
	"fmt"
	"log"
	"strconv"

	entity "chat/entities"
)

type UserService interface {
	CreateUser(newUser User) (string, error)
	UserByMail(user User) (User, error)
	GetUserById(id string) (User, error)
}

func New() UserService {
	return &User{}
}

// CreateUser Creats new User
func (u *User) CreateUser(newUser User) (string, error) {
	db, err := entity.MySQLConnect()
	if err != nil {
		log.Panicf("Couldn't connect to db, err: %v", err)
	}
	defer db.Close()
	query := fmt.Sprintf("INSERT INTO chat.users(name,email,password) VALUES('%v','%v','%v');", newUser.Name, newUser.Mail, newUser.Password)
	result, err := db.Exec(query)
	if err != nil {
		log.Panicf("Couldn't create new user, err: %v", err)
		return "", err
	}
	id, err := result.LastInsertId()
	return strconv.Itoa(int(id)), nil
}

// UserByMail returns user by his mail
func (u *User) UserByMail(user User) (User, error) {
	db, err := entity.MySQLConnect()
	if err != nil {
		log.Panicf("Couldn't connect to db, err: %v", err)
	}
	defer db.Close()
	query := fmt.Sprintf("SELECT id,password FROM chat.users WHERE email = '%v';", user.Mail)
	raws, err := db.Query(query)
	if err != nil {
		err = fmt.Errorf("Couldn't get user, err: %v", err)
		log.Println(err)
		return User{}, err
	}
	defer raws.Close()
	for raws.Next() {
		err := raws.Scan(&user.Id, &user.Password)
		if err != nil {
			log.Printf("Couldn't parse user's id, err: %v", err)
		}
	}
	return user, nil

}

// GetUserById returns user by his шd
func (u *User) GetUserById(id string) (User, error) {
	db, err := entity.MySQLConnect()
	if err != nil {
		log.Panicf("Couldn't connect to db, err: %v", err)
	}
	defer db.Close()
	query := fmt.Sprintf("SELECT name,email FROM chat.users WHERE id = %v;", id)
	raws, err := db.Query(query)
	if err != nil {
		log.Printf("No user with id %v", id)
		return User{}, err
	}
	var user User
	user.Id = id
	defer raws.Close()
	for raws.Next() {
		err = raws.Scan(&user.Name, &user.Mail)
		if err != nil {
			log.Printf("Couldn't parse data to user struct, err:%v", err)
			return User{}, err
		}
	}
	return user, err
}
