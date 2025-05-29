package generator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/smilepakawat/goat/pkg"
)

type ProjectConfig struct {
	ProjectName string
	ModuleName  string
}

func GenerateProject(config ProjectConfig) error {
	fmt.Printf("Creating project '%s' with module '%s'...\n", config.ProjectName, config.ModuleName)

	if err := os.Mkdir(config.ProjectName, 0755); err != nil {
		return fmt.Errorf("failed to create project directory %s: %w", config.ProjectName, err)
	}
	fmt.Printf("Created directory: %s\n", config.ProjectName)

	templateFiles := map[string]string{
		"templates/fiber/main.go.tmpl": filepath.Join(config.ProjectName, "main.go"),
		"templates/fiber/go.mod.tmpl":  filepath.Join(config.ProjectName, "go.mod"),
	}

	for tmplPath, outputPath := range templateFiles {
		if err := processTemplate(tmplPath, outputPath, config); err != nil {
			return fmt.Errorf("failed to process template %s: %w", tmplPath, err)
		}
		fmt.Printf("Created file: %s from template %s\n", outputPath, tmplPath)
	}

	// TODO: retructure
	fmt.Println("Running 'go mod tidy'...")
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = config.ProjectName
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to run 'go mod tidy': %w\nOutput: %s", err, string(output))
	}
	fmt.Println("'go mod tidy' completed successfully.")

	return nil
}

func processTemplate(templatePath, outputPath string, config ProjectConfig) error {
	tmplContent, err := pkg.Templates.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template file %s: %w", templatePath, err)
	}

	tmpl, err := template.New(filepath.Base(templatePath)).Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", templatePath, err)
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %w", outputPath, err)
	}
	defer outputFile.Close()

	if err := tmpl.Execute(outputFile, config); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templatePath, err)
	}

	return nil
}
