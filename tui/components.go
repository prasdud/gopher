package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const asciiHeader = `   ___  ___  ___ _  _ ___ ___
  / __|/ _ \| _ \ || | __| _ \
 | (_ | (_) |  _/ __ | _||   /
  \___|\___/|_| |_||_|___|_|_\`

// RenderHeader returns the ASCII art header styled with the theme.
func RenderHeader(theme Theme) string {
	return theme.Header.Render(asciiHeader)
}

// RenderFooter returns the footer line styled with the theme.
func RenderFooter(theme Theme, width int) string {
	text := "made by prasdud"
	pad := ""
	if width > len(text) {
		pad = strings.Repeat(" ", width-len(text))
	}
	return theme.Footer.Render(pad + text)
}

// ProgressBar renders a horizontal bar of the given width.
// percent should be 0-100.
func ProgressBar(percent float64, width int, fillColor, emptyColor lipgloss.Color) string {
	if width < 2 {
		return ""
	}

	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}

	filled := int(percent / 100 * float64(width))
	if filled > width {
		filled = width
	}
	empty := width - filled

	fillStyle := lipgloss.NewStyle().Foreground(fillColor)
	emptyStyle := lipgloss.NewStyle().Foreground(emptyColor)

	bar := fillStyle.Render(strings.Repeat("█", filled)) +
		emptyStyle.Render(strings.Repeat("░", empty))

	return bar
}

// BorderedBox wraps content in a bordered box with a title.
func BorderedBox(title, content string, width int, borderColor lipgloss.Color) string {
	innerWidth := width - 4 // account for "│ " and " │"
	if innerWidth < 1 {
		innerWidth = 1
	}

	border := lipgloss.NewStyle().Foreground(borderColor)

	// Top border
	top := border.Render("┌─ " + title + " " + strings.Repeat("─", max(0, innerWidth-len(title)-1)) + "┐")

	// Content lines
	var lines []string
	for _, line := range strings.Split(content, "\n") {
		// Pad or truncate each line to innerWidth
		visible := lipgloss.Width(line)
		if visible < innerWidth {
			line = line + strings.Repeat(" ", innerWidth-visible)
		}
		lines = append(lines, border.Render("│ ")+line+border.Render(" │"))
	}

	// Bottom border
	bottom := border.Render("└" + strings.Repeat("─", innerWidth+2) + "┘")

	return top + "\n" + strings.Join(lines, "\n") + "\n" + bottom
}

// FormatUptime returns a human-readable uptime string.
func FormatUptime(hours, minutes, seconds int64) string {
	return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
}

// FormatRAM returns "used/total GB".
func FormatRAM(used, total int64) string {
	return fmt.Sprintf("%d/%d GB", used, total)
}

// FormatCPU returns "XX.X%".
func FormatCPU(percent float64) string {
	return fmt.Sprintf("%.1f%%", percent)
}
