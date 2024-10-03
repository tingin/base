package retry

import (
	"time"

	"github.com/avast/retry-go"
)

type RetryFunc func() error

func Do(fn RetryFunc, attempts int, delay time.Duration) error {
	return retry.Do(
		func() error {
			return fn()
		},
		retry.Attempts(uint(attempts)),
		retry.Delay(delay),
		retry.DelayType(retry.FixedDelay),
	)
}
