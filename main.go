package main

import (
	"flag"
	"fmt"

	"github.com/apsamuel/brainiac/pkg/cache"
	"github.com/apsamuel/brainiac/pkg/database"
)

var ConfigFile string
var Debug bool

func init() {
	flag.StringVar(&ConfigFile, "config", "./config/brainiac.yaml", "Path to the configuration file")
	flag.BoolVar(&Debug, "debug", false, "Enable debug mode")
	flag.Parse()
}

/*
 * This is the main entry point for the brainiac application.
 */
func main() {

	var databaseConfig database.Config
	var cacheConfig cache.Config
	if Debug {
		fmt.Println("debug mode enabled")
	}

	err := databaseConfig.Configure(ConfigFile)
	if err != nil {
		fmt.Println(err)
	}

	_, err = database.MakeStorage(databaseConfig)
	if err != nil {
		fmt.Println(err)
	}

	err = cacheConfig.Configure(ConfigFile)
	if err != nil {
		fmt.Println(err)
	}
	_, err = cache.MakeStorage(cacheConfig)
	if err != nil {
		fmt.Println(err)
	}
}
