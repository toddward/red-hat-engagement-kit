package ui

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/toddward/red-hat-engagement-kit/tui/viewer/protocol"
)

// newTestApp creates an App with in-memory buffers for testing.
// The returned buffer captures commands sent to bashIn.
func newTestApp() (App, *bytes.Buffer) {
	cmdBuf := &bytes.Buffer{}
	app := NewApp(cmdBuf, strings.NewReader(""))
	return app, cmdBuf
}

// isQuitCmd executes the tea.Cmd and checks whether it returns tea.QuitMsg.
func isQuitCmd(cmd tea.Cmd) bool {
	if cmd == nil {
		return false
	}
	msg := cmd()
	if msg == nil {
		return false
	}
	_, ok := msg.(tea.QuitMsg)
	return ok
}

// isBatchContainingQuit checks if a tea.Cmd is (or contains via Batch) a quit command.
// tea.Batch returns a BatchMsg ([]Cmd) which we can inspect.
func containsQuit(cmd tea.Cmd) bool {
	if cmd == nil {
		return false
	}
	msg := cmd()
	if msg == nil {
		return false
	}
	// Direct quit
	if _, ok := msg.(tea.QuitMsg); ok {
		return true
	}
	// BatchMsg is a slice of Cmd
	if batch, ok := msg.(tea.BatchMsg); ok {
		for _, c := range batch {
			if c == nil {
				continue
			}
			inner := c()
			if inner == nil {
				continue
			}
			if _, ok2 := inner.(tea.QuitMsg); ok2 {
				return true
			}
		}
	}
	return false
}

func TestNewApp(t *testing.T) {
	app, _ := newTestApp()

	if app.currentView != ViewMenu {
		t.Errorf("expected currentView=ViewMenu, got %d", app.currentView)
	}
	if app.showPalette {
		t.Error("expected showPalette=false")
	}
	if app.engagement != "" {
		t.Errorf("expected engagement=\"\", got %q", app.engagement)
	}

	items := app.menu.Items()
	if len(items) != 5 {
		t.Fatalf("expected 5 menu items, got %d", len(items))
	}

	expectedActions := []string{"show_skills", "show_agents", "show_engagements", "show_artifacts", "show_checklists"}
	for i, action := range expectedActions {
		if items[i].Action != action {
			t.Errorf("menu item %d: expected action %q, got %q", i, action, items[i].Action)
		}
	}
}

func TestApp_WindowResize(t *testing.T) {
	app, _ := newTestApp()

	model, _ := app.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	app = model.(App)

	if app.width != 120 {
		t.Errorf("expected width=120, got %d", app.width)
	}
	if app.height != 40 {
		t.Errorf("expected height=40, got %d", app.height)
	}
}

func TestApp_KeyQ_QuitsFromMenu(t *testing.T) {
	app, _ := newTestApp()
	app.currentView = ViewMenu
	app.showPalette = false

	_, cmd := app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

	if !containsQuit(cmd) {
		t.Error("expected tea.Quit when pressing 'q' from menu with palette closed")
	}
}

func TestApp_KeyQ_DoesNotQuitFromActivity(t *testing.T) {
	app, _ := newTestApp()
	app.currentView = ViewActivity

	_, cmd := app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

	if containsQuit(cmd) {
		t.Error("expected no quit when pressing 'q' from activity view")
	}
}

func TestApp_KeySlash_OpensPalette(t *testing.T) {
	app, _ := newTestApp()
	app.showPalette = false

	model, _ := app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	app = model.(App)

	if !app.showPalette {
		t.Error("expected showPalette=true after pressing '/'")
	}
}

func TestApp_KeyEsc_ClosesPalette(t *testing.T) {
	app, _ := newTestApp()
	app.showPalette = true

	model, _ := app.Update(tea.KeyMsg{Type: tea.KeyEscape})
	app = model.(App)

	if app.showPalette {
		t.Error("expected showPalette=false after pressing Esc with palette open")
	}
}

func TestApp_KeyEsc_ReturnsToMenuFromActivity(t *testing.T) {
	app, _ := newTestApp()
	app.currentView = ViewActivity
	app.showPalette = false
	app.showHelp = false

	model, _ := app.Update(tea.KeyMsg{Type: tea.KeyEscape})
	app = model.(App)

	if app.currentView != ViewMenu {
		t.Errorf("expected currentView=ViewMenu after Esc from activity, got %d", app.currentView)
	}
}

func TestApp_KeyCtrlC_Quits(t *testing.T) {
	app, _ := newTestApp()
	app.currentView = ViewMenu

	_, cmd := app.Update(tea.KeyMsg{Type: tea.KeyCtrlC})

	if !containsQuit(cmd) {
		t.Error("expected tea.Quit when pressing Ctrl+C from menu")
	}
}

func TestApp_KeyCtrlC_CancelsRunning(t *testing.T) {
	app, _ := newTestApp()
	app.currentView = ViewActivity
	app.activity.SetRunning(true)

	_, cmd := app.Update(tea.KeyMsg{Type: tea.KeyCtrlC})

	if containsQuit(cmd) {
		t.Error("expected no quit when pressing Ctrl+C with running activity")
	}
	// The cmd should be non-nil (sends a cancel command)
	if cmd == nil {
		t.Error("expected non-nil cmd (cancel command) when Ctrl+C with running activity")
	}
}

func TestApp_HandleAction_ShowSkills(t *testing.T) {
	app, _ := newTestApp()
	app.skills = []protocol.Skill{
		{Name: "setup", Description: "Initialize engagement"},
		{Name: "discover", Description: "Run discovery"},
	}

	model, _ := app.handleAction("show_skills")
	app = *model.(*App)

	items := app.menu.Items()
	if len(items) != 2 {
		t.Fatalf("expected 2 skill menu items, got %d", len(items))
	}
	if items[0].Label != "/setup" {
		t.Errorf("expected first item label '/setup', got %q", items[0].Label)
	}
	if items[0].Action != "run_skill:setup" {
		t.Errorf("expected first item action 'run_skill:setup', got %q", items[0].Action)
	}
	if items[1].Label != "/discover" {
		t.Errorf("expected second item label '/discover', got %q", items[1].Label)
	}
}

func TestApp_HandleAction_ShowEngagements(t *testing.T) {
	app, _ := newTestApp()
	app.engagements = []protocol.Engagement{
		{Slug: "acme-corp", HasContext: true},
		{Slug: "globex", HasContext: false},
	}

	model, _ := app.handleAction("show_engagements")
	app = *model.(*App)

	items := app.menu.Items()
	// 2 engagements + "Create new engagement"
	if len(items) != 3 {
		t.Fatalf("expected 3 engagement menu items, got %d", len(items))
	}
	if items[0].Label != "acme-corp" {
		t.Errorf("expected first item label 'acme-corp', got %q", items[0].Label)
	}
	if items[1].Label != "globex (no context)" {
		t.Errorf("expected second item label 'globex (no context)', got %q", items[1].Label)
	}
	if items[2].Label != "Create new engagement" {
		t.Errorf("expected last item label 'Create new engagement', got %q", items[2].Label)
	}
	if items[2].Action != "run_skill:setup" {
		t.Errorf("expected last item action 'run_skill:setup', got %q", items[2].Action)
	}
}

func TestApp_HandleAction_SelectEngagement(t *testing.T) {
	app, _ := newTestApp()

	model, cmd := app.handleAction("select_engagement:acme")
	app = *model.(*App)

	if app.engagement != "acme" {
		t.Errorf("expected engagement='acme', got %q", app.engagement)
	}
	// Sidebar should reflect the engagement
	view := app.sidebar.View()
	if !strings.Contains(view, "acme") {
		t.Error("expected sidebar view to contain 'acme'")
	}
	// Menu should be back at root (5 items)
	items := app.menu.Items()
	if len(items) != 5 {
		t.Errorf("expected 5 root menu items after select_engagement, got %d", len(items))
	}
	// cmd should be non-nil (get_phase + set_state commands)
	if cmd == nil {
		t.Error("expected non-nil cmd after select_engagement")
	}
}

func TestApp_HandleAction_RunSkill(t *testing.T) {
	app, _ := newTestApp()

	model, cmd := app.handleAction("run_skill:setup")
	app = *model.(*App)

	if app.currentView != ViewActivity {
		t.Errorf("expected currentView=ViewActivity, got %d", app.currentView)
	}
	if !app.activity.running {
		t.Error("expected activity.running=true")
	}
	if cmd == nil {
		t.Error("expected non-nil cmd (execute_skill command)")
	}
}

func TestApp_HandleAction_RunAgent(t *testing.T) {
	app, _ := newTestApp()

	model, _ := app.handleAction("run_agent:architect")
	app = *model.(*App)

	if app.currentView != ViewInput {
		t.Errorf("expected currentView=ViewInput, got %d", app.currentView)
	}
}

func TestApp_HandleBashResponse_InitResponse(t *testing.T) {
	app, _ := newTestApp()

	initResp := protocol.InitResponse{
		Skills: []protocol.Skill{
			{Name: "setup", Description: "Initialize"},
			{Name: "discover", Description: "Discovery"},
		},
		Engagements: []protocol.Engagement{
			{Slug: "acme", HasContext: true},
			{Slug: "globex", HasContext: false},
		},
		Agents: []protocol.Agent{
			{Name: "architect", Model: "opus", Role: "Team lead"},
		},
	}
	payload, _ := json.Marshal(initResp)

	// Register the command ID so handleBashResponse routes correctly
	app.pendingCmds["cmd-init"] = "init"
	resp := protocol.Response{Type: "response", ID: "cmd-init", Payload: payload}
	app.handleBashResponse(resp)

	if len(app.skills) != 2 {
		t.Errorf("expected 2 skills, got %d", len(app.skills))
	}
	if len(app.engagements) != 2 {
		t.Errorf("expected 2 engagements, got %d", len(app.engagements))
	}
	if len(app.agents) != 1 {
		t.Errorf("expected 1 agent, got %d", len(app.agents))
	}
	// Palette items should be built (2 skills + 1 agent = 3)
	selected := app.palette.Selected()
	// Palette should have items - we verify by checking filtered list is not empty
	if selected == nil {
		t.Error("expected palette to have items after init response")
	}
}

func TestApp_HandleBashResponse_SingleEngagement(t *testing.T) {
	app, _ := newTestApp()

	initResp := protocol.InitResponse{
		Skills: []protocol.Skill{
			{Name: "setup", Description: "Initialize"},
		},
		Engagements: []protocol.Engagement{
			{Slug: "solo-eng", HasContext: true},
		},
		Agents: []protocol.Agent{},
	}
	payload, _ := json.Marshal(initResp)

	app.pendingCmds["cmd-single"] = "init"
	resp := protocol.Response{Type: "response", ID: "cmd-single", Payload: payload}
	app.handleBashResponse(resp)

	if app.engagement != "solo-eng" {
		t.Errorf("expected auto-selected engagement 'solo-eng', got %q", app.engagement)
	}
}

func TestApp_HandleBashResponse_PhaseInfo(t *testing.T) {
	app, _ := newTestApp()
	app.sidebar.SetEngagement("acme")

	phaseInfo := protocol.PhaseInfo{
		Phase: protocol.PhaseLive,
		ArtifactCounts: protocol.ArtifactCounts{
			Discovery:    3,
			Assessments:  1,
			Deliverables: 0,
		},
	}
	payload, _ := json.Marshal(phaseInfo)

	app.pendingCmds["cmd-phase"] = "get_phase"
	resp := protocol.Response{Type: "response", ID: "cmd-phase", Payload: payload}
	app.handleBashResponse(resp)

	view := app.sidebar.View()
	if !strings.Contains(view, "Live") {
		t.Error("expected sidebar to show 'Live' phase after phase info response")
	}
}

func TestApp_HandleBashResponse_StreamEvent(t *testing.T) {
	app, _ := newTestApp()
	app.activity.SetSize(80, 24)

	event := protocol.StreamEvent{
		Event: protocol.EventAssistant,
		Text:  "Processing your request",
	}
	payload, _ := json.Marshal(event)

	resp := protocol.Response{Type: "event", Payload: payload}
	app.handleBashResponse(resp)

	if len(app.activity.entries) != 1 {
		t.Fatalf("expected 1 activity entry, got %d", len(app.activity.entries))
	}
	if app.activity.entries[0].Text != "Processing your request" {
		t.Errorf("expected entry text 'Processing your request', got %q", app.activity.entries[0].Text)
	}
}

func TestApp_HandleBashResponse_Question(t *testing.T) {
	app, _ := newTestApp()

	event := protocol.StreamEvent{
		Event:   protocol.EventQuestion,
		Text:    "What is the customer name?",
		Options: []string{"Acme", "Globex"},
	}
	payload, _ := json.Marshal(event)

	resp := protocol.Response{Type: "event", Payload: payload}
	app.handleBashResponse(resp)

	if app.currentView != ViewInput {
		t.Errorf("expected currentView=ViewInput after question event, got %d", app.currentView)
	}
}

func TestApp_HandleBashResponse_Error(t *testing.T) {
	app, _ := newTestApp()
	app.activity.SetSize(80, 24)

	errPayload := struct {
		Message string `json:"message"`
	}{Message: "skill execution failed"}
	payload, _ := json.Marshal(errPayload)

	resp := protocol.Response{Type: "error", Payload: payload}
	app.handleBashResponse(resp)

	if len(app.activity.entries) != 1 {
		t.Fatalf("expected 1 activity entry for error, got %d", len(app.activity.entries))
	}
	if app.activity.entries[0].Event != protocol.EventError {
		t.Errorf("expected EventError, got %v", app.activity.entries[0].Event)
	}
	if app.activity.entries[0].Text != "skill execution failed" {
		t.Errorf("expected error text 'skill execution failed', got %q", app.activity.entries[0].Text)
	}
}

func TestApp_KeyDispatch_PaletteWhenOpen(t *testing.T) {
	app, _ := newTestApp()

	// Populate palette with items so cursor movement is visible
	app.palette.SetItems([]PaletteItem{
		{Name: "item-a", Action: "action-a", Category: "Test"},
		{Name: "item-b", Action: "action-b", Category: "Test"},
		{Name: "item-c", Action: "action-c", Category: "Test"},
	})

	app.showPalette = true
	app.palette.Open()

	// Verify palette cursor starts at 0
	sel := app.palette.Selected()
	if sel == nil || sel.Name != "item-a" {
		t.Fatal("expected palette cursor at item-a initially")
	}

	// Press down - should move palette cursor, not menu cursor
	model, _ := app.Update(tea.KeyMsg{Type: tea.KeyDown})
	app = model.(App)

	sel = app.palette.Selected()
	if sel == nil {
		t.Fatal("expected non-nil palette selection after down")
	}
	if sel.Name != "item-b" {
		t.Errorf("expected palette cursor to move to 'item-b', got %q", sel.Name)
	}
}

func TestApp_KeyDispatch_MenuWhenClosed(t *testing.T) {
	app, _ := newTestApp()
	app.currentView = ViewMenu
	app.showPalette = false

	// Menu cursor starts at 0
	sel := app.menu.Selected()
	if sel == nil {
		t.Fatal("expected non-nil menu selection initially")
	}
	initialLabel := sel.Label

	// Press down - should move menu cursor
	model, _ := app.Update(tea.KeyMsg{Type: tea.KeyDown})
	app = model.(App)

	sel = app.menu.Selected()
	if sel == nil {
		t.Fatal("expected non-nil menu selection after down")
	}
	if sel.Label == initialLabel {
		t.Error("expected menu cursor to move after pressing down, but it stayed at the same item")
	}
}
