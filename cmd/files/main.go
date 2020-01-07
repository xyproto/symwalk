package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/xyproto/symwalk"
)

// This is not a pretty function, but the use of `symwalk.Walk` is straightforward enough

func main() {
	// max concurrency out
	runtime.GOMAXPROCS(runtime.NumCPU())

	var seenLock sync.Mutex
	seen := make(map[string]string)
	walkFunc := func(p string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return nil
		}
		dir := filepath.Dir(p)
		seenLock.Lock()
		if info.IsDir() {
			seen[p+"/"] = dir + "/"
		} else {
			seen[p] = dir
		}
		seenLock.Unlock()
		return nil
	}

	flag.Parse()
	dirname := "."
	if len(os.Args) > 1 {
		dirname = os.Args[1]
	}

	symwalk.Walk(dirname, walkFunc)

	// Find all values (directory names) and use them as keys in a new map
	dirnameKeys := make(map[string]bool)
	for _, dir := range seen {
		dirnameKeys[dir] = true
	}
	dirnames := make([]string, len(dirnameKeys))
	i := 0
	for dir := range dirnameKeys {
		dirnames[i] = dir
		i++
	}
	sort.Strings(dirnames)

	printedDir := make(map[string]bool)

	var sb strings.Builder

	// Output a list of all found files, per directory
	for _, dir := range dirnames {
		if dir != "./" {
			if strings.HasSuffix(dir, "/") {
				if _, found := printedDir[dir]; !found {
					if dir != "./" {
						sb.WriteString(dir + "\n")
					}
					printedDir[dir] = true
				}
			} else {
				if _, found := printedDir[dir+"/"]; !found {
					if dir != "." {
						sb.WriteString(dir + "/\n")
					}
					printedDir[dir+"/"] = true
				}
			}
		}
		for p, seenDir := range seen {
			if seenDir == dir {
				depth := len(filepath.SplitList(p))
				if dir == "." {
					depth--
				}
				if !strings.HasSuffix(p, "/") {
					indent := strings.Repeat("\t", depth)
					filename := filepath.Base(p)
					sb.WriteString(indent + filename + "\n")
				}
				delete(seen, p)
			}
		}
	}

	// Output the text
	fmt.Print(sb.String())
}
