package retry

import (
	"context"
	"math"
	"math/rand"
	"time"
)

type Action int

const (
	Fail Action = iota
	Succeed
	Retry
)

const Stop time.Duration = -1

type FnBackoff func(AttemptNum int, min, max time.Duration) time.Duration

type Worker func(ctx context.Context) error

type RPolicy func(err error) Action

type Backoff struct {
	min, max               time.Duration
	attemptNum, maxAttempt int
	backoff                FnBackoff
}

type Retrier struct {
	backoff     *Backoff
	retryPolicy RPolicy
}

func (r *Retrier) sleep(ctx context.Context, t <-chan time.Time) error {
	select {
	case <-t:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (r *Retrier) Run(ctx context.Context, worker Worker) error {
	defer r.backoff.Reset()
	for {
		err := worker(ctx)
		switch r.retryPolicy(err) {
		case Retry:
			var delay time.Duration
			if delay = r.backoff.Next(); delay == Stop {
				return err
			}
			timeout := time.After(delay)
			if err := r.sleep(ctx, timeout); err != nil {
				return err
			}
		case Succeed, Fail:
			return err
		}
	}
}

func NewRetrier(backoff *Backoff, retryPolicy RPolicy) *Retrier {
	if retryPolicy == nil {
		retryPolicy = DefaultRetryPolicy
	}
	return &Retrier{backoff: backoff, retryPolicy: retryPolicy}
}

func (b *Backoff) Reset() {
	b.attemptNum = 0
}

func (b *Backoff) Next() time.Duration {
	if b.attemptNum >= b.maxAttempt {
		return Stop
	}
	b.attemptNum++
	return b.backoff(b.attemptNum, b.min, b.max)
}

func NewBackoff(min, max time.Duration, maxAttempt int, backoff FnBackoff) *Backoff {
	if backoff == nil {
		backoff = ExponentialBackoff
	}
	return &Backoff{min: min, max: max, maxAttempt: maxAttempt, backoff: backoff}
}

func ConstantBackoff(factor time.Duration) FnBackoff {
	return func(AttemptNum int, min, max time.Duration) time.Duration {
		if factor < min {
			return min
		}
		if factor > max {
			return max
		}
		return factor
	}
}

func LinerBackoff(factor time.Duration) FnBackoff {
	return func(AttemptNum int, min, max time.Duration) time.Duration {
		delay := factor * time.Duration(AttemptNum)
		if delay < min {
			delay = min
		}
		jitter := delay * time.Duration(rand.Float64()*float64(delay-min))

		delay = delay + jitter
		if delay > max {
			delay = max
		}
		return delay
	}
}

func ExponentialBackoff(AttemptNum int, min, max time.Duration) time.Duration {
	factor := 2.0
	rand.NewSource(time.Now().UnixNano())
	delay := time.Duration(math.Pow(factor, float64(AttemptNum)) * float64(min))
	jitter := time.Duration(rand.Float64() * float64(AttemptNum) * float64(min))

	delay = delay + jitter
	if delay > max {
		delay = max
	}
	return delay
}

func DefaultRetryPolicy(err error) Action {
	if err != nil {
		return Retry
	}
	return Succeed
}
