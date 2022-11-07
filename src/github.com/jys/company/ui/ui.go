package ui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
	"time"
)

type Model struct {
	requests  []string
	index     int
	width     int
	height    int
	spinner   spinner.Model
	progress  progress.Model
	done      bool
	startTime time.Time
}

var (
	currentReqNameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	doneStyle           = lipgloss.NewStyle().Margin(1, 2)
	checkMark           = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")
)

func NewModel(requests []string) Model {
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)
	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	return Model{
		requests:  requests,
		spinner:   s,
		progress:  p,
		startTime: time.Now(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tea.Tick(0, func(time time.Time) tea.Msg {
		return m.NextTick()
	}), m.spinner.Tick)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit
		}
	case DoneRequestMsg:
		if m.index >= len(m.requests)-1 {
			m.done = true
			return m, tea.Quit
		}

		progressCmd := m.progress.SetPercent(float64(m.index) / float64(len(m.requests)-3))

		temp := m.index

		m.index++

		return m, tea.Batch(
			progressCmd,
			tea.Printf("%s %s", checkMark, m.requests[temp]))
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

func (m Model) View() string {
	n := len(m.requests)
	w := lipgloss.Width(fmt.Sprintf("%d", n))

	if m.done {
		return doneStyle.Render(fmt.Sprintf("%d 건의 요청 성공\n총 소요 시간 : %s\n", n-2, time.Since(m.startTime)))
	}

	reqCount := fmt.Sprintf(" %*d/%*d", w, m.index, w, n-2)

	spin := m.spinner.View() + " "
	prog := m.progress.View()
	cellsAvail := max(0, m.width-lipgloss.Width(spin+prog+reqCount))

	reqName := currentReqNameStyle.Render(m.requests[m.index])
	info := lipgloss.NewStyle().MaxWidth(cellsAvail).Render("요청중 " + reqName)

	cellsRemaining := max(0, m.width-lipgloss.Width(spin+info+prog+reqCount))
	gap := strings.Repeat(" ", cellsRemaining)

	return spin + info + gap + prog + reqCount
}

func (m Model) NextTick() DoneRequestMsg {
	return "done"
}

type DoneRequestMsg string

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
