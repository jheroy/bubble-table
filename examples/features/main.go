// This file contains a full demo of most available features, for both testing
// and for reference
package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

// "os"

var listBorder = lipgloss.Border{
	Top:    "─",
	Left:   "│",
	Right:  "│",
	Bottom: "─",

	TopRight:    "╮",
	TopLeft:     "╭",
	BottomRight: "╯",
	BottomLeft:  "╰",

	MiddleTop:    "┬",
	MiddleLeft:   "├",
	MiddleRight:  "┤",
	MiddleBottom: "┴",
	// InnerJunction:  "╫",
	Middle: "┼",

	// InnerDivider: "│",
}

var docStyle = lipgloss.NewStyle().Margin(1, 2).BorderStyle(listBorder).BorderForeground(lipgloss.Color("63"))

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type listModel struct {
	list list.Model
}

func (m listModel) Init() tea.Cmd {
	return nil
}

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m listModel) View() string {
	return docStyle.Render(m.list.View())
}

func main() {
	items := []list.Item{
		item{title: "Raspberry Pi’s", desc: "I have ’em all over my house"},
		item{title: "Nutella", desc: "It's good on toast"},
		item{title: "Bitter melon", desc: "It cools you down"},
		item{title: "Nice socks", desc: "And by that I mean socks without holes"},
		item{title: "Eight hours of sleep", desc: "I had this once"},
		item{title: "Cats", desc: "Usually"},
		item{title: "Plantasia, the album", desc: "My plants love it too"},
		item{title: "Pour over coffee", desc: "It takes forever to make though"},
		item{title: "VR", desc: "Virtual reality...what is there to say?"},
		item{title: "Noguchi Lamps", desc: "Such pleasing organic forms"},
		item{title: "Linux", desc: "Pretty much the best OS"},
		item{title: "Business school", desc: "Just kidding"},
		item{title: "Pottery", desc: "Wet clay is a great feeling"},
		item{title: "Shampoo", desc: "Nothing like clean hair"},
		item{title: "Table tennis", desc: "It’s surprisingly exhausting"},
		item{title: "Milk crates", desc: "Great for packing in your extra stuff"},
		item{title: "Afternoon tea", desc: "Especially the tea sandwich part"},
		item{title: "Stickers", desc: "The thicker the vinyl the better"},
		item{title: "20° Weather", desc: "Celsius, not Fahrenheit"},
		item{title: "Warm light", desc: "Like around 2700 Kelvin"},
		item{title: "The vernal equinox", desc: "The autumnal equinox is pretty good too"},
		item{title: "Gaffer’s tape", desc: "Basically sticky fabric"},
		item{title: "Terrycloth", desc: "In other words, towel fabric"},
	}

	m := listModel{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "My Fave Things"
	// return m

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

const (
	columnKeyID          = "id"
	columnKeyName        = "name"
	columnKeyDescription = "description"
	columnKeyCount       = "count"
)

const (
	fixedVerticalMargin = 6
)

type state int

const (
	initializing state = iota
	ready
)

var (
	// customBorder = table.Border{
	// 	Top:    "─",
	// 	Left:   "│",
	// 	Right:  "│",
	// 	Bottom: "─",
	//
	// 	TopRight:    "╮",
	// 	TopLeft:     "╭",
	// 	BottomRight: "╯",
	// 	BottomLeft:  "╰",
	//
	// 	TopJunction:    "╥",
	// 	LeftJunction:   "├",
	// 	RightJunction:  "┤",
	// 	BottomJunction: "╨",
	// 	// InnerJunction:  "╫",
	// 	InnerJunction:  "┼",
	//
	// 	InnerDivider: "║",
	// }
	customBorder = table.Border{
		Top:    "─",
		Left:   "│",
		Right:  "│",
		Bottom: "─",

		TopRight:    "╮",
		TopLeft:     "╭",
		BottomRight: "╯",
		BottomLeft:  "╰",

		TopJunction:    "┬",
		LeftJunction:   "├",
		RightJunction:  "┤",
		BottomJunction: "┴",
		// InnerJunction:  "╫",
		InnerJunction: "┼",

		InnerDivider: "│",
	}
)

type Model struct {
	tableModel table.Model

	databaseList listModel

	state state

	totalWidth  int
	totalHeight int

	// Table dimensions
	horizontalMargin int
	verticalMargin   int
}

func NewModel() Model {
	columns := []table.Column{
		table.NewColumn(columnKeyID, "ID", 5).WithStyle(
			lipgloss.NewStyle().
				Faint(true).
				Foreground(lipgloss.Color("#88f")).
				Align(lipgloss.Center)),
		table.NewColumn(columnKeyName, "Name", 10),
		table.NewColumn(columnKeyDescription, "Description", 30),
		table.NewColumn(columnKeyCount, "#", 5),
		table.NewColumn(columnKeyCount, "test", 5),
		table.NewColumn(columnKeyCount, "test1", 5),
		table.NewColumn(columnKeyCount, "test2", 5),
		table.NewColumn(columnKeyCount, "test3", 5),
		table.NewColumn(columnKeyCount, "test4", 5),
		table.NewColumn(columnKeyCount, "test5", 5),
		table.NewColumn(columnKeyCount, "test6", 5),
		table.NewColumn(columnKeyCount, "test1", 5),
		table.NewColumn(columnKeyCount, "test2", 5),
		table.NewColumn(columnKeyCount, "test3", 5),
		table.NewColumn(columnKeyCount, "test4", 5),
		table.NewColumn(columnKeyCount, "test5", 5),
		table.NewColumn(columnKeyCount, "test6", 5),
		table.NewColumn(columnKeyCount, "test1", 5),
		table.NewColumn(columnKeyCount, "test2", 5),
		table.NewColumn(columnKeyCount, "test3", 5),
		table.NewColumn(columnKeyCount, "test4", 5),
		table.NewColumn(columnKeyCount, "test5", 5),
		table.NewColumn(columnKeyCount, "test6", 5),
		table.NewColumn(columnKeyCount, "test6", 5),
	}

	rows := []table.Row{
		table.NewRow(table.RowData{
			columnKeyID: "abc",
			// Missing name
			columnKeyDescription: "The first table entry, ever",
			columnKeyCount:       4,
		}),
		table.NewRow(table.RowData{
			columnKeyID:          "123",
			columnKeyName:        "Oh no",
			columnKeyDescription: "Super bold!",
			columnKeyCount:       17,
		}).WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)),
		table.NewRow(table.RowData{
			columnKeyID: "def",
			// Apply a style to this cell
			columnKeyName:        table.NewStyledCell("Styled", lipgloss.NewStyle().Foreground(lipgloss.Color("#8ff"))),
			columnKeyDescription: "This is a really, really, really long description that will get cut off",
			columnKeyCount:       table.NewStyledCell(0, lipgloss.NewStyle().Faint(true)),
		}),
		table.NewRow(table.RowData{
			columnKeyID:          "spg",
			columnKeyName:        "Page 2",
			columnKeyDescription: "Second page",
			columnKeyCount:       2,
		}),
		table.NewRow(table.RowData{
			columnKeyID:          "spg2",
			columnKeyName:        "Page 2.1",
			columnKeyDescription: "Second page again",
			columnKeyCount:       4,
		}),
		table.NewRow(table.RowData{
			columnKeyID:          "spg",
			columnKeyName:        "Page 2",
			columnKeyDescription: "Second page",
			columnKeyCount:       2,
		}),
		table.NewRow(table.RowData{
			columnKeyID:          "spg2",
			columnKeyName:        "Page 2.1",
			columnKeyDescription: "Second page again",
			columnKeyCount:       4,
		}),
		table.NewRow(table.RowData{
			columnKeyID:          "spg",
			columnKeyName:        "Page 2",
			columnKeyDescription: "Second page",
			columnKeyCount:       2,
		}),
		table.NewRow(table.RowData{
			columnKeyID:          "spg2",
			columnKeyName:        "Page 2.1",
			columnKeyDescription: "Second page again",
			columnKeyCount:       4,
		}),
		table.NewRow(table.RowData{
			columnKeyID:          "spg",
			columnKeyName:        "Page 2",
			columnKeyDescription: "Second page",
			columnKeyCount:       2,
		}),
		table.NewRow(table.RowData{
			columnKeyID:          "spg2",
			columnKeyName:        "Page 2.1",
			columnKeyDescription: "Second page again",
			columnKeyCount:       4,
		}),
		table.NewRow(table.RowData{
			columnKeyID:          "spg",
			columnKeyName:        "Page 2",
			columnKeyDescription: "Second page",
			columnKeyCount:       2,
		}),
		table.NewRow(table.RowData{
			columnKeyID:          "spg2",
			columnKeyName:        "Page 2.1",
			columnKeyDescription: "Second page again",
			columnKeyCount:       4,
		}),
		table.NewRow(table.RowData{
			columnKeyID:          "spg",
			columnKeyName:        "Page 2",
			columnKeyDescription: "Second page",
			columnKeyCount:       2,
		}),
		table.NewRow(table.RowData{
			columnKeyID:          "spg2",
			columnKeyName:        "Page 2.1",
			columnKeyDescription: "Second page again",
			columnKeyCount:       4,
		}),
	}

	// Start with the default key map and change it slightly, just for demoing
	keys := table.DefaultKeyMap()
	keys.RowDown.SetKeys("j", "down", "s")
	keys.RowUp.SetKeys("k", "up", "w")

	model := Model{
		// Throw features in... the point is not to look good, it's just reference!
		tableModel: table.New(columns).
			WithRows(rows).
			HeaderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)).
			SelectableRows(true).
			Focused(true).
			Border(customBorder).
			WithKeyMap(keys).
			WithStaticFooter("Footer!").
			// WithPageSize(3).
			WithSelectedText(" ", "✓").
			WithBaseStyle(
				lipgloss.NewStyle().
					BorderForeground(lipgloss.Color("#a38")).
					Foreground(lipgloss.Color("#a7a")).
					Align(lipgloss.Left),
			).
			SortByAsc(columnKeyID).
			WithHorizontalFreezeColumnCount(1).
			// WithMaxTotalWidth(66).
			// WithTargetWidth(33).
			WithMissingDataIndicatorStyled(table.StyledCell{
				Style: lipgloss.NewStyle().Foreground(lipgloss.Color("#faa")),
				Data:  "-",
			}),
		// databaseList: NewListModel(),
	}

	model.updateFooter()

	return model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) updateFooter() {
	highlightedRow := m.tableModel.HighlightedRow()

	footerText := fmt.Sprintf(
		"Pg. %d/%d - Currently looking at ID: %s",
		m.tableModel.CurrentPage(),
		m.tableModel.MaxPages(),
		highlightedRow.Data[columnKeyID],
	)

	m.tableModel = m.tableModel.WithStaticFooter(footerText)
}

func (m *Model) recalculateTable() {
	m.tableModel = m.tableModel.
		// WithTargetWidth(m.totalWidth).
		WithMaxTotalWidth(m.totalWidth)
	// WithMinimumHeight(m.calculateHeight())

}

func (m Model) calculateWidth() int {
	return m.totalWidth - m.horizontalMargin
}

func (m Model) calculateHeight() int {
	return m.totalHeight - m.verticalMargin - fixedVerticalMargin
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.tableModel, cmd = m.tableModel.Update(msg)
	cmds = append(cmds, cmd)

	// We control the footer text, so make sure to update it
	m.updateFooter()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			cmds = append(cmds, tea.Quit)

		case "i":
			m.tableModel = m.tableModel.WithHeaderVisibility(!m.tableModel.GetHeaderVisibility())
		}
		// m.recalculateTable()
		m.tableModel.CalcLastShowColumnIndex()
	case tea.WindowSizeMsg:
		m.state = ready
		m.totalWidth = msg.Width
		m.totalHeight = msg.Height
		m.tableModel = m.tableModel.WithMaxTotalWidth(m.totalWidth)
		m.tableModel = m.tableModel.WithPageSize(m.calculateHeight() - fixedVerticalMargin)
		m.tableModel.CalcLastShowColumnIndex()
		//
		// m.recalculateTable()
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.state != ready {
		return "Initializing..."
	}

	body := strings.Builder{}

	body.WriteString("A (chaotic) table demo with all features enabled!\n")
	body.WriteString("Press left/right or page up/down to move pages\n")
	body.WriteString("Press 'i' to toggle the header visibility\n")
	body.WriteString("Press space/enter to select a row, q or ctrl+c to quit\n")

	selectedIDs := []string{}

	for _, row := range m.tableModel.SelectedRows() {
		// Slightly dangerous type assumption but fine for demo
		selectedIDs = append(selectedIDs, row.Data[columnKeyID].(string))
	}

	body.WriteString(fmt.Sprintf("columns:%v offset:%v maxcolindex: %v current %v:%v,last: %v\n",
		m.tableModel.Columns(),
		m.tableModel.HorizontalScrollOffsetCol(),
		m.tableModel.MaxHorizontalColumnIndex(),
		m.tableModel.CurrentRowIndex(),
		m.tableModel.CurrentColumnIndex(),
		m.tableModel.LastShowColumnIndex()))

	body.WriteString(m.tableModel.View())

	body.WriteString("\n")

	return body.String()
}

func main1() {
	p := tea.NewProgram(NewModel(), tea.WithAltScreen())

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
