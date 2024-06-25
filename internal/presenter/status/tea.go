package status

import (
	"os"
	"strings"

	"github.com/carbonetes/diggity/internal/log"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Margin(1, 0)
	dotStyle     = helpStyle.UnsetMargins()
	appStyle     = lipgloss.NewStyle().Margin(1, 2, 0, 2)
	exited       = true
)

type resultMsg struct {
	file string
	done bool
}

type errorMsg struct {
	err  error
	quit bool
}

func (r resultMsg) String() string {
	if r.file == "" {
		return dotStyle.Render(strings.Repeat(".", 30))
	}
	return r.file
}

type model struct {
	spinner  spinner.Model
	results  []resultMsg
	errorM   errorMsg
	quitting bool
}

func New() model {
	s := spinner.New()
	s.Style = spinnerStyle
	return model{
		spinner: s,
		results: make([]resultMsg, 5),
	}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc {
			os.Exit(0)
			return m, nil
		} else {
			return m, nil
		}
	case resultMsg:
		if !msg.done {
			m.results = append(m.results[1:], msg)
			return m, nil
		}
		m.quitting = true
		return m, tea.Quit

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case errorMsg:
		m.errorM = msg
		return m, nil
	default:
		return m, nil
	}
}

func (m model) View() string {
	var s string

	if m.errorM.err != nil {
		s += appStyle.Render(m.errorM.err.Error())
		if m.errorM.quit {
			s += "\n"
		}
		return s
	}

	if m.quitting {
		return ""
	} else {
		s += m.spinner.View() + " Diggity is searching the files..."
	}

	s += "\n\n"

	for _, res := range m.results {
		s += res.String() + "\n"
	}

	if !m.quitting {
		s += helpStyle.Render("Hang tight, we're scanning at the speed of light! 🚀")
	}

	if m.quitting {
		s += "\n"
	}

	return appStyle.Render(s)
}

var (
	m = New()
	p = tea.NewProgram(m)
)

func Run() {
	exited = false
	go func() {
		if _, err := p.Run(); err != nil {
			log.Debugf("Failed to start program: %v", err)
		}
	}()
}

func Done() {
	if exited {
		return
	}
	p.Send(resultMsg{done: true})
}

func AddFile(file string) {
	if exited {
		return
	}
	p.Send(resultMsg{file: file, done: false})
}
