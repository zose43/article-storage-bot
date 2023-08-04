package event_consumer

import (
	"article-storage-bot/events"
	"article-storage-bot/lib/retry"
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

type Consumer struct {
	events.Fetcher
	events.Processor
	batchSize int
}

var Shutdown = fmt.Errorf("too many errors, Shutdown")

type errCounter struct {
	errNums int
	ticker  time.Ticker
}

func (c Consumer) Handle() error {
	log.SetPrefix("error: ")
	ch := make(chan []events.Event, 25)

	backoff := retry.NewBackoff(100*time.Millisecond, 3*time.Second, 3, nil)
	ctx, cancel := context.WithCancel(context.Background())
	r := retry.NewRetrier(backoff, nil)
	ec := errCounter{ticker: *(time.NewTicker(30 * time.Second))}
	go func() {
		defer log.Print("close ch")
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
			case <-ec.ticker.C:
				if ec.errNums > 10 {
					cancel()
					close(ch)
					return
				}
			}
		}
	}()
	for {
		time.Sleep(500 * time.Millisecond)
		go func() {
			err := r.Run(ctx, func(ctx context.Context) error {
				select {
				case <-ctx.Done():
					return Shutdown
				default:
					gotEvents, err := c.Fetch(c.batchSize)
					if err != nil {
						return err
					}
					if len(gotEvents) > 0 {
						ch <- gotEvents
					}
					return nil
				}
			})
			if err != nil {
				ec.errNums++
				log.Printf("can't fetch event %s", err)
			}
			if errors.Is(err, Shutdown) {
				log.Fatal(err)
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
