package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewInitModel(t *testing.T) {
	model := NewInitModel()

	// Test initial state
	if model.State != inputProjectName {
		t.Errorf("Expected initial state to be inputProjectName (%d), got %d", inputProjectName, model.State)
	}

	// Test ProjectInput initialization
	if model.ProjectInput.Width != 30 {
		t.Errorf("Expected ProjectInput width to be 30, got %d", model.ProjectInput.Width)
	}

	if !model.ProjectInput.Focused() {
		t.Error("Expected ProjectInput to be focused initially")
	}

	// Test ModuleInput initialization
	if model.ModuleInput.Width != 50 {
		t.Errorf("Expected ModuleInput width to be 50, got %d", model.ModuleInput.Width)
	}

	if !model.ModuleInput.Focused() {
		t.Error("Expected ProjectInput to be focused initially")
	}
}

func TestInit(t *testing.T) {
	model := NewInitModel()
	cmd := model.Init()

	if cmd == nil {
		t.Error("Expected Init() to return a command, got nil")
	}
}

func TestUpdate_ProjectNameInput(t *testing.T) {
	tests := []struct {
		name          string
		keyMsg        string
		expectedState int
		shouldQuit    bool
		shouldGetCmd  bool
	}{
		{
			name:          "ctrl+c quits",
			keyMsg:        "ctrl+c",
			expectedState: inputProjectName,
			shouldQuit:    true,
		},
		{
			name:          "enter advances to module input",
			keyMsg:        "enter",
			expectedState: inputModuleName,
			shouldQuit:    false,
		},
		{
			name:          "regular key input",
			keyMsg:        "a",
			expectedState: inputProjectName,
			shouldQuit:    false,
			shouldGetCmd:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewInitModel()
			model.State = inputProjectName

			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.keyMsg)}
			if tt.keyMsg == "ctrl+c" {
				keyMsg = tea.KeyMsg{Type: tea.KeyCtrlC}
			} else if tt.keyMsg == "enter" {
				keyMsg = tea.KeyMsg{Type: tea.KeyEnter}
			}

			newModel, cmd := model.Update(keyMsg)
			updatedModel := newModel.(Model)

			if updatedModel.State != tt.expectedState {
				t.Errorf("Expected state %d, got %d", tt.expectedState, updatedModel.State)
			}

			if tt.shouldQuit && cmd == nil {
				t.Error("Expected quit command")
			}

			if !tt.shouldQuit && !tt.shouldGetCmd && cmd != nil {
				t.Error("Expected no command")
			}

			if tt.shouldGetCmd && cmd == nil {
				t.Error("Expected a command but got nil")
			}
		})
	}
}

func TestUpdate_ModuleNameInput(t *testing.T) {
	tests := []struct {
		name          string
		keyMsg        string
		expectedState int
		shouldQuit    bool
		shouldGetCmd  bool
	}{
		{
			name:          "ctrl+c quits",
			keyMsg:        "ctrl+c",
			expectedState: inputModuleName,
			shouldQuit:    true,
		},
		{
			name:          "enter completes and quits",
			keyMsg:        "enter",
			expectedState: done,
			shouldQuit:    true,
		},
		{
			name:          "regular key input",
			keyMsg:        "g",
			expectedState: inputModuleName,
			shouldQuit:    false,
			shouldGetCmd:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewInitModel()
			model.State = inputModuleName

			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.keyMsg)}
			if tt.keyMsg == "ctrl+c" {
				keyMsg = tea.KeyMsg{Type: tea.KeyCtrlC}
			} else if tt.keyMsg == "enter" {
				keyMsg = tea.KeyMsg{Type: tea.KeyEnter}
			}

			newModel, cmd := model.Update(keyMsg)
			updatedModel := newModel.(Model)

			if updatedModel.State != tt.expectedState {
				t.Errorf("Expected state %d, got %d", tt.expectedState, updatedModel.State)
			}

			if tt.shouldQuit && cmd == nil {
				t.Error("Expected quit command")
			}

			if !tt.shouldQuit && !tt.shouldGetCmd && cmd != nil {
				t.Error("Expected no command")
			}

			if tt.shouldGetCmd && cmd == nil {
				t.Error("Expected a command but got nil")
			}
		})
	}
}

func TestUpdate_NonKeyMessage(t *testing.T) {
	model := NewInitModel()

	// Test with a non-key message
	newModel, cmd := model.Update("not a key message")
	updatedModel := newModel.(Model)

	if updatedModel.State != model.State {
		t.Error("Expected model state to remain unchanged for non-key messages")
	}

	if cmd != nil {
		t.Error("Expected no command for non-key messages")
	}
}

func TestView_ProjectNameState(t *testing.T) {
	model := NewInitModel()
	model.State = inputProjectName

	view := model.View()

	expectedPrefix := "Project name:"
	if len(view) < len(expectedPrefix) || view[:len(expectedPrefix)] != expectedPrefix {
		t.Errorf("Expected view to start with '%s', got: %s", expectedPrefix, view)
	}

	if view == "" {
		t.Error("Expected non-empty view for project name state")
	}
}

func TestView_ModuleNameState(t *testing.T) {
	model := NewInitModel()
	model.State = inputModuleName

	view := model.View()

	expectedPrefix := "Module path:"
	if len(view) < len(expectedPrefix) || view[:len(expectedPrefix)] != expectedPrefix {
		t.Errorf("Expected view to start with '%s', got: %s", expectedPrefix, view)
	}

	if view == "" {
		t.Error("Expected non-empty view for module name state")
	}
}

func TestView_DoneState(t *testing.T) {
	model := NewInitModel()
	model.State = done

	view := model.View()

	if view != "" {
		t.Errorf("Expected empty view for done state, got: %s", view)
	}
}

func TestView_InvalidState(t *testing.T) {
	model := NewInitModel()
	model.State = 999 // Invalid state

	view := model.View()

	if view != "" {
		t.Errorf("Expected empty view for invalid state, got: %s", view)
	}
}

func TestConstants(t *testing.T) {
	// Test that constants have expected values
	if inputProjectName != 0 {
		t.Errorf("Expected inputProjectName to be 0, got %d", inputProjectName)
	}

	if inputModuleName != 1 {
		t.Errorf("Expected inputModuleName to be 1, got %d", inputModuleName)
	}

	if done != 2 {
		t.Errorf("Expected done to be 2, got %d", done)
	}
}

func TestStateTransitions(t *testing.T) {
	model := NewInitModel()

	// Test full workflow: project name -> module name -> done

	// Start at project name input
	if model.State != inputProjectName {
		t.Fatalf("Expected to start at inputProjectName state")
	}

	// Press enter to advance to module name
	enterKey := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ := model.Update(enterKey)
	model = newModel.(Model)

	if model.State != inputModuleName {
		t.Errorf("Expected to advance to inputModuleName state, got %d", model.State)
	}

	// Press enter again to complete
	newModel, cmd := model.Update(enterKey)
	model = newModel.(Model)

	if model.State != done {
		t.Errorf("Expected to advance to done state, got %d", model.State)
	}

	if cmd == nil {
		t.Error("Expected quit command when transitioning to done state")
	}
}

func TestTextInputIntegration(t *testing.T) {
	model := NewInitModel()

	// Test that text input works in project name state
	model.State = inputProjectName

	// Simulate typing "test"
	for _, char := range "test" {
		keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{char}}
		newModel, _ := model.Update(keyMsg)
		model = newModel.(Model)
	}

	projectValue := model.ProjectInput.Value()
	if projectValue != "test" {
		t.Errorf("Expected project input value to be 'test', got '%s'", projectValue)
	}

	// Switch to module name input
	enterKey := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ := model.Update(enterKey)
	model = newModel.(Model)

	// Test that text input works in module name state
	for _, char := range "github.com/test/project" {
		keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{char}}
		newModel, _ := model.Update(keyMsg)
		model = newModel.(Model)
	}

	moduleValue := model.ModuleInput.Value()
	if moduleValue != "github.com/test/project" {
		t.Errorf("Expected module input value to be 'github.com/test/project', got '%s'", moduleValue)
	}
}

func TestModelFields(t *testing.T) {
	model := NewInitModel()

	// Test that all fields are properly initialized
	if model.ProjectInput.Width == 0 {
		t.Error("ProjectInput should be initialized with non-zero width")
	}

	if model.ModuleInput.Width == 0 {
		t.Error("ModuleInput should be initialized with non-zero width")
	}

	// Test that ProjectInput is focused initially
	if !model.ProjectInput.Focused() {
		t.Error("ProjectInput should be focused initially")
	}
}

// Benchmark tests
func BenchmarkNewInitModel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewInitModel()
	}
}

func BenchmarkUpdate(b *testing.B) {
	model := NewInitModel()
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		model.Update(keyMsg)
	}
}

func BenchmarkView(b *testing.B) {
	model := NewInitModel()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		model.View()
	}
}
