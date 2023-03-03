package chat

type Message struct {
	Id        int    `json:"messageId"`
	UserName  string `json:"user_name"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
