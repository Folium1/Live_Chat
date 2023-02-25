package entity

import "time"

type Message struct {
	Id        int       `json:"messageId"`
	UserId    string    `json:"userId"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
