package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// DirInfo holds information about a directory and its size
type DirInfo struct {
	Path     string
	Size     int64
	Children []*DirInfo
	Parent   *DirInfo
	IsDir    bool
	Error    error
}

// Model represents the application state
type model struct {
	rootDir     *DirInfo
	currentDir  *DirInfo
	cursor      int
	scanning    bool
	scanError   error
	width       int
	height      int
	startTime   time.Time
	sortBy      string // "size" or "name"
}

type scanCompleteMsg struct {
	root  *DirInfo
	err   error
}

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	dirStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)

	fileStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA"))

	sizeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA500")).
			Bold(true)

	selectedStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#3A3A3A")).
			Foreground(lipgloss.Color("#FFFFFF"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Italic(true)
)

func initialModel(path string) model {
	return model{
		scanning:   true,
		startTime:  time.Now(),
		sortBy:     "size",
	}
}

func (m model) Init() tea.Cmd {
	return scanDirectory
}

// scanDirectory scans the directory tree
func scanDirectory() tea.Msg {
	path := "."
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	root, err := calculateDirSize(path, nil)
	if err != nil && root == nil {
		return scanCompleteMsg{nil, err}
	}

	return scanCompleteMsg{root, nil}
}

// calculateDirSize recursively calculates directory sizes
func calculateDirSize(path string, parent *DirInfo) (*DirInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	dirInfo := &DirInfo{
		Path:   path,
		Parent: parent,
		IsDir:  info.IsDir(),
	}

	if !info.IsDir() {
		dirInfo.Size = info.Size()
		return dirInfo, nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		dirInfo.Error = err
		return dirInfo, nil
	}

	var totalSize int64
	for _, entry := range entries {
		childPath := filepath.Join(path, entry.Name())

		child, err := calculateDirSize(childPath, dirInfo)
		if err != nil {
			// Skip entries we can't access
			continue
		}

		dirInfo.Children = append(dirInfo.Children, child)
		totalSize += child.Size
	}

	dirInfo.Size = totalSize

	// Sort children by size (descending)
	sort.Slice(dirInfo.Children, func(i, j int) bool {
		return dirInfo.Children[i].Size > dirInfo.Children[j].Size
	})

	return dirInfo, nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case scanCompleteMsg:
		m.scanning = false
		m.rootDir = msg.root
		m.currentDir = msg.root
		m.scanError = msg.err
		return m, nil

	case tea.KeyMsg:
		if m.scanning {
			if msg.String() == "ctrl+c" || msg.String() == "q" {
				return m, tea.Quit
			}
			return m, nil
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.currentDir != nil && m.cursor < len(m.currentDir.Children)-1 {
				m.cursor++
			}

		case "enter", "right", "l":
			if m.currentDir != nil && len(m.currentDir.Children) > 0 {
				selected := m.currentDir.Children[m.cursor]
				if selected.IsDir && len(selected.Children) > 0 {
					m.currentDir = selected
					m.cursor = 0
				}
			}

		case "left", "h", "backspace":
			if m.currentDir != nil && m.currentDir.Parent != nil {
				// Find cursor position in parent
				parent := m.currentDir.Parent
				for i, child := range parent.Children {
					if child == m.currentDir {
						m.cursor = i
						break
					}
				}
				m.currentDir = parent
			}

		case "s":
			// Toggle sort
			if m.sortBy == "size" {
				m.sortBy = "name"
				m.sortChildren(m.rootDir)
			} else {
				m.sortBy = "size"
				m.sortChildren(m.rootDir)
			}
			m.cursor = 0

		case "home":
			m.currentDir = m.rootDir
			m.cursor = 0
		}
	}

	return m, nil
}

func (m *model) sortChildren(dir *DirInfo) {
	if dir == nil {
		return
	}

	if m.sortBy == "name" {
		sort.Slice(dir.Children, func(i, j int) bool {
			return strings.ToLower(filepath.Base(dir.Children[i].Path)) <
				   strings.ToLower(filepath.Base(dir.Children[j].Path))
		})
	} else {
		sort.Slice(dir.Children, func(i, j int) bool {
			return dir.Children[i].Size > dir.Children[j].Size
		})
	}

	for _, child := range dir.Children {
		m.sortChildren(child)
	}
}

func (m model) View() string {
	if m.scanning {
		elapsed := time.Since(m.startTime).Seconds()
		return fmt.Sprintf("\n  Scanning directories... %.1fs\n\n  Press 'q' to quit\n", elapsed)
	}

	if m.scanError != nil {
		return errorStyle.Render(fmt.Sprintf("\n  Error: %v\n\n", m.scanError))
	}

	if m.rootDir == nil {
		return "\n  No data available\n"
	}

	var b strings.Builder

	// Title
	title := fmt.Sprintf(" Dusty - %s ", m.currentDir.Path)
	b.WriteString(titleStyle.Render(title))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("Total Size: %s | Sort: %s\n\n",
		formatSize(m.currentDir.Size), m.sortBy))

	// Directory listing
	maxItems := m.height - 8
	if maxItems < 5 {
		maxItems = 5
	}

	startIdx := 0
	if m.cursor >= maxItems {
		startIdx = m.cursor - maxItems + 1
	}

	for i := startIdx; i < len(m.currentDir.Children) && i < startIdx+maxItems; i++ {
		child := m.currentDir.Children[i]
		name := filepath.Base(child.Path)

		// Calculate percentage
		percentage := float64(0)
		if m.currentDir.Size > 0 {
			percentage = float64(child.Size) / float64(m.currentDir.Size) * 100
		}

		// Format line
		sizeStr := formatSize(child.Size)
		percentStr := fmt.Sprintf("%5.1f%%", percentage)

		// Create bar
		barWidth := 20
		filledWidth := int(percentage / 100.0 * float64(barWidth))
		if filledWidth > barWidth {
			filledWidth = barWidth
		}
		bar := strings.Repeat("█", filledWidth) + strings.Repeat("░", barWidth-filledWidth)

		// Style the name
		var nameStr string
		if child.IsDir {
			nameStr = dirStyle.Render(name + "/")
		} else {
			nameStr = fileStyle.Render(name)
		}

		line := fmt.Sprintf("  %s %s [%s] %s",
			sizeStyle.Render(fmt.Sprintf("%10s", sizeStr)),
			percentStr,
			bar,
			nameStr)

		if i == m.cursor {
			line = selectedStyle.Render(line)
		}

		b.WriteString(line)
		b.WriteString("\n")
	}

	// Help text
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("  ↑/↓: Navigate | ←/→: Enter/Exit | s: Sort | Home: Root | q: Quit"))
	b.WriteString("\n")

	return b.String()
}

// formatSize formats bytes into human-readable format
func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func main() {
	path := "."
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	// Verify path exists
	if _, err := os.Stat(path); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(initialModel(path), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
