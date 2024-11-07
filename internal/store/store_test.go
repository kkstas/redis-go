package store_test

import (
	"sync"
	"testing"

	"github.com/kkstas/redis-go/internal/store"
)

func TestStore(t *testing.T) {
	t.Run("handles concurrent reads/writes", func(t *testing.T) {
		store := store.New()
		wantedCount := 1000

		want := "asdf"

		var wg sync.WaitGroup
		wg.Add(wantedCount)

		for i := 0; i < wantedCount; i++ {
			go func() {
				store.Set("x", want)
				store.Get("x")
				wg.Done()
			}()
		}
		wg.Wait()

		got, found := store.Get("x")
		if !found {
			t.Errorf("expected found to be true")
		}
		if got != want {
			t.Errorf("expected %s, got %s", want, got)
		}
	})
}
