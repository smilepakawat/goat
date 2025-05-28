package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	inputProjectName int = iota
	inputModuleName
	done
)

type Model struct {
	State        int
	ProjectInput textinput.Model
	ModuleInput  textinput.Model
}

func NewInitModel() Model {
	pi := textinput.New()
	pi.Focus()
	pi.Width = 30

	mi := textinput.New()
	mi.Focus()
	mi.Width = 50

	return Model{
		State:        inputProjectName,
		ProjectInput: pi,
		ModuleInput:  mi,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.State {
		case inputProjectName:
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "enter":
				m.State = inputModuleName
			default:
				var cmd tea.Cmd
				m.ProjectInput, cmd = m.ProjectInput.Update(msg)
				return m, cmd
			}
		case inputModuleName:
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "enter":
				m.State = done
				return m, tea.Quit
			default:
				var cmd tea.Cmd
				m.ModuleInput, cmd = m.ModuleInput.Update(msg)
				return m, cmd
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	switch m.State {
	case inputProjectName:
		return fmt.Sprintf("Project name:\n%s\n\n(press Enter)", m.ProjectInput.View())
	case inputModuleName:
		return fmt.Sprintf("Module path:\n%s\n\n(press Enter)", m.ModuleInput.View())
	}
	return ""
}
