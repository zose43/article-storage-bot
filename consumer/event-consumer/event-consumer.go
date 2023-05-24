package event_consumer

import (
	"article-storage-bot/events"
	"log"
	"time"
)

type Consumer struct {
	events.Fetcher
	events.Processor
	batchSize int
}

func (c Consumer) Handle() error {
	// TODO fallback list(retry), parallel handler, counter consistent errors
	log.SetPrefix("error: ")
	for {
		gotEvents, err := c.Fetch(c.batchSize)
		if err != nil {
			log.Printf("can't handle event %s", err)
			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		for _, event := range gotEvents {
			err := c.Processor.Processor(event)
			if err != nil {
				log.Printf("can't handle event %s", err)
				continue
			}
		}
	}
}

func NewConsumer(fetcher events.Fetcher, processor events.Processor, batchSize int) *Consumer {
	return &Consumer{
		Fetcher:   fetcher,
		Processor: processor,
		batchSize: batchSize}
}
