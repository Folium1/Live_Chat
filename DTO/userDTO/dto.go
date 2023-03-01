package userdto

type CreateUserDTO struct {
	Name     string `json:"name"`
	Mail     string `json:"mail"`
	Password string `json:"pass"`
	Online   bool   `json:"online"`
}
