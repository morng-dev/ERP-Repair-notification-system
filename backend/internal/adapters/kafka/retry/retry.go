package retry

import (
	"context"
	"fmt"
	"time"
)

type Config struct {
	MaxRetries   int
	InitialDelay time.Duration
	MaxDelay     time.Duration
}

func Do(ctx context.Context, cfg Config, fn func() error) error {

	var err error

	delay := cfg.InitialDelay

	for attempt := 0; attempt <= cfg.MaxRetries; attempt++ {

		if ctx.Err() != nil {
			return ctx.Err()
		}

		err = fn()

		if err == nil {
			return nil
		}

		if attempt == cfg.MaxRetries {
			break
		}

		timer := time.NewTimer(delay)

		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()

		case <-timer.C:
		}

		delay *= 2

		if delay > cfg.MaxDelay {
			delay = cfg.MaxDelay
		}
	}

	return fmt.Errorf("operation failed after %d retries: %w", cfg.MaxRetries, err)
}
