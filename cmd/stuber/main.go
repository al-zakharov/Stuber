package main

import (
	"os"
	"stuber/cmd/stuber/command"
)

func main() {
	err := command.Execute()
	if err != nil {
		os.Exit(1)
	}
}
