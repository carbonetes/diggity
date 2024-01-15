package table

import (
	"fmt"
	"os"

	"github.com/carbonetes/diggity/pkg/stream"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Margin(1, 0)
	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))
)

type model struct {
	table    table.Model
	duration float64
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, tea.Quit
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if len(m.table.Rows()) == 0 {
		return baseStyle.Render("No component had been found!") + fmt.Sprintf("\nDuration: %.3f sec", m.duration) + "\n" + helpStyle.Render("Press esc to quit... 🐱🐓")
	}
	return baseStyle.Render(m.table.View()) + fmt.Sprintf("\nDuration: %.3f sec", m.duration) + "\n" + helpStyle.Render("Press esc to quit... 🐱🐓")
}

func Create() table.Model {
	columns := []table.Column{
		{Title: "Component", Width: 30},
		{Title: "Type", Width: 12},
		{Title: "Version", Width: 16},
	}

	var rows []table.Row
	components := stream.GetComponents()
	for _, component := range components {
		rows = append(rows, table.Row{
			component.Name,
			component.Type,
			component.Version,
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return t
}

func Show(t table.Model, duration float64) {
	m := model{table: t, duration: duration}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
