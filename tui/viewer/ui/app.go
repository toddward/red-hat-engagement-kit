package ui

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/toddward/red-hat-engagement-kit/tui/viewer/protocol"
)

// ViewType identifies the current view
type ViewType int

const (
	ViewMenu ViewType = iota
	ViewActivity
	ViewInput
	ViewArtifacts
	ViewChecklists
)

// BashMsg wraps a message from the bash layer
type BashMsg struct {
	Response protocol.Response
}

// App is the root model
type App struct {
	width  int
	height int

	sidebar    Sidebar
	menu       Menu
	activity   Activity
	input      Input
	palette    Palette
	artifacts  ArtifactBrowser
	checklists ChecklistBrowser

	currentView ViewType
	showPalette bool
	engagement  string
	skills      []protocol.Skill
	engagements []protocol.Engagement
	agents      []protocol.Agent

	pendingCmds map[string]string

	bashIn  io.Writer
	bashOut io.Reader
	scanner *bufio.Scanner
}

// NewApp creates a new app
func NewApp(bashIn io.Writer, bashOut io.Reader) App {
	menu := NewMenu("Red Hat Engagement Kit", []MenuItem{
		{Key: "1", Label: "Skills", Action: "show_skills"},
		{Key: "2", Label: "Agents", Action: "show_agents"},
		{Key: "3", Label: "Engagements", Action: "show_engagements"},
		{Key: "4", Label: "Artifacts", Action: "show_artifacts"},
		{Key: "5", Label: "Checklists", Action: "show_checklists"},
	})

	return App{
		sidebar:     NewSidebar(),
		menu:        menu,
		activity:    NewActivity(),
		input:       NewInput(),
		palette:     NewPalette(),
		artifacts:   NewArtifactBrowser(),
		checklists:  NewChecklistBrowser(),
		currentView: ViewMenu,
		pendingCmds: make(map[string]string),
		bashIn:      bashIn,
		bashOut:     bashOut,
		scanner:     bufio.NewScanner(bashOut),
	}
}

func (a App) Init() tea.Cmd {
	return tea.Batch(
		a.sendCommand("init", nil),
		a.listenToBash(),
	)
}

var cmdCounter int

func (a *App) sendCommand(cmd string, args interface{}) tea.Cmd {
	return func() tea.Msg {
		cmdCounter++
		id := fmt.Sprintf("cmd-%d", cmdCounter)
		a.pendingCmds[id] = cmd

		command := protocol.Command{
			Cmd: cmd,
			ID:  id,
		}
		if args != nil {
			argsJSON, _ := json.Marshal(args)
			command.Args = argsJSON
		}

		data, _ := json.Marshal(command)
		fmt.Fprintln(a.bashIn, string(data))
		return nil
	}
}

func (a *App) listenToBash() tea.Cmd {
	return func() tea.Msg {
		if a.scanner.Scan() {
			line := a.scanner.Text()
			var resp protocol.Response
			if err := json.Unmarshal([]byte(line), &resp); err == nil {
				return BashMsg{Response: resp}
			}
		}
		return nil
	}
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.sidebar.SetSize(SidebarWidth, msg.Height)
		mainWidth := msg.Width - SidebarWidth - 2
		a.menu.SetSize(mainWidth, msg.Height)
		a.activity.SetSize(mainWidth, msg.Height)
		a.input.SetSize(mainWidth, msg.Height)
		a.artifacts.SetSize(mainWidth, msg.Height)
		a.checklists.SetSize(mainWidth, msg.Height)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			if a.currentView == ViewActivity && a.activity.running {
				cmds = append(cmds, a.sendCommand("cancel", nil))
				return a, tea.Batch(cmds...)
			}
			return a, tea.Quit
		case "q":
			if a.currentView == ViewMenu && !a.showPalette {
				return a, tea.Quit
			}
		case "/":
			if !a.showPalette {
				a.showPalette = true
				a.palette.Open()
				return a, nil
			}
		case "esc":
			if a.showPalette {
				a.showPalette = false
				a.palette.Close()
				return a, nil
			}
			if a.currentView == ViewArtifacts {
				if a.artifacts.IsViewing() {
					a.artifacts.CloseContent()
					return a, nil
				}
				a.currentView = ViewMenu
				return a, nil
			}
			if a.currentView == ViewChecklists {
				if a.checklists.IsViewing() {
					a.checklists.CloseDetail()
					return a, nil
				}
				a.currentView = ViewMenu
				return a, nil
			}
			if a.currentView == ViewMenu {
				if a.menu.PopMenu() {
					return a, nil
				}
				return a, nil
			}
			if a.currentView != ViewMenu {
				a.currentView = ViewMenu
				return a, nil
			}
		case "enter":
			if a.showPalette {
				if item := a.palette.Selected(); item != nil {
					a.showPalette = false
					a.palette.Close()
					return a.handleAction(item.Action)
				}
			}
			if a.currentView == ViewMenu {
				if item := a.menu.Selected(); item != nil {
					return a.handleAction(item.Action)
				}
			}
			if a.currentView == ViewInput {
				val := a.input.Value()
				if val != "" {
					cmds = append(cmds, a.sendCommand("user_input", map[string]string{"text": val}))
					a.currentView = ViewActivity
					return a, tea.Batch(cmds...)
				}
			}
			if a.currentView == ViewArtifacts && !a.artifacts.IsViewing() {
				if a.artifacts.SelectedIsFile() {
					return a, a.sendCommand("read_artifact", map[string]string{"path": a.artifacts.SelectedPath()})
				}
			}
			if a.currentView == ViewChecklists && !a.checklists.IsViewing() {
				name := a.checklists.SelectedName()
				if name != "" {
					return a, a.sendCommand("get_checklist", map[string]string{"name": name})
				}
			}
		case " ":
			if a.currentView == ViewChecklists && a.checklists.IsViewing() {
				item := a.checklists.SelectedItem()
				clName := a.checklists.ChecklistName()
				if item != nil && clName != "" {
					return a, tea.Batch(
						a.sendCommand("toggle_checklist", map[string]string{
							"name": clName,
							"line": fmt.Sprintf("%d", item.Line),
						}),
						a.sendCommand("get_checklist", map[string]string{"name": clName}),
					)
				}
			}
		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			if a.currentView == ViewMenu && !a.showPalette {
				for _, item := range a.menu.Items() {
					if item.Key == msg.String() {
						return a.handleAction(item.Action)
					}
				}
			}
		}

		if a.showPalette {
			var cmd tea.Cmd
			a.palette, cmd = a.palette.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			switch a.currentView {
			case ViewMenu:
				var cmd tea.Cmd
				a.menu, cmd = a.menu.Update(msg)
				cmds = append(cmds, cmd)
			case ViewActivity:
				var cmd tea.Cmd
				a.activity, cmd = a.activity.Update(msg)
				cmds = append(cmds, cmd)
			case ViewInput:
				var cmd tea.Cmd
				a.input, cmd = a.input.Update(msg)
				cmds = append(cmds, cmd)
			case ViewArtifacts:
				var cmd tea.Cmd
				a.artifacts, cmd = a.artifacts.Update(msg)
				cmds = append(cmds, cmd)
			case ViewChecklists:
				var cmd tea.Cmd
				a.checklists, cmd = a.checklists.Update(msg)
				cmds = append(cmds, cmd)
			}
		}

	case BashMsg:
		cmds = append(cmds, a.handleBashResponse(msg.Response))
		cmds = append(cmds, a.listenToBash())
	}

	return a, tea.Batch(cmds...)
}

func (a *App) handleAction(action string) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	parts := strings.SplitN(action, ":", 2)
	actionType := parts[0]
	actionArg := ""
	if len(parts) > 1 {
		actionArg = parts[1]
	}

	switch actionType {
	case "show_skills":
		items := make([]MenuItem, len(a.skills))
		for i, s := range a.skills {
			items[i] = MenuItem{
				Label:       "/" + s.Name,
				Description: s.Description,
				Action:      "run_skill:" + s.Name,
			}
		}
		a.menu.PushMenu("Skills", items)
	case "show_agents":
		items := make([]MenuItem, len(a.agents))
		for i, ag := range a.agents {
			items[i] = MenuItem{
				Label:       ag.Name,
				Description: ag.Role + " (" + ag.Model + ")",
				Action:      "run_agent:" + ag.Name,
			}
		}
		a.menu.PushMenu("Agents", items)
	case "show_engagements":
		items := make([]MenuItem, len(a.engagements))
		for i, e := range a.engagements {
			label := e.Slug
			if !e.HasContext {
				label += " (no context)"
			}
			items[i] = MenuItem{
				Label:  label,
				Action: "select_engagement:" + e.Slug,
			}
		}
		items = append(items, MenuItem{
			Label:  "Create new engagement",
			Action: "run_skill:setup",
		})
		a.menu.PushMenu("Engagements", items)
	case "show_artifacts":
		a.currentView = ViewArtifacts
		cmds = append(cmds, a.sendCommand("list_artifacts", map[string]string{"engagement": a.engagement}))
	case "show_checklists":
		a.currentView = ViewChecklists
		cmds = append(cmds, a.sendCommand("list_checklists", nil))
	case "select_engagement":
		a.engagement = actionArg
		a.sidebar.SetEngagement(actionArg)
		// Pop back to root menu
		for a.menu.PopMenu() {
		}
		cmds = append(cmds, a.sendCommand("get_phase", map[string]string{"engagement": actionArg}))
		cmds = append(cmds, a.sendCommand("set_state", map[string]string{"key": "lastEngagement", "value": actionArg}))
	case "run_skill":
		a.activity.Clear()
		a.activity.SetRunning(true)
		a.currentView = ViewActivity
		cmds = append(cmds, a.sendCommand("execute_skill", map[string]string{
			"skill":      actionArg,
			"engagement": a.engagement,
		}))
	case "run_agent":
		a.input.SetPrompt("Enter prompt for "+actionArg+":", nil)
		a.currentView = ViewInput
	}

	return a, tea.Batch(cmds...)
}

func (a *App) handleBashResponse(resp protocol.Response) tea.Cmd {
	cmdName := a.pendingCmds[resp.ID]
	delete(a.pendingCmds, resp.ID)

	switch resp.Type {
	case "response":
		switch cmdName {
		case "init":
			var initResp protocol.InitResponse
			if err := json.Unmarshal(resp.Payload, &initResp); err == nil {
				a.skills = initResp.Skills
				a.engagements = initResp.Engagements
				a.agents = initResp.Agents
				if len(a.engagements) == 1 {
					a.engagement = a.engagements[0].Slug
					a.sidebar.SetEngagement(a.engagement)
					return a.sendCommand("get_phase", map[string]string{"engagement": a.engagement})
				}
				a.buildPaletteItems()
			}
		case "get_phase":
			var phaseInfo protocol.PhaseInfo
			if err := json.Unmarshal(resp.Payload, &phaseInfo); err == nil {
				a.sidebar.SetPhase(phaseInfo)
			}
		case "set_state":
			// no-op
		case "list_artifacts":
			var treeResp struct {
				Tree []protocol.ArtifactNode `json:"tree"`
			}
			if err := json.Unmarshal(resp.Payload, &treeResp); err == nil {
				a.artifacts.SetTree(treeResp.Tree)
			}
		case "read_artifact":
			var contentResp struct {
				Content string `json:"content"`
			}
			if err := json.Unmarshal(resp.Payload, &contentResp); err == nil {
				a.artifacts.SetContent(contentResp.Content)
			}
		case "list_checklists":
			var clResp struct {
				Checklists []struct {
					Name string `json:"name"`
				} `json:"checklists"`
			}
			if err := json.Unmarshal(resp.Payload, &clResp); err == nil {
				names := make([]string, len(clResp.Checklists))
				for i, cl := range clResp.Checklists {
					names[i] = cl.Name
				}
				a.checklists.SetNames(names)
			}
		case "get_checklist":
			var cl protocol.Checklist
			if err := json.Unmarshal(resp.Payload, &cl); err == nil && cl.Name != "" {
				a.checklists.SetChecklist(cl)
			}
		case "toggle_checklist":
			// refresh handled by the follow-up get_checklist command
		}

	case "event":
		var event protocol.StreamEvent
		if err := json.Unmarshal(resp.Payload, &event); err == nil {
			if event.Event == protocol.EventQuestion {
				a.input.SetPrompt(event.Text, event.Options)
				a.currentView = ViewInput
			} else {
				a.activity.AddEvent(event)
			}
		}

	case "error":
		var errPayload struct {
			Message string `json:"message"`
		}
		if err := json.Unmarshal(resp.Payload, &errPayload); err == nil {
			a.activity.AddEvent(protocol.StreamEvent{
				Event: protocol.EventError,
				Text:  errPayload.Message,
			})
		}
	}

	return nil
}

func (a *App) buildPaletteItems() {
	items := make([]PaletteItem, 0)
	for _, s := range a.skills {
		items = append(items, PaletteItem{
			Name:        "/" + s.Name,
			Description: s.Description,
			Action:      "run_skill:" + s.Name,
			Category:    "Skills",
		})
	}
	for _, ag := range a.agents {
		items = append(items, PaletteItem{
			Name:        ag.Name,
			Description: ag.Role,
			Action:      "run_agent:" + ag.Name,
			Category:    "Agents",
		})
	}
	a.palette.SetItems(items)
}

func (a App) View() string {
	sidebar := a.sidebar.View()

	var main string
	switch a.currentView {
	case ViewMenu:
		main = a.menu.View()
	case ViewActivity:
		main = a.activity.View()
	case ViewInput:
		main = a.input.View()
	case ViewArtifacts:
		main = a.artifacts.View()
	case ViewChecklists:
		main = a.checklists.View()
	}

	content := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, main)

	if a.showPalette {
		paletteView := a.palette.View()
		x := (a.width - 60) / 2
		y := 5
		content = placeOverlay(x, y, paletteView, content)
	}

	return content
}

func placeOverlay(x, y int, overlay, background string) string {
	bgLines := strings.Split(background, "\n")
	ovLines := strings.Split(overlay, "\n")

	for i, ovLine := range ovLines {
		bgY := y + i
		if bgY >= 0 && bgY < len(bgLines) {
			bgLine := bgLines[bgY]
			if x >= 0 && x < len(bgLine) {
				bgLines[bgY] = bgLine[:x] + ovLine
			}
		}
	}

	return strings.Join(bgLines, "\n")
}
