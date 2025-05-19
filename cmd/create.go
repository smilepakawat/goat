package cmd

import (
	"fmt"
	"os"

	"github.com/smilepakawat/goat/internal/generator"
	"github.com/spf13/cobra"
)

var projectName string
var moduleName string

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Go Fiber project",
	Long:  `Creates a new Go Fiber project with a basic structure and specified options.`,
	Run: func(cmd *cobra.Command, args []string) {
		if projectName == "" {
			fmt.Println("Error: Project name is required. Use --name flag.")
			os.Exit(1)
		}
		if moduleName == "" {
			fmt.Println("Error: Module name is required. Use --module flag.")
			os.Exit(1)
		}

		fmt.Printf("Creating project '%s' with module '%s'...\n", projectName, moduleName)

		config := generator.ProjectConfig{
			ProjectName:  projectName,
			ModuleName:   moduleName,
		}

		err := generator.GenerateProject(config)
		if err != nil {
			fmt.Printf("Error generating project: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Project '%s' created successfully!\n", projectName)
		fmt.Printf("Next steps:\n")
		fmt.Printf("  cd %s\n", projectName)
		fmt.Printf("  go run cmd/api/main.go\n")
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&projectName, "name", "n", "", "Name of the project (required)")
	createCmd.Flags().StringVarP(&moduleName, "module", "m", "", "Go module name (e.g., github.com/user/project) (required)")

	createCmd.MarkFlagRequired("name")
	createCmd.MarkFlagRequired("module")
}
