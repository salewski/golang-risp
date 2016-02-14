package util

import (
	"fmt"
	"os"
)

func ReportError(err error, recoverable bool) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())

		if !recoverable {
			os.Exit(1)
		}
	}
}
