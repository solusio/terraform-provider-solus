// Copyright 1999-2022. Plesk International GmbH. All rights reserved.

package timer

import (
	"context"
	"time"
)

// WaitFor waits until `fn` return `true, nil`, or `false, error`.
func WaitFor(ctx context.Context, d time.Duration, fn func() (bool, error)) error {
	t := time.NewTicker(d)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
			isFinished, err := fn()
			if isFinished || err != nil {
				return err
			}
		}
	}
}
