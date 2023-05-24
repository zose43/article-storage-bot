package consumer

type Consumer interface {
	Handle() error
}
