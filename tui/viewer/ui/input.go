package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Input struct {
	textInput textinput.Model
	prompt    string
	options   []string
	cursor    int
	width     int
	height    int
}

func NewInput() Input {
	ti := textinput.New()
	ti.Placeholder = "Type your response..."
	ti.CharLimit = 500
	ti.Width = 60
	return Input{textInput: ti}
}

func (i *Input) SetSize(width, height int) {
	i.width = width
	i.height = height
	i.textInput.Width = width - 10
}

func (i *Input) SetPrompt(prompt string, options []string) {
	i.prompt = prompt
	i.options = options
	i.cursor = 0
	i.textInput.SetValue("")
	if len(options) == 0 {
		i.textInput.Focus()
	}
}

func (i *Input) Focus() { i.textInput.Focus() }
func (i *Input) Blur()  { i.textInput.Blur() }

func (i Input) Value() string {
	if len(i.options) > 0 && i.cursor < len(i.options) {
		return i.options[i.cursor]
	}
	return i.textInput.Value()
}

func (i Input) Update(msg tea.Msg) (Input, tea.Cmd) {
	if len(i.options) > 0 {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "up", "k":
				if i.cursor > 0 {
					i.cursor--
				}
			case "down", "j":
				if i.cursor < len(i.options)-1 {
					i.cursor++
				}
			}
		}
		return i, nil
	}
	var cmd tea.Cmd
	i.textInput, cmd = i.textInput.Update(msg)
	return i, cmd
}

func (i Input) View() string {
	var b strings.Builder

	b.WriteString(InputLabelStyle.Render(i.prompt))
	b.WriteString("\n\n")

	if len(i.options) > 0 {
		for idx, opt := range i.options {
			cursor := "  "
			style := MenuItemStyle
			if idx == i.cursor {
				cursor = "> "
				style = MenuItemSelectedStyle
			}
			b.WriteString(style.Render(cursor + opt))
			b.WriteString("\n")
		}
	} else {
		b.WriteString(InputStyle.Render(i.textInput.View()))
	}

	b.WriteString("\n\n")
	b.WriteString(HelpKeyStyle.Render("enter"))
	b.WriteString(HelpDescStyle.Render(" submit  "))
	b.WriteString(HelpKeyStyle.Render("Esc"))
	b.WriteString(HelpDescStyle.Render(" cancel"))

	return MainStyle.Render(b.String())
}
