package tui

import (
	"encoding/json"
	"os"

	"github.com/charmbracelet/lipgloss"
)

type metricColors struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type barColors struct {
	Label    string `json:"label"`
	Value    string `json:"value"`
	BarFill  string `json:"bar_filled"`
	BarEmpty string `json:"bar_empty"`
}

type themeJSON struct {
	Name   string       `json:"name"`
	Header string       `json:"header"`
	Footer string       `json:"footer"`
	Border string       `json:"border"`
	Uptime metricColors `json:"uptime"`
	CPU    barColors    `json:"cpu"`
	RAM    barColors    `json:"ram"`
}

// Theme holds resolved lipgloss styles for the TUI.
type Theme struct {
	Header lipgloss.Style
	Footer lipgloss.Style
	Border lipgloss.Color

	UptimeLabel lipgloss.Style
	UptimeValue lipgloss.Style

	CPULabel    lipgloss.Style
	CPUValue    lipgloss.Style
	CPUBarFill  lipgloss.Color
	CPUBarEmpty lipgloss.Color

	RAMLabel    lipgloss.Style
	RAMValue    lipgloss.Style
	RAMBarFill  lipgloss.Color
	RAMBarEmpty lipgloss.Color
}

// LoadTheme reads a JSON theme file and returns a resolved Theme.
func LoadTheme(path string) (Theme, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Theme{}, err
	}

	var tj themeJSON
	if err := json.Unmarshal(data, &tj); err != nil {
		return Theme{}, err
	}

	return Theme{
		Header: lipgloss.NewStyle().Foreground(lipgloss.Color(tj.Header)),
		Footer: lipgloss.NewStyle().Foreground(lipgloss.Color(tj.Footer)),
		Border: lipgloss.Color(tj.Border),

		UptimeLabel: lipgloss.NewStyle().Foreground(lipgloss.Color(tj.Uptime.Label)),
		UptimeValue: lipgloss.NewStyle().Foreground(lipgloss.Color(tj.Uptime.Value)),

		CPULabel:    lipgloss.NewStyle().Foreground(lipgloss.Color(tj.CPU.Label)),
		CPUValue:    lipgloss.NewStyle().Foreground(lipgloss.Color(tj.CPU.Value)),
		CPUBarFill:  lipgloss.Color(tj.CPU.BarFill),
		CPUBarEmpty: lipgloss.Color(tj.CPU.BarEmpty),

		RAMLabel:    lipgloss.NewStyle().Foreground(lipgloss.Color(tj.RAM.Label)),
		RAMValue:    lipgloss.NewStyle().Foreground(lipgloss.Color(tj.RAM.Value)),
		RAMBarFill:  lipgloss.Color(tj.RAM.BarFill),
		RAMBarEmpty: lipgloss.Color(tj.RAM.BarEmpty),
	}, nil
}

// DefaultTheme returns a hardcoded fallback theme.
func DefaultTheme() Theme {
	return Theme{
		Header: lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6F61")),
		Footer: lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")),
		Border: lipgloss.Color("#444444"),

		UptimeLabel: lipgloss.NewStyle().Foreground(lipgloss.Color("#61AFEF")),
		UptimeValue: lipgloss.NewStyle().Foreground(lipgloss.Color("#ABB2BF")),

		CPULabel:    lipgloss.NewStyle().Foreground(lipgloss.Color("#E06C75")),
		CPUValue:    lipgloss.NewStyle().Foreground(lipgloss.Color("#ABB2BF")),
		CPUBarFill:  lipgloss.Color("#E06C75"),
		CPUBarEmpty: lipgloss.Color("#3E4451"),

		RAMLabel:    lipgloss.NewStyle().Foreground(lipgloss.Color("#98C379")),
		RAMValue:    lipgloss.NewStyle().Foreground(lipgloss.Color("#ABB2BF")),
		RAMBarFill:  lipgloss.Color("#98C379"),
		RAMBarEmpty: lipgloss.Color("#3E4451"),
	}
}
