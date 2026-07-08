package clipboard

import (
	"sync"
	"github.com/atotto/clipboard"
)

// mu ensures atomic access to the clipboard to prevent race conditions
var mu sync.Mutex

// WriteText safely writes text to the system clipboard
func WriteText(text string) error {
	mu.Lock()
	defer mu.Unlock()
	return clipboard.WriteAll(text)
}
