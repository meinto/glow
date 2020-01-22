package util

import (
	"errors"
	"os"

	l "github.com/meinto/glow/logging"
)

func ExitOnError(args ...interface{}) {
	if hasError(args...) {
		os.Exit(1)
	}
}

func ExitOnErrorWithMessage(msg string) func(args ...interface{}) {
	return func(args ...interface{}) {
		if hasError(args...) {
			l.Log().Error(errors.New(msg))
			os.Exit(1)
		}
	}
}

func hasError(args ...interface{}) bool {
	last := args[len(args)-1]
	if last != nil {
		return true
	}
	return false
}
