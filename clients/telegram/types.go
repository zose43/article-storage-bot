package telegram

import "time"

type Update struct {
	Id      int              `json:"update_id"`
	Message *IncomingMessage `json:"message"`
}

type IncomingMessage struct {
	Id   int       `json:"message_id"`
	Date time.Time `json:"date"`
	Text string    `json:"text"`
	User `json:"from"`
	Chat `json:"chat"`
}

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"first_name"`
	LastName string `json:"last_name"`
	Username string `json:"username"`
	Language string `json:"language_code"`
}

type Chat struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

type UpdateResponse struct {
	Issue   bool     `json:"ok"`
	Results []Update `json:"result"`
}
