package frame

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/lmika/awstools/internal/dynamo-browse/ui/teamodels/utils"
)

var (
	activeHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#ffffff")).
				Background(lipgloss.Color("#4479ff"))

	inactiveHeaderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#000000")).
				Background(lipgloss.Color("#d1d1d1"))
)

// Frame is a frame that appears in the
type FrameTitle struct {
	header string
	active bool
	width  int
}

func NewFrameTitle(header string, active bool) FrameTitle {
	return FrameTitle{header, active, 0}
}

func (f *FrameTitle) SetTitle(title string) {
	f.header = title
}

func (f FrameTitle) View() string {
	return f.headerView()
}

func (f *FrameTitle) Resize(w, h int) {
	f.width = w
}

func (f FrameTitle) HeaderHeight() int {
	return lipgloss.Height(f.headerView())
}

func (f FrameTitle) headerView() string {
	style := inactiveHeaderStyle
	if f.active {
		style = activeHeaderStyle
	}

	titleText := f.header
	title := style.Render(titleText)
	line := style.Render(strings.Repeat(" ", utils.Max(0, f.width-lipgloss.Width(title))))
	return lipgloss.JoinHorizontal(lipgloss.Left, title, line)
}