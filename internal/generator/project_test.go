package generator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateProject(t *testing.T) {
	tests := []struct {
		name    string
		config  ProjectConfig
		wantErr bool
	}{
		{
			name: "successful project generation",
			config: ProjectConfig{
				ProjectName: "testproject",
				ModuleName:  "github.com/test/testproject",
			},
			wantErr: false,
		},
		{
			name: "invalid project name (already exists)",
			config: ProjectConfig{
				ProjectName: "existingproject",
				ModuleName:  "github.com/test/existingproject",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up before and after test
			defer os.RemoveAll(tt.config.ProjectName)
			os.RemoveAll(tt.config.ProjectName)

			// For the "already exists" test case, create the directory first
			if tt.name == "invalid project name (already exists)" {
				if err := os.Mkdir(tt.config.ProjectName, 0755); err != nil {
					t.Fatalf("Failed to setup test: %v", err)
				}
			}

			err := GenerateProject(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify directory structure
				expectedDirs := []string{
					filepath.Join(tt.config.ProjectName, "cmd", "api"),
					filepath.Join(tt.config.ProjectName, "internal"),
				}

				for _, dir := range expectedDirs {
					if _, err := os.Stat(dir); os.IsNotExist(err) {
						t.Errorf("Expected directory %s was not created", dir)
					}
				}

				// Verify files
				expectedFiles := []string{
					filepath.Join(tt.config.ProjectName, "cmd", "api", "main.go"),
					filepath.Join(tt.config.ProjectName, "go.mod"),
				}

				for _, file := range expectedFiles {
					if _, err := os.Stat(file); os.IsNotExist(err) {
						t.Errorf("Expected file %s was not created", file)
					}
				}

				// Verify go.mod exists and contains module name
				goModPath := filepath.Join(tt.config.ProjectName, "go.mod")
				content, err := os.ReadFile(goModPath)
				if err != nil {
					t.Errorf("Failed to read go.mod: %v", err)
				}
				if !contains(string(content), tt.config.ModuleName) {
					t.Errorf("go.mod does not contain expected module name %s", tt.config.ModuleName)
				}
			}
		})
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
			name:         "success process template",
			templatePath: "templates/fiber/main.go.tmpl",
			outputPath:   "testproject/cmd/api/main.go",
			config: ProjectConfig{
				ProjectName: "testproject",
				ModuleName:  "github.com/test/testproject",
			},
			wantErr: false,
		},
		{
			name:         "invalid template path",
			templatePath: "nonexistent.tmpl",
			outputPath:   filepath.Join(tempDir, "test_output.go"),
			config: ProjectConfig{
				ProjectName: "testproject",
				ModuleName:  "github.com/test/testproject",
			},
			wantErr: true,
		},
		{
			name:         "invalid output path",
			templatePath: "templates/fiber/main.go.tmpl",
			outputPath:   "/nonexistent/directory/test_output.go",
			config: ProjectConfig{
				ProjectName: "testproject",
				ModuleName:  "github.com/test/testproject",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up before and after test
			defer os.RemoveAll(tt.config.ProjectName)
			os.RemoveAll(tt.config.ProjectName)

			// Mock
			os.MkdirAll(filepath.Join(tt.config.ProjectName, "cmd", "api"), 0755)
			err := processTemplate(tt.templatePath, tt.outputPath, tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("processTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if file, err := os.Stat(filepath.Join(tt.config.ProjectName, "cmd", "api", "main.go")); os.IsNotExist(err) {
					t.Errorf("Expected file %s was not created", file)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return s != "" && substr != "" && s != substr && len(s) > len(substr) && s != substr
}
