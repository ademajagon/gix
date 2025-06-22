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

func (s *Spinner) Start() {
	frames := []string{
		"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏",
	}

	go func() {
		i := 0
		for atomic.LoadInt32(&s.stopFlag) == 0 {
			fmt.Printf("\r%s", frames[i%len(frames)])
			time.Sleep(100 * time.Millisecond)
			i++
		}
	}()
}

func (s *Spinner) Stop() {
	atomic.StoreInt32(&s.stopFlag, 1)
	fmt.Print("\r")
}
