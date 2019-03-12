package sql

import (
	"context"
	"time"
)

type ReactorType int

const PrimaryReactor ReactorType = 1

func scatter(n int, fn func(i int) error) error {
	errors := make(chan error, n)

	var i int
	for i = 0; i < n; i++ {
		go func(i int) { errors <- fn(i) }(i)
	}

	var err, innerErr error
	for i = 0; i < cap(errors); i++ {
		if innerErr = <-errors; innerErr != nil {
			err = innerErr
		}
	}

	return err
}

// Lock database latch with context
func LockGlobal(ctx context.Context, tx Tx, latch string, timeout int) error {
	var res int
	query := "SELECT GET_LOCK(?, ?)"
	err := tx.QueryRow(ctx, query, latch, timeout).Scan(&res)
	if err != nil {
		return err
	}
	if res != 1 {
		return ErrCaptureLock
	}
	return nil
}

// Unlock database with context
func UnlockGlobal(ctx context.Context, tx Tx, latch string) error {
	var res NullInt64
	err := tx.QueryRow(ctx, "SELECT RELEASE_LOCK(?)", latch).Scan(&res)
	if err != nil {
		return err
	}
	if !res.Valid {
		return ErrReleaseInvalid
	}
	if res.Int64 == 0 {
		return ErrReleaseLock
	}
	return nil
}

// Extract database context from context
func FromContext(ctx context.Context, key ReactorType) Reactor {
	val := ctx.Value(key)
	if c, valid := val.(Reactor); valid {
		return c
	}
	return nil
}

// Append database context into context
func WithContext(
	ctx context.Context,
	reactor Reactor,
) context.Context {
	return context.WithValue(ctx, reactor.Type(), reactor)
}

func Heartbeart(
	ctx context.Context,
	db DB,
	interval time.Duration,
) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(interval):
			db.Ping()
		}
	}
}
