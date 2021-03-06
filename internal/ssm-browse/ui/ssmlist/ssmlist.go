package ssmlist

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lmika/audax/internal/dynamo-browse/ui/teamodels/frame"
	"github.com/lmika/audax/internal/dynamo-browse/ui/teamodels/layout"
	"github.com/lmika/audax/internal/ssm-browse/models"
	table "github.com/lmika/go-bubble-table"
)

type Model struct {
	frameTitle frame.FrameTitle
	table      table.Model

	parameters *models.SSMParameters

	w, h int
}

func New(style frame.Style) *Model {
	frameTitle := frame.NewFrameTitle("SSM: /", true, style)
	table := table.New(table.SimpleColumns{"name", "type", "value"}, 0, 0)

	return &Model{
		frameTitle: frameTitle,
		table:      table,
	}
}

func (m *Model) SetPrefix(newPrefix string) {
	m.frameTitle.SetTitle("SSM: " + newPrefix)
}

func (m *Model) SetParameters(parameters *models.SSMParameters) {
	m.parameters = parameters
	cols := table.SimpleColumns{"name", "type", "value"}

	newTbl := table.New(cols, m.w, m.h-m.frameTitle.HeaderHeight())
	newRows := make([]table.Row, len(parameters.Items))
	for i, r := range parameters.Items {
		newRows[i] = itemTableRow{r}
	}
	newTbl.SetRows(newRows)

	m.table = newTbl
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	//var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "i", "up":
			m.table.GoUp()
			return m, m.emitNewSelectedParameter()
		case "k", "down":
			m.table.GoDown()
			return m, m.emitNewSelectedParameter()
		}
		//m.table, cmd = m.table.Update(msg)
		//return m, cmd
	}
	return m, nil
}

func (m *Model) emitNewSelectedParameter() tea.Cmd {
	return func() tea.Msg {
		if row, ok := m.table.SelectedRow().(itemTableRow); ok {
			return NewSSMParameterSelected(&(row.item))
		}

		return nil
	}
}

func (m *Model) CurrentParameter() *models.SSMParameter {
	if row, ok := m.table.SelectedRow().(itemTableRow); ok {
		return &(row.item)
	}

	return nil
}

func (m *Model) View() string {
	return lipgloss.JoinVertical(lipgloss.Top, m.frameTitle.View(), m.table.View())
}

func (m *Model) Resize(w, h int) layout.ResizingModel {
	m.w, m.h = w, h
	m.frameTitle.Resize(w, h)
	m.table.SetSize(w, h-m.frameTitle.HeaderHeight())
	return m
}
