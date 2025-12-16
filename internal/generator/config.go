package generator

import (
	"fmt"
	"os"

	"github.com/smilepakawat/goat/internal/ui"
)

func BuildProjectConfig(m ui.Model, t []string) ProjectConfig {
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

	return ProjectConfig{
		ProjectName: projectName,
		ModuleName:  moduleName,
		Templates:   t,
	}
}
