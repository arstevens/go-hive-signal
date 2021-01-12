package main

import (
	"log"
	"os"

	"github.com/arstevens/go-hive-signal/internal/configuration"
)

var DefaultConfigLocation = "./sigconf.json"

func main() {
	//Read Configuration
	configFname := DefaultConfigLocation
	if len(os.Args) > 1 {
		configFname = os.Args[1]
	}
	componentConfigs, err := configuration.ReadConfiguration(configFname)
	if err != nil {
		log.Fatalf("Failed to start: %v", err)
	}

	//Run package configurations
	configurators := configuration.GenerateConfiguratorMap()
	for component, confMap := range componentConfigs {
		configurer, ok := configurators[component]
		if !ok {
			log.Fatalf("Failed to start: No configurer associated with %s", component)
		}
		configurer(confMap)
	}

	//Create program
	start, err := configuration.LinkProgram()
	if err != nil {
		log.Fatalf("Failed to start: %v", err)
	}

	//Start program
	start()
}
