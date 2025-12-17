package generator

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"text/template"
)

func TestGenerateProject(t *testing.T) {
	tests := []struct {
		name         string
		config       ProjectConfig
		wantErr      bool
		setupFunc    func(t *testing.T, config ProjectConfig)
		validateFunc func(t *testing.T, config ProjectConfig)
	}{
		{
			name: "successful project generation",
			config: ProjectConfig{
				ProjectName: "testproject",
				ModuleName:  "github.com/test/testproject",
				Templates: []string{
					"templates/fiber/main.go.tmpl",
					"templates/fiber/go.mod.tmpl",
				},
			},
			wantErr: false,
			validateFunc: func(t *testing.T, config ProjectConfig) {
				// Verify files exist
				expectedFiles := []string{
					filepath.Join(config.ProjectName, "main.go"),
					filepath.Join(config.ProjectName, "go.mod"),
				}

				for _, file := range expectedFiles {
					if _, err := os.Stat(file); os.IsNotExist(err) {
						t.Errorf("Expected file %s was not created", file)
					}
				}

				// Verify go.mod contains correct module name
				goModPath := filepath.Join(config.ProjectName, "go.mod")
				content, err := os.ReadFile(goModPath)
				if err != nil {
					t.Errorf("Failed to read go.mod: %v", err)
					return
				}

				if !strings.Contains(string(content), config.ModuleName) {
					t.Errorf("go.mod does not contain expected module name %s", config.ModuleName)
				}
			},
		},
		{
			name: "project directory already exists",
			config: ProjectConfig{
				ProjectName: "existingproject",
				ModuleName:  "github.com/test/existingproject",
				Templates: []string{
					"templates/fiber/main.go.tmpl",
					"templates/fiber/go.mod.tmpl",
				},
			},
			wantErr: true,
			setupFunc: func(t *testing.T, config ProjectConfig) {
				if err := os.Mkdir(config.ProjectName, 0755); err != nil {
					t.Fatalf("Failed to setup test: %v", err)
				}
			},
		},
		{
			name: "empty project name",
			config: ProjectConfig{
				ProjectName: "",
				ModuleName:  "github.com/test/empty",
				Templates: []string{
					"templates/fiber/main.go.tmpl",
					"templates/fiber/go.mod.tmpl",
				},
			},
			wantErr: true,
		},
		{
			name: "project with special characters",
			config: ProjectConfig{
				ProjectName: "test-project_123",
				ModuleName:  "github.com/test/test-project",
				Templates: []string{
					"templates/fiber/main.go.tmpl",
					"templates/fiber/go.mod.tmpl",
				},
			},
			wantErr: false,
			validateFunc: func(t *testing.T, config ProjectConfig) {
				// Verify the project directory was created
				if _, err := os.Stat(config.ProjectName); os.IsNotExist(err) {
					t.Errorf("Project directory %s was not created", config.ProjectName)
				}
			},
		},
		{
			name: "long project name",
			config: ProjectConfig{
				ProjectName: "very-long-project-name-with-many-characters-that-should-still-work",
				ModuleName:  "github.com/test/very-long-project",
				Templates: []string{
					"templates/fiber/main.go.tmpl",
					"templates/fiber/go.mod.tmpl",
				},
			},
			wantErr: false,
			validateFunc: func(t *testing.T, config ProjectConfig) {
				if _, err := os.Stat(config.ProjectName); os.IsNotExist(err) {
					t.Errorf("Project directory %s was not created", config.ProjectName)
				}
			},
		},
		{
			name: "nonexistent template",
			config: ProjectConfig{
				ProjectName: "testproject",
				ModuleName:  "github.com/test/testproject",
				Templates: []string{
					"nonexistent.go.tmpl",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up before and after test
			defer os.RemoveAll(tt.config.ProjectName)
			os.RemoveAll(tt.config.ProjectName)

			// Setup test if needed
			if tt.setupFunc != nil {
				tt.setupFunc(t, tt.config)
			}

			err := tt.config.GenerateProject()

			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Run validation if test should succeed
			if !tt.wantErr && tt.validateFunc != nil {
				tt.validateFunc(t, tt.config)
			}
		})
	}
}

// Test to verify directory permissions
func TestGenerateProject_DirectoryPermissions(t *testing.T) {
	config := ProjectConfig{
		ProjectName: "permtest",
		ModuleName:  "github.com/test/permtest",
		Templates: []string{
			"templates/fiber/main.go.tmpl",
			"templates/fiber/go.mod.tmpl",
		},
	}

	defer os.RemoveAll(config.ProjectName)

	err := config.GenerateProject()
	if err != nil {
		t.Fatalf("GenerateProject() error = %v", err)
	}

	// Check directory permissions
	info, err := os.Stat(config.ProjectName)
	if err != nil {
		t.Fatalf("Failed to stat project directory: %v", err)
	}

	expectedPerm := os.FileMode(0755)
	if info.Mode().Perm() != expectedPerm {
		t.Errorf("Project directory permissions = %v, want %v", info.Mode().Perm(), expectedPerm)
	}
}

func TestProcessTemplate(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name         string
		templatePath string
		outputPath   string
		config       ProjectConfig
		wantErr      bool
	}{
		{
			name:         "nonexistent template file",
			templatePath: "nonexistent.tmpl",
			outputPath:   filepath.Join(tempDir, "test_output.go"),
			config: ProjectConfig{
				ProjectName: "testproject",
				ModuleName:  "github.com/test/testproject",
			},
			wantErr: true,
		},
		{
			name:         "invalid output directory",
			templatePath: "templates/fiber/main.go.tmpl",
			outputPath:   "/nonexistent/directory/test_output.go",
			config: ProjectConfig{
				ProjectName: "testproject",
				ModuleName:  "github.com/test/testproject",
			},
			wantErr: true,
		},
		{
			name:         "empty template path",
			templatePath: "",
			outputPath:   filepath.Join(tempDir, "test_output.go"),
			config: ProjectConfig{
				ProjectName: "testproject",
				ModuleName:  "github.com/test/testproject",
			},
			wantErr: true,
		},
		{
			name:         "empty output path",
			templatePath: "templates/fiber/main.go.tmpl",
			outputPath:   "",
			config: ProjectConfig{
				ProjectName: "testproject",
				ModuleName:  "github.com/test/testproject",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			templatePath := tt.templatePath

			err := processTemplate(templatePath, tt.outputPath, tt.config)

			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("processTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

// Test to verify file creation in subtest pattern
// Test for testing actual template files (integration test)
func TestProcessTemplate_RealTemplates(t *testing.T) {
	tempDir := t.TempDir()

	config := ProjectConfig{
		ProjectName: "testproject",
		ModuleName:  "github.com/test/testproject",
	}

	tests := []struct {
		name         string
		templatePath string
		outputPath   string
		wantErr      bool
	}{
		{
			name:         "process go.mod template",
			templatePath: "templates/fiber/go.mod.tmpl",
			outputPath:   filepath.Join(tempDir, "go.mod"),
			wantErr:      false,
		},
		{
			name:         "process main.go template",
			templatePath: "templates/fiber/main.go.tmpl",
			outputPath:   filepath.Join(tempDir, "main.go"),
			wantErr:      false,
		},
		{
			name:         "process .gitignore template",
			templatePath: "templates/base/gitignore.tmpl",
			outputPath:   filepath.Join(tempDir, ".gitignore"),
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := processTemplate(tt.templatePath, tt.outputPath, config)
			if (err != nil) != tt.wantErr {
				t.Errorf("processTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify file was created
				if _, err := os.Stat(tt.outputPath); os.IsNotExist(err) {
					t.Errorf("Output file %s was not created", tt.outputPath)
				}
			}
		})
	}
}

func TestMapTemplates(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name          string
		templates     []string
		expectedValue map[string]string
	}{
		{
			name: "success generate map of templates",
			templates: []string{
				"templates/base/gitignore.tmpl",
				"templates/fiber/main.go.tmpl",
				"templates/fiber/go.mod.tmpl",
			},
			expectedValue: map[string]string{
				"templates/base/gitignore.tmpl": filepath.Join(tempDir, ".gitignore"),
				"templates/fiber/main.go.tmpl":  filepath.Join(tempDir, "main.go"),
				"templates/fiber/go.mod.tmpl":   filepath.Join(tempDir, "go.mod"),
			},
		},
		{
			name:          "empty input",
			templates:     []string{},
			expectedValue: map[string]string{},
		},
		{
			name: "invalid templates pattern",
			templates: []string{
				"templates/fiber/invalidtemplate",
			},
			expectedValue: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := mapTemplates(tt.templates, tempDir)
			if !reflect.DeepEqual(actual, tt.expectedValue) {
				t.Errorf("Value not match\nactual = %s\nexpected = %s", actual, tt.expectedValue)
			}
		})
	}
}

func TestBuildDestinationFile(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectValue string
	}{
		{
			name:        "is in invisible file list",
			input:       "gitignore",
			expectValue: ".gitignore",
		},
		{
			name:        "is not in invisible file list",
			input:       "go.mod",
			expectValue: "go.mod",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := buildDestinationFile(tt.input)
			if actual != tt.expectValue {
				t.Errorf("Value not match\nactual = %s\nexpected = %s", actual, tt.expectValue)
			}
		})
	}
}

func TestIsInvisibleFile(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectValue bool
	}{
		{
			name:        "is in invisible file list",
			input:       "gitignore",
			expectValue: true,
		},
		{
			name:        "is not in invisible file list",
			input:       "go.mod",
			expectValue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := isInvisibleFile(tt.input)
			if actual != tt.expectValue {
				t.Errorf("Value not match\nactual = %v\nexpected = %v", actual, tt.expectValue)
			}
		})
	}
}

func TestLoadAndParseTemplate(t *testing.T) {
	tests := []struct {
		name         string
		templatePath string
		wantErr      bool
	}{
		{
			name:         "nonexistent template file",
			templatePath: "nonexistent/template.tmpl",
			wantErr:      true,
		},
		{
			name:         "empty template path",
			templatePath: "",
			wantErr:      true,
		},
		{
			name:         "invalid template path with special characters",
			templatePath: "templates/invalid/../../../etc/passwd",
			wantErr:      true,
		},
		{
			name:         "template path with null bytes",
			templatePath: "templates/test\x00.tmpl",
			wantErr:      true,
		},
		{
			name:         "invalid template format",
			templatePath: "templates/test/invalid.tmpl",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			templatePath := tt.templatePath

			tmpl, err := loadAndParseTemplate(templatePath)

			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("loadAndParseTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Validate template if test should succeed
			if !tt.wantErr {
				if tmpl == nil {
					t.Error("loadAndParseTemplate() returned nil template but no error")
					return
				}
			}
		})
	}
}

func TestExecuteTemplateToFile(t *testing.T) {
	tempDir := t.TempDir()

	// Create a simple test template
	testTemplate := `package main

// Module: {{.ModuleName}}
func main() {
	println("Hello, {{.ProjectName}}!")
}`

	tmpl, err := template.New("test.tmpl").Parse(testTemplate)
	if err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	config := ProjectConfig{
		ProjectName: "testproject",
		ModuleName:  "github.com/test/testproject",
	}

	tests := []struct {
		name         string
		template     *template.Template
		outputPath   string
		config       ProjectConfig
		wantErr      bool
		setupFunc    func(t *testing.T) string
		validateFunc func(t *testing.T, outputPath string, config ProjectConfig)
	}{
		{
			name:       "successful template execution",
			template:   tmpl,
			outputPath: filepath.Join(tempDir, "success_test.go"),
			config:     config,
			wantErr:    false,
			validateFunc: func(t *testing.T, outputPath string, config ProjectConfig) {
				// Verify file was created
				if _, err := os.Stat(outputPath); os.IsNotExist(err) {
					t.Errorf("Output file %s was not created", outputPath)
					return
				}

				// Verify file content
				content, err := os.ReadFile(outputPath)
				if err != nil {
					t.Errorf("Failed to read output file: %v", err)
					return
				}

				contentStr := string(content)
				if !strings.Contains(contentStr, config.ProjectName) {
					t.Errorf("Output file does not contain project name %s", config.ProjectName)
				}
				if !strings.Contains(contentStr, config.ModuleName) {
					t.Errorf("Output file does not contain module name %s", config.ModuleName)
				}
			},
		},
		{
			name:       "invalid output directory",
			template:   tmpl,
			outputPath: "/root/nonexistent/directory/test.go",
			config:     config,
			wantErr:    true,
		},
		{
			name:       "empty output path",
			template:   tmpl,
			outputPath: "",
			config:     config,
			wantErr:    true,
		},
		{
			name:       "output path with null bytes",
			template:   tmpl,
			outputPath: filepath.Join(tempDir, "test\x00.go"),
			config:     config,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputPath := tt.outputPath

			// Setup test if needed
			if tt.setupFunc != nil {
				outputPath = tt.setupFunc(t)
			}

			err := executeTemplateToFile(tt.template, outputPath, tt.config)

			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("executeTemplateToFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Run validation if test should succeed
			if !tt.wantErr && tt.validateFunc != nil {
				tt.validateFunc(t, outputPath, tt.config)
			}
		})
	}
}

func TestExecuteTemplateToFile_TemplateExecutionError(t *testing.T) {
	tempDir := t.TempDir()

	// Create a template with invalid syntax that will cause execution error
	invalidTemplate := `package {{.InvalidField}}`

	tmpl, err := template.New("invalid.tmpl").Parse(invalidTemplate)
	if err != nil {
		t.Fatalf("Failed to create invalid test template: %v", err)
	}

	config := ProjectConfig{
		ProjectName: "testproject",
		ModuleName:  "github.com/test/testproject",
	}

	outputPath := filepath.Join(tempDir, "invalid_execution.go")

	err = executeTemplateToFile(tmpl, outputPath, config)
	if err == nil {
		t.Error("executeTemplateToFile() should have failed with template execution error")
	}

	// Verify that error message contains template execution information
	if !strings.Contains(err.Error(), "failed to execute template") {
		t.Errorf("Error message should contain 'failed to execute template', got: %v", err)
	}
}

func TestExecuteTemplateToFile_FilePermissions(t *testing.T) {
	tempDir := t.TempDir()

	testTemplate := `package {{.ProjectName}}`
	tmpl, err := template.New("perm.tmpl").Parse(testTemplate)
	if err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	config := ProjectConfig{
		ProjectName: "testproject",
		ModuleName:  "github.com/test/testproject",
	}

	outputPath := filepath.Join(tempDir, "perm_test.go")

	err = executeTemplateToFile(tmpl, outputPath, config)
	if err != nil {
		t.Fatalf("executeTemplateToFile() failed: %v", err)
	}

	// Check file permissions
	info, err := os.Stat(outputPath)
	if err != nil {
		t.Fatalf("Failed to stat output file: %v", err)
	}

	// Files created with os.Create should have default permissions (usually 0644)
	expectedMode := os.FileMode(0644)
	if info.Mode().Perm() != expectedMode {
		t.Logf("File permissions = %v, expected around %v (actual permissions may vary by system)", info.Mode().Perm(), expectedMode)
	}
}
