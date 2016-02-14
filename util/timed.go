package util

import (
	"fmt"
	"time"
)

func Timed(name string, display bool, fn func()) {
	if display {
		fmt.Printf("Started stage " + Yellow(name) + "\n")
	}

	start := time.Now()

	fn()

	duration := time.Since(start)

	if display {
		fmt.Printf("Stage "+Yellow(name)+" took %.2fms\n", float32(duration)/1000000)
	}
}
