package main

import (
	"article-storage-bot/clients/telegram"
	event_consumer "article-storage-bot/consumer/event-consumer"
	manager "article-storage-bot/events/telegram"
	"article-storage-bot/storage/files"
	"flag"
	"log"
)

const BatchSize = 10

func main() {
	t, h := mustToken()
	tgClient := telegram.NewClient(*t, *h)
	fetcher := manager.NewManager(&tgClient, files.NewStorage(""))
	consumer := event_consumer.NewConsumer(fetcher, fetcher, BatchSize)

	log.Print("service started")
	log.Fatal(consumer.Handle())
}

func mustToken() (t, h *string) {
	t = flag.String(
		"t",
		"",
		"token for access to telegram bot api",
	)
	h = flag.String(
		"h",
		"api.telegram.org",
		"telegram api host")
	flag.Parse()

	if *t == "" {
		log.Fatal("No one token for telegram bot")
	}
	return
}
