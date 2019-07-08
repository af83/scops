package main

import (
	"plugin"

	"github.com/af83/scops/api"
	"github.com/af83/scops/config"
	"github.com/af83/scops/logger"
)

func main() {
	config.LoadConfig()

	// load module
	logger.Log.Printf("Load module %v", config.Config.Plugin)
	plug, err := plugin.Open(config.Config.Plugin)
	if err != nil {
		logger.Log.Panicf("Can't open plugin %v: %v", config.Config.Plugin, err)
	}
	// find the feeder
	f, err := plug.Lookup("Feeder")
	if err != nil {
		logger.Log.Panicf("Invalid plugin: %v", err)
	}
	// assert the interface
	feeder, ok := f.(api.Feeder)
	if !ok {
		logger.Log.Panicf("Invalid plugin: invalid type for feeder")
	}

	api.NewProbe(feeder).Run()
}
