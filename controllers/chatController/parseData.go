package chatController

import (
	dto "chat/DTO/chatDTO"
	"chat/entities/chat"
	"encoding/json"
)

func parseToDTO[T dto.MessagesDTO | dto.MessageDTO](message chat.Message, dtoData *T) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &dtoData)
	if err != nil {
		return err
	}
	return nil
}
