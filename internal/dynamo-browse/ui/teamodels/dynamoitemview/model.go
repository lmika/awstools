package dynamoitemview

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lmika/awstools/internal/dynamo-browse/models"
	"github.com/lmika/awstools/internal/dynamo-browse/ui/teamodels/frame"
	"github.com/lmika/awstools/internal/dynamo-browse/ui/teamodels/layout"
)

type Model struct {
	ready      bool
	frameTitle frame.FrameTitle
	viewport   viewport.Model
	w, h       int

	// model state
	currentResultSet *models.ResultSet
	selectedItem     models.Item
}

func New() Model {
	return Model{
		frameTitle: frame.NewFrameTitle("Item", false),
		viewport:   viewport.New(100, 100),
	}
}

func (Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case NewItemSelected:
		m.currentResultSet = msg.ResultSet
		m.selectedItem = msg.Item
		m.updateViewportToSelectedMessage()
		return m, nil
	}
	return m, nil
}

func (m Model) View() string {
	if !m.ready {
		return ""
	}
	return lipgloss.JoinVertical(lipgloss.Top, m.frameTitle.View(), m.viewport.View())
}

func (m Model) Resize(w, h int) layout.ResizingModel {
	m.w, m.h = w, h
	if !m.ready {
		m.viewport = viewport.New(w, h-m.frameTitle.HeaderHeight())
		m.viewport.SetContent("")
		m.ready = true
	} else {
		m.viewport.Width = w
		m.viewport.Height = h - m.frameTitle.HeaderHeight()
	}
	m.frameTitle.Resize(w, h)
	return m
}

func (m *Model) updateViewportToSelectedMessage() {
	if m.selectedItem == nil {
		m.viewport.SetContent("")
	}

	viewportContent := &strings.Builder{}
	tabWriter := tabwriter.NewWriter(viewportContent, 0, 1, 1, ' ', 0)
	for _, colName := range m.currentResultSet.Columns {
		switch colVal := m.selectedItem[colName].(type) {
		case nil:
			break
		case *types.AttributeValueMemberS:
			fmt.Fprintf(tabWriter, "%v\tS\t%s\n", colName, colVal.Value)
		case *types.AttributeValueMemberN:
			fmt.Fprintf(tabWriter, "%v\tN\t%s\n", colName, colVal.Value)
		default:
			fmt.Fprintf(tabWriter, "%v\t?\t%s\n", colName, "(other)")
		}
	}

	tabWriter.Flush()
	m.viewport.Width = m.w
	m.viewport.Height = m.h - m.frameTitle.HeaderHeight()
	m.viewport.SetContent(viewportContent.String())
}