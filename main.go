package main

import (
	"github.com/apsamuel/brainiac/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
