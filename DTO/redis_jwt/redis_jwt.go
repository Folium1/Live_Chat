package redisjwt

type RdbDTO struct {
	Id    string `json:"user_id"`
	Token string `json:"jwt"`
}
