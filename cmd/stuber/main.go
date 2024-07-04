package main

import (
	"Stuber/cmd/stuber/command"
	"os"
)

func main() {
	err := command.Execute()
	if err != nil {
		os.Exit(1)
	}
}
