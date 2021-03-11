package main

import (
	"os"

	"github.com/starshine-sys/tribble/types"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

func getConfig(sugar *zap.SugaredLogger) (config *types.BotConfig) {
	config = &types.BotConfig{}

	configFile, err := os.ReadFile("data/config.yaml")
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		sugar.Fatal(err)
	}
	sugar.Infof("Loaded configuration file.")

	return config
}
