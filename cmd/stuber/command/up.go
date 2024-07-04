package command

import (
	"github.com/spf13/cobra"
	"os"
)

var filePath string

func init() {
	var cmdInitStub = &cobra.Command{
		Use:   "up",
		Short: "Up stub",
		Long:  `Run a stub from the directory`,
		Run: func(cmd *cobra.Command, args []string) {
			if filePath == "" {
				err := cmd.Help()
				if err != nil {
					os.Exit(1)
				}
			}
		},
	}
	cmdInitStub.Flags().StringVarP(&filePath, "file", "f", "", "Path to the file")
	rootCmd.AddCommand(cmdInitStub)
}
