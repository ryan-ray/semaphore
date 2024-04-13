package semaphore

import (
	"errors"
)

var ErrNotOpen = errors.New("semaphore is not open")

type semaphore chan struct{}

func _aquire(s semaphore) func() error {
	return func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = ErrNotOpen
			}
		}()

		s <- struct{}{}

		return err
	}
}

func _release(s semaphore) func() error {
	return func() error {
		select {
		case _, ok := <-s:
			if !ok {
				return ErrNotOpen
			}
		default:
		}
		return nil
	}
}

func _cleanup(s semaphore) func() {
	return func() {
		close(s)
		for range s {
		}
	}
}

func _new(size int) (semaphore, func() error, func() error) {
	if size <= 0 {
		size = 1
	}
	sem := make(semaphore, size)
	return sem, _aquire(sem), _release(sem)
}

// New returns aquire and release methods that operate on an underlying channel
// acting as a semaphore, with a capacity indicated by the size parameter.
func New(size int) (func() error, func() error) {
	_, aquire, release := _new(size)
	return aquire, release
}

// WithCancel provides an additional returned close function that closes and
// drains the semaphore.
//
// The underlying channel will be closed and cannot be reused. Any calls to
// the aquire and release functions will return ErrNotOpen. Calling the cancel
// function more that once will panic
func WithCancel(size int) (func() error, func() error, func()) {
	sem, aquire, release := _new(size)
	return aquire, release, _cleanup(sem)
}
