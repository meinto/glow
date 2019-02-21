package utils

import "log"

func CheckForError(err error, desc string) {
	if err != nil {
		log.Println(desc)
		log.Println("---")
		log.Panicln(err)
	}
}
