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
			if filePath == "" {
				err := cmd.Help()
				if err != nil {
					os.Exit(1)
				}
			}

			//TODO handle error
			yamlStubCollection, err := yaml.NewStubCollection(filePath)
			if err != nil {
				fmt.Println(err)
			}
			stubCollection := make([]*router.Stub, 0)

			for _, yamlStub := range yamlStubCollection.Items {
				stubCollection = append(stubCollection, &router.Stub{
					HttpMethod: yamlStub.HttpMethod,
					Path:       yamlStub.Path,
					Body:       yamlStub.Body,
					Status:     yamlStub.Status,
				})
			}
			router.Register(stubCollection)
		},
	}
	cmdInitStub.Flags().StringVarP(&filePath, "file", "f", "", "Path to the file")
	rootCmd.AddCommand(cmdInitStub)
}
