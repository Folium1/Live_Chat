package redis_jwt

type UserJwt struct {
	Id    string `json:"user_id"`
	Token string `json:"jwt"`
}
