package chat

import (
	"fmt"
	"log"

	entity "chat/entities"
)

type MessageService interface {
	SendMsg(msg Message) (int, error)
	DeleteMsg(id string) error
	GetAllMessages() ([]Message, error)
}

func New() MessageService {
	return &Message{}
}

// SendMsg creates message
func (m *Message) SendMsg(msg Message) (int, error) {
	db, err := entity.DbConnect()
	if err != nil {
		log.Fatalf("Couldn't connect to db, err: '%v'", err)
	}
	defer db.Close()
	query := fmt.Sprintf("INSERT INTO chat.messages(user_name,user_id,text,created_at,updated_at) VALUES('%v','%v','%v',now(),now());", msg.UserName, msg.UserId, msg.Text)
	result, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Err: '%v'", err)
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Fatalf("Err: '%v'", err)
		return 0, err
	}
	return int(id), nil
}

// DeleteMsg deletes message from db
func (m *Message) DeleteMsg(id string) error {
	db, err := entity.DbConnect()
	if err != nil {
		log.Fatalf("Couldn't connect to db, err: '%v'", err)
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM chat.messages WHERE id = ?;", id)
	if err != nil {
		return err
	}
	return nil
}

// GetAllMessages returns all messages from db
func (m *Message) GetAllMessages() ([]Message, error) {
	db, err := entity.DbConnect()
	if err != nil {
		log.Fatalf("Couldn't connect to db, err: '%v'", err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT id,text,user_name,created_at,updated_at FROM chat.messages;")
	if err != nil {
		log.Panicf("Couldn't make a query, err: '%v'", err)
	}
	defer rows.Close()
	var messages []Message
	for rows.Next() {
		var currentMessage Message
		err = rows.Scan(&currentMessage.Id, &currentMessage.Text, &currentMessage.UserName, &currentMessage.CreatedAt, &currentMessage.UpdatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, currentMessage)
	}
	return messages, nil
}
