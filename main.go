package main

import (
	"article-storage-bot/clients/telegram"
	"flag"
	"log"
)

func main() {
	t, h := mustToken()
	tgClient := telegram.NewClient(*t, *h)
	// fetcher:= fetcher.New(t)
	// processor:= processor.New(t)
	// consumer.Start(processor,fetcher)
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
