package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/xyproto/symwalk"
)

func main() {
	var mut sync.Mutex
	symwalk.Walk(".", func(p string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return nil
		}
		filename := filepath.Base(p)
		mut.Lock()
		fmt.Printf("%s\n", filename)
		mut.Unlock()
		return nil
	})
}
