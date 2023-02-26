package chat

type Message struct {
	Id        int    `json:"messageId"`
	UserId    string `json:"userId"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
