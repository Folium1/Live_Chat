package chatdto

type CreateMessageDTO struct {
	Id       int    `json:"messageId"`
	UserName string `json:"user_name"`
	UserId   string `json:"userId"`
	Text     string `json:"text"`
}

type UpdateMessageDTO struct {
	Id   int    `json:"messageId"`
	Text string `json:"text"`
}

type MessageDTO struct {
	Id        int    `json:"messageId"`
	UserName  string `json:"user_name"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type MessagesDTO struct {
	Id        int    `json:"messageId"`
	Text      string `json:"text"`
	UserName  string `json:"user_name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
