package chat

import (
	entity "chat/entities"
	"fmt"
	"log"
	"time"
)

const messageTable = "messages"

type MessageService interface {
	SendMsg(msg Message) error
	EditmMsg(newMessageData Message) error
	DeleteMsg(id string) error
	GetAllMessages() ([]Message, error)
}

func New() MessageService {
	return &Message{}
}

// Creates message
func (m *Message) SendMsg(msg Message) error {
	db, err := entity.DbConnect()
	if err != nil {
		log.Panicf("Couldn't connect to db, err: %v", err)
	}
	defer db.Close()
	// Setting up local time for sent message
	msg.CreatedAt = time.Now().Format("2006-01-02 15:04")
	msg.UpdatedAt = time.Now().Format("2006-01-02 15:04")
	query := fmt.Sprintf("INSERT INTO messages(user_name,text,created_at,updated_at) VALUES(%v,%v,%v,%v);", msg.UserName, msg.Text, msg.CreatedAt, msg.UpdatedAt)
	_, err = db.Query(query)
	if err != nil {
		return err
	}
	return nil
}

// Changes text of the msg
func (m *Message) EditmMsg(newMessageData Message) error {
	db, err := entity.DbConnect()
	if err != nil {
		log.Panicf("Couldn't connect to db, err: %v", err)
	}
	defer db.Close()
	// Changing UpdatedAt value to current time
	newMessageData.UpdatedAt = time.Now().Format("2006-01-02 15:04")
	query := fmt.Sprintf("UPDATE messages SET text = %v, updated_at = %v WHERE id = %v;", newMessageData.Text, newMessageData.UpdatedAt, newMessageData.Id)
	_, err = db.Query(query)
	if err != nil {
		return err
	}
	return nil
}

// Deletes message from db
func (m *Message) DeleteMsg(id string) error {
	db, err := entity.DbConnect()
	if err != nil {
		log.Panicf("Couldn't connect to db, err: %v", err)
	}
	defer db.Close()

	_, err = db.Query("DELETE * FROM messages WHERE id = %v;", id)
	if err != nil {
		return err
	}
	return nil
}

func (m *Message) GetAllMessages() ([]Message, error) {
	db, err := entity.DbConnect()
	if err != nil {
		log.Panicf("Couldn't connect to db, err: %v", err)
	}
	defer db.Close()
	query := fmt.Sprintf("SELECT text,created_at,updated_at FROM messages ORDER BY -created_at;")
	q, err := db.Query(query)
	if err != nil {
		log.Panicf("Couldn't make a query, err: %v", err)
	}
	var messages []Message
	for q.Next() {
		var currentMessage Message
		err = q.Scan(&currentMessage.Text, currentMessage.CreatedAt, currentMessage.UpdatedAt)
		if err != nil {

			return nil, err
		}
		messages = append(messages, currentMessage)
	}
	return messages, nil
}
