package get

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SpinModel struct {
	spinner  spinner.Model
	title    string
	quit     *bool
	done     chan bool
	quitting bool
}

func initialSpinModel(quit *bool, title string, dn chan bool) SpinModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return SpinModel{
		spinner: s,
		title:   title,
		quit:    quit,
		done:    dn,
	}
}

func (m SpinModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, waitForClose(m.done))
}

type smoothQuit string
type waitMore string

func waitForClose(done chan bool) tea.Cmd {
	d := time.Nanosecond * time.Duration(1)
	return tea.Tick(d, func(t time.Time) tea.Msg {
		isDone := <-done
		if isDone {
			return smoothQuit("quit")
		} else {
			return waitMore("wait")
		}
	})

}

func (m SpinModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case waitMore:
		return m, tea.Batch(waitForClose(m.done))

	case smoothQuit:
		*m.quit = true
		m.quitting = true
		return m, tea.Quit

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			*m.quit = false
			return m, tea.Quit
		default:
			return m, nil
		}

	case errMsg:
		return m, nil

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m SpinModel) View() string {
	if m.quitting {
		n := len(m.title) + 4
		prt := ""
		for n > 0 {
			n--
			prt = prt + " "
		}
		return prt
	}

	str := fmt.Sprintf(" %s %s", m.spinner.View(), m.title)
	return str
}

func runSpinnerTeaProgram(ok *bool, title string, dn chan bool, end chan bool) {
	m := initialSpinModel(ok, title, dn)
	p := tea.NewProgram(m)
	_, _ = p.Run()

	end <- true
}

type TUISpinnerMessenger struct {
	quitOk bool
	fin    chan bool
	end    chan bool
}

func NewTUISpinnerMessenger() *TUISpinnerMessenger {
	return &TUISpinnerMessenger{
		quitOk: true,
		fin:    make(chan bool),
		end:    make(chan bool),
	}
}

func (ctx *TUISpinnerMessenger) ShowMessage(msg string) {
	ctx.quitOk = true
	ctx.fin = make(chan bool)
	ctx.end = make(chan bool)

	go runSpinnerTeaProgram(&ctx.quitOk, msg, ctx.fin, ctx.end)
}

func (ctx *TUISpinnerMessenger) StopMessage() bool {
	ctx.fin <- true
	ret := <-ctx.end
	for ret {
		break
	}

	return ctx.quitOk
}
