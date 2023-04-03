package userdto

type UserDTO struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Mail     string `json:"mail"`
	Password string `json:"pass"`
}

type CreateUserDTO struct {
	Name     string `json:"name"`
	Mail     string `json:"mail"`
	Password string `json:"pass"`
}

type GetUserDTO struct {
	Id       string `json:"id"`
	Mail     string `json:"mail"`
	Password string `json:"pass"`
}

type ChatUserDTO struct {
	Name string `json:"name"`
	Mail string `json:"mail"`
}
