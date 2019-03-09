package util

import (
	"log"
	"os"
)

func CheckForError(err error, desc string) {
	if err != nil {
		log.Printf("%+v", err)
		os.Exit(1)
	}
}
