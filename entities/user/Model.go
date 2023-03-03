package user

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"user_name"`
	Mail     string `json:"mail"`
	Password string `json:"pass"`
}
