package main

import (
	"os"

	"github.com/things-go/ormat/cmd"
)

func main() {
	err := cmd.NewRootCmd().Execute()
	if err != nil {
		os.Exit(1)
	}
}
