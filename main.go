package main

import (
	"log"

	"github.com/meinto/glow/cmd"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("glow")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Println("there is no glow config")
	}

	cmd.Execute()
}
