package entity

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Mail     string `json:"mail"`
	Password string `json:"pass"`
	Online   bool   `json:"online"`
}
