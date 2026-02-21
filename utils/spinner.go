package utils

import (
	"fmt"
	"sync/atomic"
	"time"
)

type Spinner struct {
	stopFlag int32
}

func NewSpinner() *Spinner {
	return &Spinner{}
}

// Start begins animating the spinner in a background goroutine
func (s *Spinner) Start() {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	go func() {
		for i := 0; atomic.LoadInt32(&s.stopFlag) == 0; i++ {
			fmt.Printf("\r%s", frames[i%len(frames)])
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

// Stop terminates the spinner and clears the line
func (s *Spinner) Stop() {
	atomic.StoreInt32(&s.stopFlag, 1)
	fmt.Print("\r\033[K")
}
