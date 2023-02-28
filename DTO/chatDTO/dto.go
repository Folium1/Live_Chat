package chatdto

type CreateMessageDTO struct {
	UserId string `json:"userId"`
	Text   string `json:"text"`
}

type UpdateMessageDTO struct {
	Id   int    `json:"messageId"`
	Text string `json:"text"`
}

type MessageDTO struct {
	Id        int    `json:"messageId"`
	UserId    string `json:"userId"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
