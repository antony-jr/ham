package get

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func runQuestionTeaProgram(resp *ResponseT, question string, placeholder string) error {
	p := tea.NewProgram(questionModel(resp, question, placeholder))
	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}

type (
	errMsg error
)

type ResponseT struct {
	answer   string
	err      error
	required bool
	secret   bool
}

func NewQuestionResponse(req bool, sec bool) *ResponseT {
	return &ResponseT{
		answer:   "",
		err:      nil,
		required: req,
		secret:   sec,
	}
}

type QuestionModel struct {
	question  string
	textInput textinput.Model
	resp      *ResponseT
}

func questionModel(resp *ResponseT, ques string, holder string) QuestionModel {
	ti := textinput.New()
	ti.Placeholder = holder
	ti.Focus()
	ti.CharLimit = 250
	ti.Width = 20
	if resp.secret {
		ti.EchoMode = textinput.EchoPassword
	}

	return QuestionModel{
		question:  ques,
		textInput: ti,
		resp:      resp,
	}
}

func (m QuestionModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m QuestionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
		case tea.KeyEsc:
			if m.resp.required {
				m.resp.err = errors.New("User Did Not Answer Question")
			}
			return m, tea.Quit

		case tea.KeyEnter:
			m.resp.answer = m.textInput.Value()
			if len(m.resp.answer) == 0 && m.resp.required {
				m.resp.err = errors.New("User Did Not Answer Question")
			}
			return m, tea.Quit
		}

	case errMsg:
		m.resp.err = errMsg(msg)
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m QuestionModel) View() string {
	return fmt.Sprintf(
		"   %s\n\n%s\n\n%s",
		m.question,
		"   "+m.textInput.View(),
		"    (esc to quit)",
	) + "\n"
}
