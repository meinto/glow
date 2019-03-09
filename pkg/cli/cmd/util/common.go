package util

import (
	"log"

	"github.com/spf13/viper"
)

func ShouldUseNativeGitBinding(cmd string) bool {
	for _, c := range viper.GetStringSlice("useNativeGitBindings") {
		if c == cmd {
			return true
		}
	}
	return false
}

func CheckRequiredStringField(val, fieldName string) {
	if val == "" {
		log.Fatalf("please provide %s", fieldName)
	}
}
