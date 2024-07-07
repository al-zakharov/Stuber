package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"stuber/internal/router"
	"stuber/internal/yaml"
)

var filePath string

func init() {
	var cmdInitStub = &cobra.Command{
		Use:   "up",
		Short: "Up stub",
		Long:  `Run a stub from the directory`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := runUpStub(cmd); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	cmdInitStub.Flags().StringVarP(&filePath, "file", "f", "", "Path to the file")
	rootCmd.AddCommand(cmdInitStub)
}

func runUpStub(cmd *cobra.Command) error {
	if filePath == "" {
		if err := cmd.Help(); err != nil {
			return fmt.Errorf("failed to display help: %w", err)
		}
		return fmt.Errorf("file path is required")
	}

	yamlStubCollection, err := yaml.NewStubCollection(filePath)
	if err != nil {
		return fmt.Errorf("failed to create stub collection: %w", err)
	}

	c, err := yamlStubCollection.MapToStubs()
	if err != nil {
		return fmt.Errorf("failed to compile stub collection: %w", err)
	}

	router.Run(c)

	return nil
}
