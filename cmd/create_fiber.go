package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/smilepakawat/goat/internal/generator"
	"github.com/smilepakawat/goat/internal/ui"
	"github.com/spf13/cobra"
)

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
		projectName := model.ProjectInput.Value()
		moduleName := model.ModuleInput.Value()
		if projectName == "" {
			fmt.Println("Error: Project name is required.")
			os.Exit(1)
		}
		if moduleName == "" {
			fmt.Println("Error: Module name is required.")
			os.Exit(1)
		}

		templates := []string{
			"templates/base/gitignore.tmpl",
			"templates/fiber/main.go.tmpl",
			"templates/fiber/go.mod.tmpl",
		}
		config := generator.ProjectConfig{
			ProjectName: projectName,
			ModuleName: moduleName,
			Templates: templates,
		}
		err = config.GenerateProject()
		if err != nil {
			fmt.Printf("Error generating project: %v\n", err)
			os.Exit(1)
		}

		RunCmd(config.ProjectName, "go", "mod", "tidy")

		fmt.Printf("Project '%s' created successfully!\n", config.ProjectName)
		fmt.Printf("Next steps:\n")
		fmt.Printf("  cd %s\n", config.ProjectName)
		fmt.Printf("  go run main.go\n")
	},
}

func init() {
	rootCmd.AddCommand(createFiberCmd)
}
