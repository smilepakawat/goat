package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/smilepakawat/goat/internal/generator"
	"github.com/smilepakawat/goat/internal/ui"
	"github.com/spf13/cobra"
)

var createGinCmd = &cobra.Command{
	Use:   "create-gin",
	Short: "Create a new Go Gin project",
	Long:  `Creates a new Go Gin project with a basic structure and specified options.`,
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(ui.NewInitModel())
		teaModel, err := p.Run()
		if err != nil {
			fmt.Printf("Error, there's been an error: %v", err)
			os.Exit(1)
		}

		model, _ := teaModel.(ui.Model)
		templ := []string {
			"templates/base/gitignore.tmpl",
			"templates/gin/main.go.tmpl",
			"templates/gin/go.mod.tmpl",
		}
		config := generator.BuildProjectConfig(model, templ)

		err = generator.GenerateProject(config)
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
	rootCmd.AddCommand(createGinCmd)
}
