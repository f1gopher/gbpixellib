package ui

import (
	"go-boy/system"
	"log"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func Main(system *system.System) {
	// Log to a file. Useful in debugging since you can't really log to stdout.
	// Not required.
	logfilePath := os.Getenv("BUBBLETEA_LOG")
	if logfilePath != "" {
		if _, err := tea.LogToFile(logfilePath, "simple"); err != nil {
			log.Fatal(err)
		}
	}

	// Initialize our program
	p := tea.NewProgram(model{system: system})
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

// A model can be more or less any type of data. It holds all the data for a
// program, so often it's a struct. For this simple example, however, all
// we'll need is a simple integer.
type model struct {
	system *system.System
}

// Init optionally returns an initial command we should run. In this case we
// want to start the timer.
func (m model) Init() tea.Cmd {
	return tick
}

// Update is called when messages are received. The idea is that you inspect the
// message and send back an updated model accordingly. You can also return
// a command, which is a function that performs I/O and returns a message.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch msgType.Type {
		case tea.KeySpace:
			m.system.Tick()
		case tea.KeyEsc:
			return m, tea.Quit

		}
	case tickMsg:
		//m--
		//if m <= 0 {
		//	return m, tea.Quit
		//}
		return m, tick
	}
	return m, nil
}

// View returns a string based on data in the model. That string which will be
// rendered to the terminal.
func (m model) View() string {
	//return fmt.Sprintf("Hi. This program will exit in %d seconds. To quit sooner press any key.\n", m)
	//return m.system.Pixels()

	cpuStyle := lipgloss.NewStyle().Padding(0, 5)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.system.Pixels(),
		cpuStyle.Render(m.system.State()))
}

// Messages are events that we respond to in our Update function. This
// particular one indicates that the timer has ticked.
type tickMsg time.Time

func tick() tea.Msg {
	time.Sleep(time.Millisecond * 10)
	return tickMsg{}
}