package table

import (
	"fmt"
	"os"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Margin(1, 0)
	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.HiddenBorder()).
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

	if len(m.table.Rows()) == 0 {
		return m, tea.Quit
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if len(m.table.Rows()) == 0 {
		return baseStyle.Render("No component had been found!") + fmt.Sprintf("\nDuration: %.3f sec", m.duration) + "\n" + helpStyle.Render("Press esc to quit... 🐱🐓")
	}
	return baseStyle.Render(m.table.View()) + helpStyle.Render("Use the ↑ and ↓ arrow keys to select different components in the table.\nPress esc to quit... 🐱🐓") + fmt.Sprintf("\nDuration: %.3f sec\n", m.duration)
}

func Create(bom *cyclonedx.BOM) table.Model {
	columns := []table.Column{
		{Title: "Component", Width: 32},
		{Title: "Type", Width: 16},
		{Title: "Version", Width: 20},
	}

	var rows []table.Row
	components := bom.Components
	for _, component := range *components {

		var componentType string
		for _, property := range *component.Properties {
			if property.Name == "diggity:package:type" {
				componentType = property.Value
			}
		}

		if component.Type == cyclonedx.ComponentTypeOS {
			componentType = "operating-system"
		}

		rows = append(rows, table.Row{
			component.Name,
			componentType,
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
