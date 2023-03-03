package userdto

type UserDTO struct {
	Id       int    `json:"id"`
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
	Id       int    `json:"id"`
	Mail     string `json:"mail"`
	Password string `json:"pass"`
}
