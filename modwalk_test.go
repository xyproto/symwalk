package symwalk

import (
	"os"
	"path"
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSymWalk(t *testing.T) {

	// max concurrency out
	runtime.GOMAXPROCS(runtime.NumCPU())

	var seenLock sync.Mutex
	seen := make(map[string]bool)
	walkFunc := func(p string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			filename := path.Base(p)
			seenLock.Lock()
			defer seenLock.Unlock()
			seen[filename] = true
		}
		return nil
	}

	assert.NoError(t, Walk("testdata/root", walkFunc))

	// Check if "found.txt" was found, starting from "testdata/root", then following the symlink to "testdata/other"
	followedSymlink := false
	for fn := range seen {
		if fn == "found.txt" {
			followedSymlink = true
		}
	}

	if !followedSymlink {
		t.FailNow()
	}

	// make sure everything was seen
	if assert.NotEqual(t, len(seen), 0, "Walker should visit at least one file.") {
		for k, v := range seen {
			assert.True(t, v, k)
		}
	}

}
