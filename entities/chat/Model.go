package chat

type Message struct {
	Id        int    `json:"messageId"`
	UserId    string `json:"userId"`
	UserName  string `json:"user_name"`
	Text      string `json:"text"`
}
