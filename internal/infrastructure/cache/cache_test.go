package cache_test

import (
	"context"
	"fmt"
	"job_vacancies/internal/infrastructure/cache"
	jobvacancies "job_vacancies/internal/job_vacancies"
	"sync"
	"testing"
	"time"
)

// --- Helpers ---

func newCache[T any](ttl time.Duration) *cache.InMemoryCache[T] {
	return cache.NewInMemoryCache[T](ttl)
}

func makeJobs(n int) []jobvacancies.Job {
	jobs := make([]jobvacancies.Job, n)
	for i := range jobs {
		jobs[i] = jobvacancies.Job{
			Id:    fmt.Sprintf("job-%d", i),
			Title: fmt.Sprintf("Engineer %d", i),
		}
	}
	return jobs
}

type jobCache = *cache.InMemoryCache[[]jobvacancies.Job]

func newJobCache(ttl time.Duration) jobCache {
	return newCache[[]jobvacancies.Job](ttl)
}

func id(s string) string {
	return (s)
}

// --- Tests ---

func TestNewInMemoryCache_CanSetAndGetImmediately(t *testing.T) {
	c := newJobCache(time.Minute)

	c.Set(id("user-1"), makeJobs(1))

	got, ok := c.Get(id("user-1"))
	if !ok {
		t.Fatalf("expected value after Set")
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 job, got %d", len(got))
	}
}

// --- Set ---

func TestSet_StoresItem(t *testing.T) {
	c := newJobCache(time.Minute)
	jobs := makeJobs(3)

	c.Set(id("user-1"), jobs)

	got, ok := c.Get(id("user-1"))
	if !ok {
		t.Fatalf("unexpected miss")
	}
	if len(got) != len(jobs) {
		t.Errorf("expected %d jobs, got %d", len(jobs), len(got))
	}
}

func TestSet_OverwritesExistingItem(t *testing.T) {
	c := newJobCache(time.Minute)

	c.Set(id("user-1"), makeJobs(3))
	c.Set(id("user-1"), makeJobs(1))

	got, ok := c.Get(id("user-1"))
	if !ok {
		t.Fatalf("unexpected miss")
	}
	if len(got) != 1 {
		t.Errorf("expected 1 job after overwrite, got %d", len(got))
	}
}

func TestSet_EmptySlice(t *testing.T) {
	c := newJobCache(time.Minute)

	c.Set(id("empty"), []jobvacancies.Job{})

	got, ok := c.Get(id("empty"))
	if !ok {
		t.Fatalf("unexpected miss")
	}
	if len(got) != 0 {
		t.Errorf("expected empty slice, got %d", len(got))
	}
}

func TestSet_ResetsExpiryOnOverwrite(t *testing.T) {
	c := newJobCache(60 * time.Millisecond)

	c.Set(id("user-1"), makeJobs(1))
	time.Sleep(40 * time.Millisecond)

	c.Set(id("user-1"), makeJobs(2))
	time.Sleep(40 * time.Millisecond)

	got, ok := c.Get(id("user-1"))
	if !ok || len(got) != 2 {
		t.Errorf("expected refreshed item after overwrite")
	}
}

// --- Get ---

func TestGet_ReturnsCorrectJobs(t *testing.T) {
	c := newJobCache(time.Minute)

	jobs := makeJobs(5)
	c.Set(id("user-1"), jobs)

	got, ok := c.Get(id("user-1"))
	if !ok {
		t.Fatalf("unexpected miss")
	}
	if len(got) != len(jobs) {
		t.Errorf("expected %d jobs, got %d", len(jobs), len(got))
	}
}

func TestGet_MissingKey_ReturnsFalse(t *testing.T) {
	c := newJobCache(time.Minute)

	_, ok := c.Get(id("nonexistent"))
	if ok {
		t.Fatal("expected miss")
	}
}

func TestGet_ExpiredItem_ReturnsFalse(t *testing.T) {
	c := newJobCache(50 * time.Millisecond)

	c.Set(id("user-1"), makeJobs(2))

	time.Sleep(100 * time.Millisecond)

	_, ok := c.Get(id("user-1"))
	if ok {
		t.Fatal("expected expired item")
	}
}

func TestGet_ExpiredItem_SubsequentGetReturnsMissing(t *testing.T) {
	c := newJobCache(50 * time.Millisecond)

	c.Set(id("user-1"), makeJobs(1))

	time.Sleep(100 * time.Millisecond)
	c.Get(id("user-1"))

	_, ok := c.Get(id("user-1"))
	if ok {
		t.Fatal("expected deleted item")
	}
}

func TestGet_RefreshesItemLifetime(t *testing.T) {
	c := newJobCache(80 * time.Millisecond)

	c.Set(id("user-1"), makeJobs(1))

	time.Sleep(50 * time.Millisecond)
	_, ok := c.Get(id("user-1"))
	if !ok {
		t.Fatalf("unexpected miss")
	}

	time.Sleep(50 * time.Millisecond)
	_, ok = c.Get(id("user-1"))
	if !ok {
		t.Errorf("expected refreshed lifetime")
	}
}

func TestGet_MultipleKeys_AreIndependent(t *testing.T) {
	c := newJobCache(time.Minute)

	c.Set("a", makeJobs(1))
	c.Set("b", makeJobs(2))
	c.Set("c", makeJobs(3))

	tests := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	for k, expected := range tests {
		got, ok := c.Get(k)
		if !ok {
			t.Errorf("missing key %s", k)
			continue
		}
		if len(got) != expected {
			t.Errorf("key %s: expected %d, got %d", k, expected, len(got))
		}
	}
}

// --- Cleanup ---

func TestStartCleanup_AutomaticallyEvictsExpiredItems(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := newJobCache(30 * time.Millisecond)

	c.Set(id("auto-evict"), makeJobs(1))
	c.StartCleanup(ctx, 50*time.Millisecond)

	time.Sleep(150 * time.Millisecond)

	_, ok := c.Get(id("auto-evict"))
	if ok {
		t.Error("expected eviction")
	}
}

func TestStartCleanup_DoesNotEvictFreshItems(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := newJobCache(500 * time.Millisecond)

	c.Set(id("keep"), makeJobs(1))
	c.StartCleanup(ctx, 50*time.Millisecond)

	time.Sleep(100 * time.Millisecond)

	_, ok := c.Get(id("keep"))
	if !ok {
		t.Errorf("expected item to survive cleanup")
	}
}

// --- Concurrency ---

func TestConcurrentSets_NoDataRace(t *testing.T) {
	c := newJobCache(time.Minute)

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			c.Set(id(fmt.Sprintf("user-%d", i%10)), makeJobs(2))
		}(i)
	}

	wg.Wait()
}

func TestConcurrentGets_NoDataRace(t *testing.T) {
	c := newJobCache(time.Minute)

	for i := 0; i < 10; i++ {
		c.Set(id(fmt.Sprintf("user-%d", i)), makeJobs(2))
	}

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			c.Get(id(fmt.Sprintf("user-%d", i%10)))
		}(i)
	}

	wg.Wait()
}

func TestConcurrentSetAndGet_NoDataRace(t *testing.T) {
	c := newJobCache(time.Minute)

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(2)

		go func(i int) {
			defer wg.Done()
			c.Set(id(fmt.Sprintf("user-%d", i%10)), makeJobs(2))
		}(i)

		go func(i int) {
			defer wg.Done()
			c.Get(id(fmt.Sprintf("user-%d", i%10)))
		}(i)
	}

	wg.Wait()
}

func TestConcurrentSetAndCleanup_NoDataRace(t *testing.T) {
	c := newJobCache(10 * time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c.StartCleanup(ctx, 20*time.Millisecond)

	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			k := id(fmt.Sprintf("key-%d", i))
			c.Set(k, makeJobs(1))

			time.Sleep(15 * time.Millisecond)

			c.Get(k)
		}(i)
	}

	wg.Wait()
}
