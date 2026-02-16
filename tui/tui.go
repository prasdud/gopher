package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/prasdud/gopher/internal"
	"github.com/prasdud/gopher/internal/collector"
)

// MetricMsg wraps a CollectorResult as a bubbletea message.
type MetricMsg collector.CollectorResult

// Model is the bubbletea model for the TUI.
type Model struct {
	uptime   *internal.Uptime
	ram      *internal.RamDetails
	cpu      *internal.CpuPercent
	width    int
	height   int
	expanded bool
	theme    Theme
	results  chan collector.CollectorResult
	err      error
}

// NewModel creates a new TUI model.
func NewModel(results chan collector.CollectorResult, theme Theme) Model {
	return Model{
		results: results,
		theme:   theme,
	}
}

// waitForMetric returns a Cmd that blocks on the channel and sends a MetricMsg.
func waitForMetric(results chan collector.CollectorResult) tea.Cmd {
	return func() tea.Msg {
		result := <-results
		return MetricMsg(result)
	}
}

func (m Model) Init() tea.Cmd {
	return waitForMetric(m.results)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "ctrl+o":
			m.expanded = !m.expanded
			return m, nil
		case "ctrl+r":
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case MetricMsg:
		if msg.Err != nil {
			m.err = msg.Err
		} else {
			switch metric := msg.Metric.(type) {
			case *internal.Uptime:
				m.uptime = metric
			case *internal.RamDetails:
				m.ram = metric
			case *internal.CpuPercent:
				m.cpu = metric
			}
			m.err = nil
		}
		return m, waitForMetric(m.results)
	}

	return m, nil
}

func (m Model) View() string {
	// Minimum size check
	if m.width < 40 || m.height < 10 {
		msg := "Terminal too small"
		hint := "(need at least 40x10)"
		padX := max(0, (m.width-len(msg))/2)
		padY := max(0, (m.height-2)/2)
		return strings.Repeat("\n", padY) +
			strings.Repeat(" ", padX) + msg + "\n" +
			strings.Repeat(" ", max(0, (m.width-len(hint))/2)) + hint
	}

	var b strings.Builder

	// Header
	b.WriteString(RenderHeader(m.theme))
	b.WriteString("\n\n")

	if m.expanded {
		b.WriteString(m.viewExpanded())
	} else {
		b.WriteString(m.viewCompact())
	}

	// Footer
	b.WriteString("\n")
	b.WriteString(RenderFooter(m.theme, m.width))

	// Error line
	if m.err != nil {
		b.WriteString("\n")
		errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
		b.WriteString(errStyle.Render(fmt.Sprintf(" error: %v", m.err)))
	}

	return b.String()
}

func (m Model) viewCompact() string {
	var b strings.Builder

	barWidth := min(30, m.width-20)
	if barWidth < 5 {
		barWidth = 5
	}

	// CPU
	cpuLabel := m.theme.CPULabel.Render(" CPU")
	cpuValue := "  --.-% "
	var cpuBar string
	if m.cpu != nil {
		cpuValue = fmt.Sprintf(" %6s ", FormatCPU(m.cpu.Percent))
		cpuBar = ProgressBar(m.cpu.Percent, barWidth, m.theme.CPUBarFill, m.theme.CPUBarEmpty)
	} else {
		cpuBar = ProgressBar(0, barWidth, m.theme.CPUBarFill, m.theme.CPUBarEmpty)
	}
	b.WriteString(cpuLabel + m.theme.CPUValue.Render(cpuValue) + cpuBar + "\n")

	// RAM
	ramLabel := m.theme.RAMLabel.Render(" RAM")
	ramValue := "  --/-- GB "
	var ramBar string
	var ramPercent float64
	if m.ram != nil {
		ramValue = fmt.Sprintf(" %s ", FormatRAM(m.ram.UsedRam, m.ram.TotalRam))
		// Pad to align with CPU value
		for len(ramValue) < 8 {
			ramValue = " " + ramValue
		}
		if m.ram.TotalRam > 0 {
			ramPercent = float64(m.ram.UsedRam) / float64(m.ram.TotalRam) * 100
		}
		ramBar = ProgressBar(ramPercent, barWidth, m.theme.RAMBarFill, m.theme.RAMBarEmpty)
	} else {
		ramBar = ProgressBar(0, barWidth, m.theme.RAMBarFill, m.theme.RAMBarEmpty)
	}
	b.WriteString(ramLabel + m.theme.RAMValue.Render(ramValue) + ramBar + "\n")

	// Swap
	swapLabel := m.theme.RAMLabel.Render(" SWP")
	swapValue := "  --/-- GB "
	var swapBar string
	var swapPercent float64
	if m.ram != nil && m.ram.TotalSwap > 0 {
		swapValue = fmt.Sprintf(" %s ", FormatRAM(m.ram.UsedSwap, m.ram.TotalSwap))
		for len(swapValue) < 8 {
			swapValue = " " + swapValue
		}
		swapPercent = m.ram.UsedSwap / m.ram.TotalSwap * 100
		swapBar = ProgressBar(swapPercent, barWidth, m.theme.RAMBarFill, m.theme.RAMBarEmpty)
	} else {
		swapBar = ProgressBar(0, barWidth, m.theme.RAMBarFill, m.theme.RAMBarEmpty)
	}
	b.WriteString(swapLabel + m.theme.RAMValue.Render(swapValue) + swapBar + "\n")

	// Uptime
	uptimeLabel := m.theme.UptimeLabel.Render(" ⏱ ")
	uptimeValue := " --h --m --s"
	if m.uptime != nil {
		uptimeValue = " " + FormatUptime(m.uptime.Hours, m.uptime.Minutes, m.uptime.Seconds)
	}
	b.WriteString(uptimeLabel + m.theme.UptimeValue.Render(uptimeValue))

	return b.String()
}

func (m Model) viewExpanded() string {
	var b strings.Builder

	boxWidth := min(40, m.width-2)
	if boxWidth < 20 {
		boxWidth = 20
	}

	barWidth := boxWidth - 18
	if barWidth < 5 {
		barWidth = 5
	}

	wide := m.width >= 80

	// CPU box
	cpuContent := "Usage: --.-% "
	if m.cpu != nil {
		bar := ProgressBar(m.cpu.Percent, barWidth, m.theme.CPUBarFill, m.theme.CPUBarEmpty)
		cpuContent = fmt.Sprintf("Usage: %s  %s", FormatCPU(m.cpu.Percent), bar)
	}
	cpuBox := BorderedBox("CPU", cpuContent, boxWidth, m.theme.Border)

	// RAM box
	ramContent := "Used:  --/-- GB"
	if m.ram != nil {
		var ramPercent float64
		if m.ram.TotalRam > 0 {
			ramPercent = m.ram.UsedRam / m.ram.TotalRam * 100
		}
		bar := ProgressBar(ramPercent, barWidth, m.theme.RAMBarFill, m.theme.RAMBarEmpty)
		ramContent = fmt.Sprintf("Used:  %s  %s\nFree:  %.2f GB\nAvailable: %.2f GB",
			FormatRAM(m.ram.UsedRam, m.ram.TotalRam), bar,
			m.ram.FreeRam, m.ram.AvailableRam)
		// Swap info
		if m.ram.TotalSwap > 0 {
			var swapPercent float64
			swapPercent = m.ram.UsedSwap / m.ram.TotalSwap * 100
			swapBar := ProgressBar(swapPercent, barWidth, m.theme.RAMBarFill, m.theme.RAMBarEmpty)
			ramContent += fmt.Sprintf("\nSwap:  %s  %s",
				FormatRAM(m.ram.UsedSwap, m.ram.TotalSwap), swapBar)
		}
	}
	ramBox := BorderedBox("Memory", ramContent, boxWidth, m.theme.Border)

	// Uptime box
	uptimeContent := "--h --m --s"
	if m.uptime != nil {
		uptimeContent = FormatUptime(m.uptime.Hours, m.uptime.Minutes, m.uptime.Seconds)
	}
	uptimeBox := BorderedBox("Uptime", uptimeContent, boxWidth, m.theme.Border)

	if wide {
		// Side by side: CPU | RAM, then Uptime below
		cpuRendered := cpuBox
		ramRendered := ramBox
		b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, " "+cpuRendered, "  "+ramRendered))
		b.WriteString("\n")
		b.WriteString(" " + uptimeBox)
	} else {
		// Stacked
		b.WriteString(" " + cpuBox + "\n")
		b.WriteString(" " + ramBox + "\n")
		b.WriteString(" " + uptimeBox)
	}

	return b.String()
}
