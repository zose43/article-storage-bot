package telegram

import "time"

type Update struct {
	Id      int              `json:"update_id"`
	Message *IncomingMessage `json:"message"`
}

type IncomingMessage struct {
	Id   int           `json:"message_id"`
	Date time.Duration `json:"date"`
	Text string        `json:"text"`
	User User          `json:"from"`
	Chat Chat          `json:"chat"`
}

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"first_name"`
	Username string `json:"username"`
}

type Chat struct {
	Id int `json:"id"`
}

type UpdateResponse struct {
	Issue   bool     `json:"ok"`
	Results []Update `json:"result"`
}
