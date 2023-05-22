package events

type Fetcher interface {
	Fetch(limit int) ([]Event, error)
}

type Processor interface {
	Processor(e Event) error
}

const (
	Unknown Type = iota
	Message
)

type Type int

type Event struct {
	Type Type
	Text string
	Meta any
}
