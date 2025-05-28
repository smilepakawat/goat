package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/smilepakawat/goat/internal/generator"
	"github.com/smilepakawat/goat/internal/ui"
	"github.com/spf13/cobra"
)

var projectName string
var moduleName string

var createFiberCmd = &cobra.Command{
	Use:   "create-fiber",
	Short: "Create a new Go Fiber project",
	Long:  `Creates a new Go Fiber project with a basic structure and specified options.`,
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(ui.NewInitModel())
		teaModel, err := p.Run()
		if err != nil {
			fmt.Printf("Error, there's been an error: %v", err)
			os.Exit(1)
		}

		model, _ := teaModel.(ui.Model)
		config := buildProjectConfig(model)

		err = generator.GenerateProject(config)
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
	rootCmd.AddCommand(createFiberCmd)
}

func buildProjectConfig(m ui.Model) generator.ProjectConfig {
	projectName := m.ProjectInput.Value()
	moduleName := m.ModuleInput.Value()
	if projectName == "" {
		fmt.Println("Error: Project name is required.")
		os.Exit(1)
	}
	if moduleName == "" {
		fmt.Println("Error: Module name is required.")
		os.Exit(1)
	}

	return generator.ProjectConfig{
		ProjectName: projectName,
		ModuleName:  moduleName,
	}
}
