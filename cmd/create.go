package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/smilepakawat/goat/internal/generator"
	"github.com/smilepakawat/goat/internal/ui"
	"github.com/spf13/cobra"
)

func createProject(use string, short string, long string, templates []string) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
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

			config := generator.ProjectConfig{
				ProjectName: projectName,
				ModuleName:  moduleName,
				Templates:   templates,
			}
			err = config.GenerateProject()
			if err != nil {
				fmt.Printf("Error generating project: %v\n", err)
				os.Exit(1)
			}

			RunCmd(config.ProjectName, "go", "mod", "tidy")
		},
	}
}
