package main

import (
	"log"
	"os"

	"github.com/starshine-sys/tribble/types"
	"gopkg.in/yaml.v2"
)

func getConfig() (config *types.BotConfig) {
	config = &types.BotConfig{}

	configFile, err := os.ReadFile("data/config.yaml")
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Loaded configuration file.")

	return config
}
