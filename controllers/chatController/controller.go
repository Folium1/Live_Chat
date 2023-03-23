package chatController

import (
	chatdto "chat/DTO/chatDTO"
	"chat/entities/chat"
	"encoding/json"
)

type chatController struct {
	db chat.MessageService
}

func New(msg chat.MessageService) ControllerInterface {
	return &chatController{msg}
}

type ControllerInterface interface {
	CreateMsg(msg chatdto.CreateMessageDTO) error
	ChangeData(newMessageData chatdto.UpdateMessageDTO) error
	DeleteMsg(id string) error
	GetAllMessages() ([]chatdto.MessagesDTO, error)
}

func (c *chatController) CreateMsg(msg chatdto.CreateMessageDTO) error {
	var message chat.Message
	// Parsing dto data to Message struct variable
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &message)
	if err != nil {
		return err
	}
	err = c.db.SendMsg(message)
	if err != nil {
		return err
	}
	return nil
}

// Gets new text for message by id
func (c *chatController) ChangeData(newData chatdto.UpdateMessageDTO) error {
	var changedMessage chat.Message
	// Parsing dto data to Message struct variable
	data, err := json.Marshal(newData)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &changedMessage)
	if err != nil {
		return err
	}
	err = c.db.EditeMsg(changedMessage)
	if err != nil {
		return err
	}
	return nil
}

// Delets message from db by id
func (c *chatController) DeleteMsg(id string) error {
	err := c.db.DeleteMsg(id)
	if err != nil {
		return err
	}
	return nil
}

func (c chatController) GetAllMessages() ([]chatdto.MessagesDTO, error) {
	messages, err := c.db.GetAllMessages()
	if err != nil {
		return []chatdto.MessagesDTO{}, err
	}
	messagesDTO := make([]chatdto.MessagesDTO, 0)
	var messageDTO chatdto.MessagesDTO
	for _, dbMessage := range messages {
		err = parseToDTO(dbMessage, &messageDTO)
		if err != nil {
			return []chatdto.MessagesDTO{}, err
		}
		messagesDTO = append(messagesDTO, messageDTO)
	}
	return messagesDTO, nil
}
