package main

import (
	"flag"
	"log"
)

func main() {
	t := mustToken()
	// tgClient := telegram.New(t)
	// fetcher:= fetcher.New(t)
	// processor:= processor.New(t)
	// consumer.Start(processor,fetcher)
}

func mustToken() string {
	t := flag.String(
		"t",
		"",
		"token for access to telegram bot api",
	)
	flag.Parse()
	if *t == "" {
		log.Fatal("No one token for telegram bot")
	}
	return *t
}
