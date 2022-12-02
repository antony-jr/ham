package get

import (
   "fmt"
   "strings"
   "time"
   "math/rand"
   "os"

   "github.com/charmbracelet/bubbles/progress"
   "github.com/charmbracelet/bubbles/spinner"
   tea "github.com/charmbracelet/bubbletea"
   "github.com/charmbracelet/lipgloss"
)

type model struct {
        percentage int
	prog     string
	width    int
	height   int
	spinner  spinner.Model
	progress progress.Model
	done     bool
}

var (
	currentStatusStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	doneStyle           = lipgloss.NewStyle().Margin(1, 2)
	checkMark           = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")
)

func newModel() model {
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
	)
	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	return model{
	   	percentage: 0,
	        prog: "Building",
		spinner:  s,
		progress: p,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick)
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

	case stillMsg:
		// Update progress bar
		progressCmd := m.progress.SetPercent(float64(m.percentage))

		m.prog = "Started"

		return m, tea.Batch(
		 	progressCmd,
			downloadAndInstall(),             // download the next package
		)
	case statusMsg:
		if m.percentage == 100 {
			// Everything's been installed. We're done!
			m.done = true
			return m, tea.Quit
		}

		// Update progress bar
		progressCmd := m.progress.SetPercent(float64(m.percentage))

		return m, tea.Batch(
		 	progressCmd,
			tea.Printf(" %s %s", checkMark, m.prog), // print success message above our program
			downloadAndInstall(),             // download the next package
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
		return doneStyle.Render(fmt.Sprintf("Done!\n"))
	}

	pkgCount := fmt.Sprintf(" %*d/%*d", w, m.percentage, w, n)

	spin := " " + m.spinner.View() + " "
	prog := m.progress.View()
	cellsAvail := max(0, m.width-lipgloss.Width(spin+prog+pkgCount))

	progName := currentStatusStyle.Render(m.prog)
	info := lipgloss.NewStyle().MaxWidth(cellsAvail).Render(progName)

	cellsRemaining := max(0, m.width-lipgloss.Width(spin+info+prog+pkgCount))
	gap := strings.Repeat(" ", cellsRemaining)

	return spin + info + gap + prog + pkgCount
}

type statusMsg string
type stillMsg string

func downloadAndInstall() tea.Cmd {
        d := time.Millisecond * time.Duration(rand.Intn(2500))
	return tea.Tick(d, func(t time.Time) tea.Msg {
	   // Through SSH get the status json here.
	   // Send it via stillMsg or rename it to statusMsg
	   // as a single type.
	   // Parse the json up in the update and do your thing
	   return stillMsg("")
	})
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func runTeaProgram() {
   if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
    }
}
