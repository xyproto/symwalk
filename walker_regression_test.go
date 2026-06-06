package symwalk

import (
	"errors"
	"os"
	"sync/atomic"
	"testing"
	"time"
)

// TestWalkLimitConcurrentErrorsNoPanic stress-reproduces the data race and
// double close(kill) reported in algernon issue #121: walkFn returns an error
// from many concurrent workers. Before the rewrite this either panicked with
// "close of closed channel" or deadlocked.
func TestWalkLimitConcurrentErrorsNoPanic(t *testing.T) {
	makeTestFiles(10, 20)
	defer deleteTestFiles()

	for i := 0; i < 200; i++ {
		err := WalkLimit(testFiles, func(_ string, _ os.FileInfo, _ error) error {
			return errors.New("forced error")
		}, 32)
		if err == nil {
			t.Fatalf("iter %d: expected an error", i)
		}
	}
}

// TestWalkLimitTerminatesQuicklyOnError verifies the walker stops promptly
// once an error is recorded rather than processing every remaining file.
func TestWalkLimitTerminatesQuicklyOnError(t *testing.T) {
	makeTestFiles(20, 50)
	defer deleteTestFiles()

	var visited atomic.Int64
	done := make(chan struct{})
	go func() {
		_ = WalkLimit(testFiles, func(_ string, info os.FileInfo, _ error) error {
			if info != nil && !info.IsDir() {
				if visited.Add(1) == 1 {
					return errors.New("stop")
				}
			}
			return nil
		}, 8)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("WalkLimit did not terminate within 5s after first error")
	}
}
