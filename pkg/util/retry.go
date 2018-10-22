//
// Copyright © 2018 Aljabr, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package util

import (
	"context"
	"time"

	"github.com/cenkalti/backoff"
)

// Permanent flags the given error as a permanent failure.
func Permanent(err error) error {
	return backoff.Permanent(err)
}

type causer interface {
	Cause() error
}

// Retry the given operation until:
// - Success or
// - Given context is canceled or
// - A permanent error is returns from the given operation
func Retry(ctx context.Context, op func(context.Context) error, minAttempts ...int) error {
	if ctx == nil {
		ctx = context.Background()
	}
	divider := 1
	if len(minAttempts) > 0 {
		divider = minAttempts[0]
	}
	b := backoff.NewExponentialBackOff()
	if deadline, ok := ctx.Deadline(); ok {
		b.MaxElapsedTime = time.Until(deadline)
	}
	b.InitialInterval = b.MaxElapsedTime / time.Duration(divider*10)
	b.MaxInterval = b.MaxElapsedTime / time.Duration(divider)

	if err := backoff.Retry(func() error {
		lctx, cancel := context.WithTimeout(ctx, b.MaxInterval)
		defer cancel()
		go func() {
			select {
			case <-time.After(b.MaxInterval):
				cancel()
			case <-lctx.Done():
			}
		}()
		err := op(lctx)
		if err == nil {
			if lctx.Err() != nil {
				err = lctx.Err()
			} else {
				return nil
			}
		}
		for {
			if perr, ok := err.(*backoff.PermanentError); ok {
				return perr
			}
			if c, ok := err.(causer); ok {
				if cause := c.Cause(); cause != nil {
					err = cause
				} else {
					return err
				}
			} else {
				return err
			}
		}
	}, backoff.WithContext(b, ctx)); err != nil {
		return err
	}
	return nil
}
