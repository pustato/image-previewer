package lru

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func newItemStub(size uint64) *Item {
	return &Item{
		Size: size,
	}
}

func TestCache_Complex(t *testing.T) {
	c := NewCache(100, func(item *Item) {
	})

	val1, ok := c.Get("key1")
	require.Nil(t, val1)
	require.False(t, ok)

	item1 := newItemStub(10)
	require.False(t, c.Set("key1", item1))
	require.Equal(t, "key1", item1.key)

	val1, ok = c.Get("key1")
	require.True(t, ok)
	require.Equal(t, item1, val1)

	item2 := newItemStub(20)
	require.False(t, c.Set("key2", item2))
	require.Equal(t, "key2", item2.key)

	val2, ok := c.Get("key2")
	require.True(t, ok)
	require.Equal(t, item2, val2)

	require.True(t, c.Set("key2", item2))
	require.True(t, c.Set("key1", item1))
}

func TestCache_Purge(t *testing.T) {
	removeCounter := 0
	c := NewCache(30, func(item *Item) {
		removeCounter++
	})

	c.Set("key1", newItemStub(10))
	c.Set("key2", newItemStub(10))
	c.Set("key3", newItemStub(10))

	for _, k := range [...]string{"key3", "key2", "key1", "key1", "key2", "key3"} {
		_, hit := c.Get(k)
		require.True(t, hit)
	}

	require.Equal(t, 0, removeCounter)

	c.Set("key4", newItemStub(20))
	require.Equal(t, 2, removeCounter)

	_, hit := c.Get("key1")
	require.False(t, hit)
	_, hit = c.Get("key2")
	require.False(t, hit)

	c.Set("key5", newItemStub(20))
	require.Equal(t, 4, removeCounter)

	_, hit = c.Get("key3")
	require.False(t, hit)
	_, hit = c.Get("key4")
	require.False(t, hit)

	c.Set("key6", newItemStub(100))
	require.Equal(t, 6, removeCounter)

	_, hit = c.Get("key5")
	require.False(t, hit)
	_, hit = c.Get("key6")
	require.False(t, hit)
}

func TestCache_Multithreading(t *testing.T) {
	iterationsCount := 100_000
	c := NewCache(10000, func(item *Item) {
	})
	wg := &sync.WaitGroup{}
	wg.Add(2)
	item := newItemStub(100)

	go func() {
		defer wg.Done()
		for i := 0; i < iterationsCount; i++ {
			key := strconv.Itoa(i)
			c.Set(key, item)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < iterationsCount; i++ {
			key := strconv.Itoa(rand.Intn(iterationsCount * 2)) // nolint:gosec
			c.Get(key)
		}
	}()

	wg.Wait()
}
