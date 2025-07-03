package command

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "Stuber",
	Short: "App for stubs",
	Long:  "App for creating stubs and handling incoming requests",
}

func Execute() error {
	return rootCmd.Execute()
}
