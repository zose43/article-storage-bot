package telegram

import "time"

type Update struct {
	Id      int     `json:"update_id"`
	Message Message `json:"message"`
}

type Message struct {
	Id   int       `json:"message_id"`
	Date time.Time `json:"date"`
	Text string    `json:"text"`
}

type UpdateResponse struct {
	Issue   bool     `json:"ok"`
	Results []Update `json:"result"`
}
