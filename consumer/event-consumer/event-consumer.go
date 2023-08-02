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
	// TODO fallback list(retry, counter consistent errors
	log.SetPrefix("error: ")
	ch := make(chan []events.Event, 25)
	go func() {
		for {
			select {
			case collection := <-ch:
				for _, event := range collection {
					go func(event events.Event) {
						err := c.Processor.Processor(event)
						if err != nil {
							log.Printf("can't process event %s", err)
						}
					}(event)
				}
			}
		}
	}()
	for {
		time.Sleep(500 * time.Millisecond)
		go func() {
			gotEvents, err := c.Fetch(c.batchSize)
			if err != nil {
				log.Printf("can't fetch event %s", err)
			}
			if len(gotEvents) > 0 {
				ch <- gotEvents
			}
		}()
	}
}

func NewConsumer(fetcher events.Fetcher, processor events.Processor, batchSize int) *Consumer {
	return &Consumer{
		Fetcher:   fetcher,
		Processor: processor,
		batchSize: batchSize}
}
