//go:build unix

package playwrightdebug

import (
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestWithEnsureLockSerializesConcurrentAccess(t *testing.T) {
	cacheDir := t.TempDir()
	var active int32
	var maxActive int32
	var wg sync.WaitGroup

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := withEnsureLock(cacheDir, func() error {
				current := atomic.AddInt32(&active, 1)
				for {
					prev := atomic.LoadInt32(&maxActive)
					if current <= prev || atomic.CompareAndSwapInt32(&maxActive, prev, current) {
						break
					}
				}
				time.Sleep(20 * time.Millisecond)
				atomic.AddInt32(&active, -1)
				return nil
			})
			if err != nil {
				t.Errorf("withEnsureLock: %v", err)
			}
		}()
	}

	wg.Wait()
	if got := atomic.LoadInt32(&maxActive); got != 1 {
		t.Fatalf("max concurrent holders = %d, want 1", got)
	}
}

func TestPlaywrightInstalled(t *testing.T) {
	cacheDir := t.TempDir()
	if playwrightInstalled(cacheDir) {
		t.Fatal("expected missing playwright package")
	}

	playwrightDir := filepath.Join(cacheDir, "node_modules", "playwright")
	if err := os.MkdirAll(playwrightDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(playwrightDir, "package.json"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if !playwrightInstalled(cacheDir) {
		t.Fatal("expected installed playwright package")
	}
}