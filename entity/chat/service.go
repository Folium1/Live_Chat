package chat

import (
	"chat/entity"
	"errors"
	"fmt"
	"time"
)

const messageTable string = "message"

type MessageService interface {
	Create(msg Message) error
	GetMessage(id int) (Message, error)
	ChangeData(newMessage Message) error
}

func New() MessageService {
	return &Message{}
}

func (m *Message) Create(msg Message) error {
	db, err := entity.DbConnect(messageTable)
	if err != nil {
		return err
	}
	defer db.Close()
	msg.CreatedAt = time.Now().Format("YYYY-MM-DD hh:mm")
	msg.UpdatedAt = time.Now().Format("YYYY-MM-DD hh:mm")
	query := fmt.Sprintf("INSERT INTO message(user_id,text,created_at,updated_at) VALUES(%v,%v,%v,%v);", msg.UserId, msg.Text, msg.CreatedAt, msg.UpdatedAt)
	_, err = db.Query(query)
	if err != nil {
		return err
	}
	return nil
}

func (m *Message) GetMessage(id int) (Message, error) {
	db, err := entity.DbConnect(messageTable)
	if err != nil {
		return Message{}, err
	}
	defer db.Close()
	query := fmt.Sprintf("SELECT user_id,text,created_at,updated_at FROM message WHERE id = %v;", id)
	res, err := db.Query(query)
	if err != nil {
		return Message{}, errors.New("message id doesn't exists")
	}
	var resMessage Message
	resMessage.Id = id
	for res.Next() {
		err = res.Scan(&resMessage.UserId, &resMessage.Text, &resMessage.CreatedAt, &resMessage.UpdatedAt)
		if err != nil {
			return Message{}, err
		}
	}
	return resMessage, nil
}

func (m *Message) ChangeData(newMessage Message) error {
	db, err := entity.DbConnect(messageTable)
	if err != nil {
		return err
	}
	defer db.Close()
	newMessage.UpdatedAt = time.Now().Format("YYYY-MM-DD hh:mm")
	query := fmt.Sprintf("UPDATE message SET text = %v, updated_at = %v WHERE id = %v;", newMessage.Text, newMessage.UpdatedAt, newMessage.Id)
	_, err = db.Query(query)
	if err != nil {
		return err
	}
	return nil
}
