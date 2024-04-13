package semaphore_test

import (
	"errors"
	"testing"
	"time"

	"github.com/ryan-ray/semaphore"
)

func TestFlow(t *testing.T) {

	t.Run("Cancel->Aquire", func(t *testing.T) {
		aquire, _, cancel := semaphore.WithCancel(1)
		cancel()

		if err := aquire(); !errors.Is(err, semaphore.ErrNotOpen) {
			t.Fatalf("Got error %v, want %v", err, semaphore.ErrNotOpen)
		}
	})

	t.Run("Aquire->Cancel->Aquire", func(t *testing.T) {
		aquire, _, cancel := semaphore.WithCancel(1)
		if err := aquire(); err != nil {
			t.Fatalf("Got error %v, want %v", err, nil)
		}
		cancel()
		if err := aquire(); !errors.Is(err, semaphore.ErrNotOpen) {
			t.Fatalf("Got error %v, want %v", err, semaphore.ErrNotOpen)
		}
	})

	t.Run("Aquire->Cancel-Release", func(t *testing.T) {
		aquire, release, cancel := semaphore.WithCancel(1)
		aquire()
		cancel()
		if err := release(); !errors.Is(err, semaphore.ErrNotOpen) {
			t.Fatalf("Got error %v, want %v", err, semaphore.ErrNotOpen)
		}
	})

}

func TestSemaphore(t *testing.T) {
	limit := 10
	aquire, release := semaphore.New(limit)
	count := 0
	for range 30 {
		aquire()
		go func() {
			defer release()
			count++
			time.Sleep(time.Millisecond * 200)
			count--
		}()
		if count > limit {
			t.Fatalf("Count: %d should be under %d", count, limit)
		}
	}
}
