package clipboard

import (
	"time"

	"github.com/atotto/clipboard"
)

type Watcher interface {
	Start()
	Stop()
	Events() <-chan string
	IgnoreNext()
}

type PollingWatcher struct {
	events     chan string
	stopChan   chan struct{}
	ignoreNext bool
}

func NewPollingWatcher() *PollingWatcher {
	return &PollingWatcher{
		events:   make(chan string, 10),
		stopChan: make(chan struct{}),
	}
}

func (w *PollingWatcher) Start() {
	go func() {
		lastText, _ := clipboard.ReadAll()
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-w.stopChan:
				return
			case <-ticker.C:
				text, err := clipboard.ReadAll()
				if err != nil {
					continue
				}

				if text != lastText {
					lastText = text
					if w.ignoreNext {
						w.ignoreNext = false
						continue
					}
					if text != "" {
						w.events <- text
					}
				}
			}
		}
	}()
}

func (w *PollingWatcher) Stop() {
	close(w.stopChan)
}

func (w *PollingWatcher) Events() <-chan string {
	return w.events
}

func (w *PollingWatcher) IgnoreNext() {
	w.ignoreNext = true
}

var globalWatcher Watcher

func InitWatcher() Watcher {
	globalWatcher = NewPollingWatcher()
	return globalWatcher
}

func GetWatcher() Watcher {
	return globalWatcher
}
