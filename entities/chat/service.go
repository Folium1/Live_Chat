package chat

import (
	"fmt"
	"log"
	"time"

	entity "chat/entities"
)

type MessageService interface {
	SendMsg(msg Message) error
	EditeMsg(newMessageData Message) error
	DeleteMsg(id string) error
	GetAllMessages() ([]Message, error)
}

func New() MessageService {
	return &Message{}
}

// SendMsg creates message
func (m *Message) SendMsg(msg Message) error {
	db, err := entity.DbConnect()
	if err != nil {
		log.Fatalf("Couldn't connect to db, err: '%v'", err)
	}
	defer db.Close()
	// Setting up local time for sent message
	msg.CreatedAt = time.Now().Format("2006-01-02 15:04")
	msg.UpdatedAt = time.Now().Format("2006-01-02 15:04")
	query := fmt.Sprintf("INSERT INTO messages(user_name,user_id,text,created_at,updated_at) VALUES('%v','%v','%v','%v','%v');", msg.UserName, msg.UserId, msg.Text, msg.CreatedAt, msg.UpdatedAt)
	_, err = db.Query(query)
	if err != nil {
		log.Fatalf("Err: '%v'", err)
		return err
	}
	return nil
}

// EditMsg changes text and updated_at fields of the msg
func (m *Message) EditeMsg(newMessageData Message) error {
	db, err := entity.DbConnect()
	if err != nil {
		log.Fatalf("Couldn't connect to db, err: '%v'", err)
	}
	defer db.Close()
	// Changing UpdatedAt value to current time
	newMessageData.UpdatedAt = time.Now().Format("2006-01-02 15:04")
	query := fmt.Sprintf("UPDATE messages SET text = '%v', updated_at = '%v' WHERE id = '%v';", newMessageData.Text, newMessageData.UpdatedAt, newMessageData.Id)
	_, err = db.Query(query)
	if err != nil {
		return err
	}
	return nil
}

// DeleteMsg deletes message from db
func (m *Message) DeleteMsg(id string) error {
	db, err := entity.DbConnect()
	if err != nil {
		log.Fatalf("Couldn't connect to db, err: '%v'", err)
	}
	defer db.Close()

	_, err = db.Query("DELETE * FROM messages WHERE id = '%v';", id)
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
	q, err := db.Query("SELECT text,user_name,created_at,updated_at FROM messages ORDER BY -created_at;")
	if err != nil {
		log.Panicf("Couldn't make a query, err: '%v'", err)
	}
	var messages []Message
	for q.Next() {
		var currentMessage Message
		err = q.Scan(&currentMessage.Text, &currentMessage.UserName, &currentMessage.CreatedAt, &currentMessage.UpdatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, currentMessage)
	}
	return messages, nil
}
