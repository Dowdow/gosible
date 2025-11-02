package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type loadingModel struct {
	spinner spinner.Model
}

func newLoadingModel() loadingModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#B4BEFE"))
	return loadingModel{spinner: s}
}

func (m loadingModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m loadingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m loadingModel) View() string {
	return fmt.Sprintf("%s Loading configuration...", m.spinner.View())
}
