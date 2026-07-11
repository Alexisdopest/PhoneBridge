package clipboard

import (
	"sync"
	"github.com/atotto/clipboard"
)

// mu ensures atomic access to the clipboard to prevent race conditions
var mu sync.Mutex

// WriteText writes text to the clipboard safely
func WriteText(text string) error {
	mu.Lock()
	defer mu.Unlock()

	if globalWatcher != nil {
		globalWatcher.IgnoreNext()
	}

	return clipboard.WriteAll(text)
}
