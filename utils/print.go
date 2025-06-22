package utils

import (
	"fmt"
	"time"
)

func TypingEffect(s string, d time.Duration) {
	for _, c := range s {
		fmt.Printf("%c", c)
		time.Sleep(d)
	}
	fmt.Println()
}
