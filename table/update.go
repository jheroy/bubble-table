package table

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) moveHighlightUp() {
	m.rowCursorIndex--

	if m.rowCursorIndex < 0 {
		m.rowCursorIndex = len(m.GetVisibleRows()) - 1
	}

	m.currentPage = m.expectedPageForRowIndex(m.rowCursorIndex)
}

func (m *Model) moveHighlightDown() {
	m.rowCursorIndex++

	if m.rowCursorIndex >= len(m.GetVisibleRows()) {
		m.rowCursorIndex = 0
	}

	m.currentPage = m.expectedPageForRowIndex(m.rowCursorIndex)
}

func (m *Model) moveHighlightTop() {
	m.rowCursorIndex = 0
	m.currentPage = m.expectedPageForRowIndex(m.rowCursorIndex)
}

func (m *Model) moveHighlightBottom() {
	m.rowCursorIndex = len(m.GetVisibleRows()) - 1
	m.currentPage = m.expectedPageForRowIndex(m.rowCursorIndex)
}

func (m *Model) moveHighlightLeft() {
	m.columnCursorIndex--

	if m.columnCursorIndex < 0 {
		m.columnCursorIndex = 0
		// m.columnCursorIndex = len(m.columns) - 1
	}

	if m.columnCursorIndex > 0 && m.columnCursorIndex <= m.horizontalScrollOffsetCol {
		m.scrollLeft()
	}

	// m.currentPage = m.expectedPageForRowIndex(m.rowCursorIndex)
}

func (m *Model) moveHighlightRight() {
	m.columnCursorIndex++

	if m.columnCursorIndex >= len(m.columns) {
		m.columnCursorIndex = len(m.columns) - 1
		// m.columnCursorIndex = 0
	}

	if m.columnCursorIndex >= m.lastShowColumnIndex && m.columnCursorIndex != len(m.columns)-1 {
		m.scrollRight()
	}

	// m.currentPage = m.expectedPageForRowIndex(m.rowCursorIndex)
}

func (m *Model) moveHighlightedRowStart() {
	m.columnCursorIndex = 0
	m.horizontalScrollOffsetCol = 0
}

func (m *Model) moveHighlightedRowEnd() {
	m.columnCursorIndex = len(m.columns) - 1
	m.horizontalScrollOffsetCol = m.maxHorizontalColumnIndex
}

func (m *Model) CalcLastShowColumnIndex() {
	headerStrings := []string{}

	totalRenderedWidth := 0

	headerStyles := m.styleHeaders()

	renderHeader := func(column Column, borderStyle lipgloss.Style) string {
		borderStyle = borderStyle.Inherit(column.style).Inherit(m.baseStyle)

		headerSection := limitStr(column.title, column.width)

		return borderStyle.Render(headerSection)
	}

	for columnIndex, column := range m.columns {
		var borderStyle lipgloss.Style

		if m.horizontalScrollOffsetCol > 0 && columnIndex == m.horizontalScrollFreezeColumnsCount {
			if columnIndex == 0 {
				borderStyle = headerStyles.left.Copy()
			} else {
				borderStyle = headerStyles.inner.Copy()
			}

			rendered := renderHeader(genOverflowColumnLeft(1), borderStyle)

			totalRenderedWidth += lipgloss.Width(rendered)

			// headerStrings = append(headerStrings, rendered)
		}

		if columnIndex >= m.horizontalScrollFreezeColumnsCount &&
			columnIndex < m.horizontalScrollOffsetCol+m.horizontalScrollFreezeColumnsCount {
			continue
		}

		if len(headerStrings) == 0 {
			borderStyle = headerStyles.left.Copy()
		} else if columnIndex < len(m.columns)-1 {
			borderStyle = headerStyles.inner.Copy()
		} else {
			borderStyle = headerStyles.right.Copy()
		}

		rendered := renderHeader(column, borderStyle)

		if m.maxTotalWidth != 0 {
			renderedWidth := lipgloss.Width(rendered)

			const (
				borderAdjustment = 1
				overflowColWidth = 2
			)

			targetWidth := m.maxTotalWidth - overflowColWidth

			if columnIndex == len(m.columns)-1 {
				// If this is the last header, we don't need to account for the
				// overflow arrow column
				targetWidth = m.maxTotalWidth
			}

			if totalRenderedWidth+renderedWidth > targetWidth {
				overflowWidth := m.maxTotalWidth - totalRenderedWidth - borderAdjustment
				overflowStyle := genOverflowStyle(headerStyles.right, overflowWidth)
				overflowColumn := genOverflowColumnRight(overflowWidth)
				m.lastShowColumnIndex = columnIndex

				overflowStr := renderHeader(overflowColumn, overflowStyle)

				headerStrings = append(headerStrings, overflowStr)

				break
			}

			totalRenderedWidth += renderedWidth
		}

		headerStrings = append(headerStrings, rendered)
	}

}

func (m *Model) toggleSelect() {
	if !m.selectableRows || len(m.GetVisibleRows()) == 0 {
		return
	}

	rows := make([]Row, len(m.GetVisibleRows()))
	copy(rows, m.GetVisibleRows())

	currentSelectedState := rows[m.rowCursorIndex].selected

	rows[m.rowCursorIndex].selected = !currentSelectedState

	m.rows = rows
	m.visibleRowCacheUpdated = false

	m.appendUserEvent(UserEventRowSelectToggled{
		RowIndex:   m.rowCursorIndex,
		IsSelected: !currentSelectedState,
	})
}

func (m Model) updateFilterTextInput(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, m.keyMap.FilterBlur) {
			m.filterTextInput.Blur()
		}
	}
	m.filterTextInput, cmd = m.filterTextInput.Update(msg)
	m.pageFirst()
	m.visibleRowCacheUpdated = false

	return m, cmd
}

// This is a series of Matches tests with minimal logic
//
//nolint:cyclop
func (m *Model) handleKeypress(msg tea.KeyMsg) {
	previousRowIndex := m.rowCursorIndex

	if key.Matches(msg, m.keyMap.RowDown) {
		m.moveHighlightDown()
	}

	if key.Matches(msg, m.keyMap.RowUp) {
		m.moveHighlightUp()
	}

	if key.Matches(msg, m.keyMap.RowLeft) {
		m.moveHighlightLeft()
	}

	if key.Matches(msg, m.keyMap.RowRight) {
		m.moveHighlightRight()
	}

	if key.Matches(msg, m.keyMap.RowStart) {
		m.moveHighlightedRowStart()
	}

	if key.Matches(msg, m.keyMap.RowEnd) {
		m.moveHighlightedRowEnd()
	}

	if key.Matches(msg, m.keyMap.RowTop) {
		m.moveHighlightTop()
	}

	if key.Matches(msg, m.keyMap.RowBottom) {
		m.moveHighlightBottom()
	}

	// if key.Matches(msg, m.keyMap.RowUp) {
	// 	m.moveHighlightUp()
	// }

	if key.Matches(msg, m.keyMap.RowSelectToggle) {
		m.toggleSelect()
	}

	if key.Matches(msg, m.keyMap.PageDown) {
		m.pageDown()
	}

	if key.Matches(msg, m.keyMap.PageUp) {
		m.pageUp()
	}

	if key.Matches(msg, m.keyMap.PageFirst) {
		m.pageFirst()
	}

	if key.Matches(msg, m.keyMap.PageLast) {
		m.pageLast()
	}

	if key.Matches(msg, m.keyMap.Filter) {
		m.filterTextInput.Focus()
		m.appendUserEvent(UserEventFilterInputFocused{})
	}

	if key.Matches(msg, m.keyMap.FilterClear) {
		m.visibleRowCacheUpdated = false
		m.filterTextInput.Reset()
	}

	if key.Matches(msg, m.keyMap.ScrollRight) {
		m.scrollRight()
	}

	if key.Matches(msg, m.keyMap.ScrollLeft) {
		m.scrollLeft()
	}

	if m.rowCursorIndex != previousRowIndex {
		m.appendUserEvent(UserEventHighlightedIndexChanged{
			PreviousRowIndex: previousRowIndex,
			SelectedRowIndex: m.rowCursorIndex,
		})
	}
}

// Update responds to input from the user or other messages from Bubble Tea.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	m.clearUserEvents()

	if !m.focused {
		return m, nil
	}

	if m.filterTextInput.Focused() {
		var cmd tea.Cmd
		m, cmd = m.updateFilterTextInput(msg)

		if !m.filterTextInput.Focused() {
			m.appendUserEvent(UserEventFilterInputUnfocused{})
		}

		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.handleKeypress(msg)
	}

	return m, nil
}
