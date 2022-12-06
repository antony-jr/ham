package get

import (
	"fmt"
	"strings"
	"time"

	"encoding/json"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kyokomi/emoji/v2"
)

type model struct {
	shell      *SSHShellContext
	percentage int
	prog       string
	width      int
	height     int
	spinner    spinner.Model
	progress   progress.Model
	done       bool
}

var (
	currentStatusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	doneStyle          = lipgloss.NewStyle().Margin(1, 2)
	checkMark          = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")
	crossMark          = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).SetString(emoji.Sprintf(":prohibited:"))
)

func newModel(shell *SSHShellContext) model {
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
	)
	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	s.Spinner = spinner.Points
	return model{
		shell:      shell,
		percentage: 0,
		prog:       "Building",
		spinner:    s,
		progress:   p,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tea.Println("  Tracking Remote Build..."), m.spinner.Tick, refreshProgress(m.shell))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit
		}

	case errorCode:
		return m, tea.Quit

	case statusJson:
		raw_json := []byte(statusJson(msg))
		var result map[string]interface{}

		err := json.Unmarshal(raw_json, &result)
		if err != nil {
			return m, tea.Batch(
				tea.Printf(" %sCannot Get Progress from Remote. (%s)\n", crossMark, err.Error()),
				withErrorQuit(m.shell, SSH_SHELL_MALFORMED_JSON),
			)
		}

		isErr := result["error"].(interface{}).(bool)

		if isErr {
			erMsg := result["message"].(interface{}).(string)
			m.prog = "Build Failed"
			return m, tea.Batch(
				tea.Printf(" %s%s", crossMark, erMsg),
				tea.Printf(" %sBuild Failed.\n", crossMark),
				withErrorQuit(m.shell, SSH_SHELL_HAM_STATUS_ERRORED),
			)
		}

		m.prog = result["progress"].(interface{}).(string)
		m.percentage = int(result["percentage"].(interface{}).(float64))

		if m.percentage == 100 {
			// Everything's been installed. We're done!
			m.done = true

			return m, tea.Batch(
				tea.Printf("  %s Remote Build Completed\n", checkMark),
				withErrorQuit(m.shell, SSH_SHELL_NO_ERROR),
			)
		}

		// Update progress bar
		progressCmd := m.progress.SetPercent(float64(m.percentage))

		return m, tea.Batch(
			progressCmd,
			// tea.Printf(" %s %s", checkMark, m.prog),
			refreshProgress(m.shell),
		)
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case progress.FrameMsg:
		newModel, cmd := m.progress.Update(msg)
		if newModel, ok := newModel.(progress.Model); ok {
			m.progress = newModel
		}
		return m, cmd
	}
	return m, nil
}

func (m model) View() string {
	n := 100
	w := lipgloss.Width(fmt.Sprintf("%d", n))

	if m.done {
		return ""
	}

	pkgCount := fmt.Sprintf(" %*d/%*d", w, m.percentage, w, n)

	spin := "  " + m.spinner.View() + " "
	prog := m.progress.View()
	cellsAvail := max(0, m.width-lipgloss.Width(spin+prog+pkgCount))

	progName := currentStatusStyle.Render(m.prog)
	info := lipgloss.NewStyle().MaxWidth(cellsAvail).Render(progName)

	cellsRemaining := max(0, m.width-lipgloss.Width(spin+info+prog+pkgCount))
	gap := strings.Repeat(" ", cellsRemaining)

	return spin + info + gap + prog + pkgCount
}

type statusJson string
type errorCode SSHShellCode

func withErrorQuit(shell *SSHShellContext, code SSHShellCode) tea.Cmd {
	d := time.Second * time.Duration(5)
	return tea.Tick(d, func(t time.Time) tea.Msg {
		shell.SetCode(code)
		return errorCode(code)
	})
}

func refreshProgress(shell *SSHShellContext) tea.Cmd {
	d := time.Second * time.Duration(2)
	return tea.Tick(d, func(t time.Time) tea.Msg {
		out, err := shell.Exec("ham build-status | cat |  grep -a Status | cut -c 10-")
		if err != nil {
			shell.SetCode(SSH_SHELL_CANNOT_CONNECT)
			return errorCode(SSH_SHELL_CANNOT_CONNECT)
		}

		if len(out) == 0 {
			return statusJson("{\"error\": true, \"message\": \"Remote Server not Responding Build Status\"}")
		}
		return statusJson(out)
	})
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func runProgressTeaProgram(shell *SSHShellContext) error {
	if _, err := tea.NewProgram(newModel(shell)).Run(); err != nil {
		return err
	}

	return nil
}
