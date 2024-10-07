package cmap_test

import (
	"dtalk/internal/data/cmap"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCmap(t *testing.T) {
	cmap := cmap.New[int, int]()
	wg := sync.WaitGroup{}

	count := 100
	for i := 0; i <= count; i++ {
		wg.Add(1)
		go func() {
			cmap.Set(i, i)
			cmap.Get(i)
			wg.Done()
		}()
	}
	wg.Wait()

	total := 0
	for _, v := range cmap.Iter() {
		total += v
		// cmap.Set(0, 0) this will cause a deadlock
	}
	require.Equal(t, count*(count+1)/2, total)

	for i := 0; i <= count; i++ {
		require.Equal(t, i, cmap.Get(i))
	}
}
