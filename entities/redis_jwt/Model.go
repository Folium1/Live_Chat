package redis_jwt

type RedisData struct {
	Id    string `json:"user_id"`
	Token string `json:"jwt"`
}
