package user

type User struct {
	Id       int `json:"id"`
	Name     string `json:"name"`
	Mail     string `json:"mail"`
	Password string `json:"pass"`
}
