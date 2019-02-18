package cmd

import (
	"log"
)

func CheckRequiredStringField(val, fieldName string) {
	if val == "" {
		log.Fatalf("please provide %s", fieldName)
	}
}
