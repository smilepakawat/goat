package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"text/template"

	"github.com/smilepakawat/goat/pkg"
)

type ProjectConfig struct {
	ProjectName string
	ModuleName  string
	Templates   []string
}

func GenerateProject(config ProjectConfig) error {
	fmt.Printf("Creating project '%s' with module '%s'...\n", config.ProjectName, config.ModuleName)

	if err := os.Mkdir(config.ProjectName, 0755); err != nil {
		return fmt.Errorf("failed to create project directory %s: %w", config.ProjectName, err)
	}
	fmt.Printf("Created directory: %s\n", config.ProjectName)

	templateFiles := mapTemplates(config.Templates, config.ProjectName)

	for tmplPath, outputPath := range templateFiles {
		if err := processTemplate(tmplPath, outputPath, config); err != nil {
			return fmt.Errorf("failed to process template %s: %w", tmplPath, err)
		}
		fmt.Printf("Created file: %s from template %s\n", outputPath, tmplPath)
	}

	return nil
}

func mapTemplates(templates []string, projectName string) map[string]string {
	if len(templates) == 0 {
		return nil
	}

	res := make(map[string]string)
	reg := regexp.MustCompile(`([^/]+?)\.tmpl$`)
	for _, t := range templates {
		matches := reg.FindStringSubmatch(t)
		if len(matches) != 0 {
			res[t] = filepath.Join(projectName, matches[1])
		}
	}
	return res
}

func processTemplate(templatePath, outputPath string, config ProjectConfig) error {
	tmpl, err := loadAndParseTemplate(templatePath)
	if err != nil {
		return fmt.Errorf("failed to load template %s: %w", templatePath, err)
	}

	if err := executeTemplateToFile(tmpl, outputPath, config); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templatePath, err)
	}

	return nil
}

func loadAndParseTemplate(templatePath string) (*template.Template, error) {
	tmplContent, err := pkg.Templates.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}

	tmpl, err := template.New(filepath.Base(templatePath)).Parse(string(tmplContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	return tmpl, nil
}

func executeTemplateToFile(tmpl *template.Template, outputPath string, config ProjectConfig) error {

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	if err := tmpl.Execute(outputFile, config); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}
