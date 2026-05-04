package cache_test

import (
	"fmt"
	jobvacancies "job_vacancies/internal/job_vacancies"
	"sync"
	"testing"
	"time"

	"job_vacancies/internal/cache"
)

// --- Helpers ---

func newCache(ttl time.Duration) *cache.InMemoryCache {
	return cache.NewInMemoryCache(ttl)
}

func makeJobs(n int) []jobvacancies.Job {
	jobs := make([]jobvacancies.Job, n)
	for i := range jobs {
		jobs[i] = jobvacancies.Job{Id: fmt.Sprintf("job-%d", i), Title: fmt.Sprintf("Engineer %d", i)}
	}
	return jobs
}

// --- NewInMemoryCache ---

func TestNewInMemoryCache_CanSetAndGetImmediately(t *testing.T) {
	c := newCache(time.Minute)
	c.Set("user-1", makeJobs(1))

	_, err := c.Get("user-1")
	if err != nil {
		t.Fatalf("expected no error after Set, got: %v", err)
	}
}

// --- Set ---

func TestSet_StoresItem(t *testing.T) {
	c := newCache(time.Minute)
	jobs := makeJobs(3)

	c.Set("user-1", jobs)

	got, err := c.Get("user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != len(jobs) {
		t.Errorf("expected %d jobs, got %d", len(jobs), len(got))
	}
}

func TestSet_OverwritesExistingItem(t *testing.T) {
	c := newCache(time.Minute)
	c.Set("user-1", makeJobs(3))
	c.Set("user-1", makeJobs(1))

	got, err := c.Get("user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("expected 1 job after overwrite, got %d", len(got))
	}
}

func TestSet_EmptySlice(t *testing.T) {
	c := newCache(time.Minute)
	c.Set("empty", []jobvacancies.Job{})

	got, err := c.Get("empty")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty slice, got %d jobs", len(got))
	}
}

func TestSet_ResetsExpiryOnOverwrite(t *testing.T) {
	// A second Set should reset the TTL clock for that key
	c := newCache(60 * time.Millisecond)

	c.Set("user-1", makeJobs(1))
	time.Sleep(40 * time.Millisecond)

	c.Set("user-1", makeJobs(2)) // reset TTL
	time.Sleep(40 * time.Millisecond)

	// 80ms total elapsed, but TTL was reset at 40ms — item should still be valid
	_, err := c.Get("user-1")
	if err != nil {
		t.Errorf("expected item to be valid after TTL reset via Set, got: %v", err)
	}
}

// --- Get ---

func TestGet_ReturnsCorrectJobs(t *testing.T) {
	c := newCache(time.Minute)
	jobs := makeJobs(5)
	c.Set("user-1", jobs)

	got, err := c.Get("user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != len(jobs) {
		t.Errorf("expected %d jobs, got %d", len(jobs), len(got))
	}
}

func TestGet_MissingKey_ReturnsError(t *testing.T) {
	c := newCache(time.Minute)

	_, err := c.Get("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
	if err.Error() != "cache id doesnt exist" {
		t.Errorf("unexpected error message: %q", err.Error())
	}
}

func TestGet_ExpiredItem_ReturnsError(t *testing.T) {
	c := newCache(50 * time.Millisecond)
	c.Set("user-1", makeJobs(2))

	time.Sleep(100 * time.Millisecond)

	_, err := c.Get("user-1")
	if err == nil {
		t.Fatal("expected expiry error, got nil")
	}
	if err.Error() != "search expired" {
		t.Errorf("unexpected error message: %q", err.Error())
	}
}

func TestGet_ExpiredItem_SubsequentGetReturnsMissingError(t *testing.T) {
	// Verifies lazy deletion: after an expired Get, the item is removed
	// so the next Get returns "doesnt exist" rather than "search expired"
	c := newCache(50 * time.Millisecond)
	c.Set("user-1", makeJobs(1))

	time.Sleep(100 * time.Millisecond)
	c.Get("user-1") // triggers lazy delete

	_, err := c.Get("user-1")
	if err == nil {
		t.Fatal("expected error on second Get, got nil")
	}
	if err.Error() != "cache id doesnt exist" {
		t.Errorf("expected 'doesnt exist' after lazy delete, got: %q", err.Error())
	}
}

func TestGet_RefreshesItemLifetime(t *testing.T) {
	// A Get within TTL should extend the item's lifetime
	c := newCache(80 * time.Millisecond)
	c.Set("user-1", makeJobs(1))

	time.Sleep(50 * time.Millisecond)
	_, err := c.Get("user-1") // access resets lastAccessed
	if err != nil {
		t.Fatalf("unexpected error on first Get: %v", err)
	}

	time.Sleep(50 * time.Millisecond)
	// 100ms total elapsed, but lastAccessed was reset at 50ms — should still be valid
	_, err = c.Get("user-1")
	if err != nil {
		t.Errorf("expected item to still be valid after access refresh, got: %v", err)
	}
}

func TestGet_MultipleKeys_AreIndependent(t *testing.T) {
	c := newCache(time.Minute)
	c.Set("a", makeJobs(1))
	c.Set("b", makeJobs(2))
	c.Set("c", makeJobs(3))

	for id, expectedLen := range map[string]int{"a": 1, "b": 2, "c": 3} {
		got, err := c.Get(id)
		if err != nil {
			t.Errorf("key %q: unexpected error: %v", id, err)
		}
		if len(got) != expectedLen {
			t.Errorf("key %q: expected %d jobs, got %d", id, expectedLen, len(got))
		}
	}
}

// --- StartCleanup ---

func TestStartCleanup_AutomaticallyEvictsExpiredItems(t *testing.T) {
	c := newCache(30 * time.Millisecond)
	c.Set("auto-evict", makeJobs(1))
	c.StartCleanup(50 * time.Millisecond)

	time.Sleep(150 * time.Millisecond)

	_, err := c.Get("auto-evict")
	if err == nil {
		t.Error("expected item to be evicted by background cleanup")
	}
	// After background cleanup, the key is gone entirely (not just expired)
	if err.Error() != "cache id doesnt exist" {
		t.Errorf("expected 'doesnt exist' after background eviction, got: %q", err.Error())
	}
}

func TestStartCleanup_DoesNotEvictFreshItems(t *testing.T) {
	c := newCache(500 * time.Millisecond)
	c.Set("keep", makeJobs(1))
	c.StartCleanup(50 * time.Millisecond)

	time.Sleep(100 * time.Millisecond)

	_, err := c.Get("keep")
	if err != nil {
		t.Errorf("expected fresh item to survive cleanup, got: %v", err)
	}
}

// --- Concurrency ---

func TestConcurrentSets_NoDataRace(t *testing.T) {
	c := newCache(time.Minute)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			id := fmt.Sprintf("user-%d", i%10) // intentional key collisions
			c.Set(id, makeJobs(2))
		}(i)
	}

	wg.Wait()
}

func TestConcurrentGets_NoDataRace(t *testing.T) {
	c := newCache(time.Minute)
	for i := 0; i < 10; i++ {
		c.Set(fmt.Sprintf("user-%d", i), makeJobs(2))
	}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			c.Get(fmt.Sprintf("user-%d", i%10))
		}(i)
	}

	wg.Wait()
}

func TestConcurrentSetAndGet_NoDataRace(t *testing.T) {
	c := newCache(time.Minute)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func(i int) {
			defer wg.Done()
			c.Set(fmt.Sprintf("user-%d", i%10), makeJobs(2))
		}(i)
		go func(i int) {
			defer wg.Done()
			c.Get(fmt.Sprintf("user-%d", i%10)) // errors acceptable; races are not
		}(i)
	}

	wg.Wait()
}

func TestConcurrentSetAndCleanup_NoDataRace(t *testing.T) {
	c := newCache(10 * time.Millisecond)
	c.StartCleanup(20 * time.Millisecond)

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			c.Set(fmt.Sprintf("key-%d", i), makeJobs(1))
			time.Sleep(15 * time.Millisecond)
			c.Get(fmt.Sprintf("key-%d", i))
		}(i)
	}

	wg.Wait()
}
