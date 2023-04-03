package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	entity "chat/entities"
)

var ctx = context.Background()

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
	redisDb, err := entity.RedisConnect()
	if err != nil {
		log.Fatalf("Couldn't connect to Redis, err: '%v'", err)
	}

	id, err := redisDb.Incr(ctx, "chat:nextMsgId").Result()
	if err != nil {
		log.Fatalf("Couldn't generate new message ID, err: '%v'", err)
		return 0, err
	}

	msg.Id = int(id)

	data, err := json.Marshal(msg)
	if err != nil {
		log.Fatalf("Couldn't marshal message to JSON, err: '%v'", err)
		return 0, err
	}

	key := fmt.Sprintf("chat:msg:%d", msg.Id)
	err = redisDb.Set(ctx, key, data, 0).Err()
	if err != nil {
		log.Fatalf("Couldn't save message to Redis, err: '%v'", err)
		return 0, err
	}

	return msg.Id, nil
}

// DeleteMsg deletes message from Redis
func (m *Message) DeleteMsg(id string) error {
	redisDb, err := entity.RedisConnect()
	if err != nil {
		log.Fatalf("Couldn't connect to Redis, err: '%v'", err)
	}

	key := fmt.Sprintf("chat:msg:%s", id)
	err = redisDb.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}

// GetAllMessages returns all messages from Redis
func (m *Message) GetAllMessages() ([]Message, error) {
	redisDb, err := entity.RedisConnect()
	if err != nil {
		log.Fatalf("Couldn't connect to Redis, err: '%v'", err)
	}

	iter := redisDb.Scan(ctx, 0, "chat:msg:*", 0).Iterator()

	var messages []Message

	for iter.Next(ctx) {
		key := iter.Val()
		data, err := redisDb.Get(ctx, key).Bytes()
		if err != nil {
			log.Fatalf("Couldn't get message data from Redis, err: '%v'", err)
			return nil, err
		}
		var msg Message
		err = json.Unmarshal(data, &msg)
		if err != nil {
			log.Fatalf("Couldn't unmarshal message from JSON, err: '%v'", err)
			return nil, err
		}

		messages = append(messages, msg)
	}

	if err := iter.Err(); err != nil {
		log.Fatalf("Couldn't iterate over Redis keys, err: '%v'", err)
		return nil, err
	}

	return messages, nil
}
