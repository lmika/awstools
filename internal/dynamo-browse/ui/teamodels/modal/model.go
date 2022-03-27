package modal

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lmika/awstools/internal/dynamo-browse/ui/teamodels/utils"
	"log"
)

// Modal is a tea model which displays modes on a stack.  Only the top-level model is display and will receive
// keyboard and mouse events.
type Modal struct {
	baseMode  tea.Model
	modeStack []tea.Model
}

func New(baseMode tea.Model) Modal {
	return Modal{baseMode: baseMode}
}

func (m Modal) Init() tea.Cmd {
	return nil
}

func (m *Modal) pushMode(model tea.Model) {
	m.modeStack = append(m.modeStack, model)
	log.Printf("pusing new mode: len = %v", len(m.modeStack))
}

func (m *Modal) popMode() {
	if len(m.modeStack) > 0 {
		m.modeStack = m.modeStack[:len(m.modeStack)-1]
	}
}

func (m Modal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cc utils.CmdCollector

	switch msg := msg.(type) {
	case newModePushed:
		m.pushMode(msg)
		return m, nil
	case modePopped:
		m.popMode()
		return m, nil
	case tea.KeyMsg, tea.MouseMsg:
		// only notify top level stack
		if len(m.modeStack) > 0 {
			m.modeStack[len(m.modeStack)-1] = cc.Collect(m.modeStack[len(m.modeStack)-1].Update(msg))
		} else {
			m.baseMode = cc.Collect(m.baseMode.Update(msg))
		}
	default:
		// notify all modes of other events
		// TODO: is this right?
		m.baseMode = cc.Collect(m.baseMode.Update(msg))
		for i, s := range m.modeStack {
			m.modeStack[i] = cc.Collect(s.Update(msg))
		}
	}

	return m, cc.Cmd()
}

func (m Modal) View() string {
	// only show top level mode
	if len(m.modeStack) > 0 {
		log.Printf("viewing mode stack: len = %v", len(m.modeStack))
		return m.modeStack[len(m.modeStack)-1].View()
	}
	return m.baseMode.View()
}
