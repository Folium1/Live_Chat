package redis_jwt

type RedisData struct {
	Id    string `json:"user_id"`
	Key   string `json:"key"`
	Token string `json:"jwt"`
}
