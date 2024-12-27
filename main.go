package main

import (
	"flag"
	"fmt"

	"github.com/apsamuel/brainiac/pkg/common"
)

var ConfigFile string

func init() {
	flag.StringVar(&ConfigFile, "config", "./config/brainiac.yaml", "Path to the configuration file")
	flag.Parse()
}

func main() {
	fmt.Println(common.Foo)
}
