package chatController

import (
	chatdto "chat/DTO/chatDTO"
	"chat/entity/chat"

	"github.com/mitchellh/mapstructure"
)

type chatController struct {
	db chat.Message
}

func New(msg chat.Message) ControllerInterface {
	return &chatController{msg}
}

type ControllerInterface interface {
	CreateMsg(msg chatdto.CreateMessageDTO) error
	ChangeData(newMessageData chatdto.UpdateMessageDTO) error
	DeleteMsg(id string) error
}

func (c *chatController) CreateMsg(msg chatdto.CreateMessageDTO) error {
	var message chat.Message
	err := mapstructure.Decode(msg, &message)
	if err != nil {
		return err
	}
	err = c.db.CreateMsg(message)
	if err != nil {
		return err
	}
	return nil
}

func (c *chatController) ChangeData(newData chatdto.UpdateMessageDTO) error {
	var changedMessage chat.Message
	err := mapstructure.Decode(newData, &changedMessage)
	if err != nil {
		return err
	}
	err = c.db.ChangemMsg(changedMessage)
	if err != nil {
		return err
	}
	return nil
}

func (c *chatController) DeleteMsg(id string) error {
	err := c.db.DeleteMsg(id)
	if err != nil {
		return err
	}
	return nil
}
