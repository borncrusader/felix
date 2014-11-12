package common

import (
	"fmt"
	"os"
)

func CheckError(prefix string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal Error: %s: %s\n", prefix, err.Error())
		os.Exit(1)
	}
}
