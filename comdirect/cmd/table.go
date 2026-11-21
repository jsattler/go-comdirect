package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

type alignment int

const (
	alignLeft alignment = iota
	alignRight
)

// table renders rows in a markdown-like format:
//
//	|  ID  | BALANCE |
//	|------|---------|
//	| A1   |   10,00 |
//	caption
//
// Headers are centered, cells are left aligned unless a column alignment is
// given. Cells are never wrapped.
type table struct {
	header  []string
	align   []alignment
	rows    [][]string
	caption string
}

// newTable creates a table with the given header. Optional alignments apply
// to the columns in order; columns without one are left aligned.
func newTable(header []string, align ...alignment) *table {
	return &table{header: header, align: align}
}

func (t *table) Append(row []string) {
	t.rows = append(t.rows, row)
}

// SetCaption sets a line of text printed below the table.
func (t *table) SetCaption(caption string) {
	t.caption = caption
}

// Render writes the table to stdout.
func (t *table) Render() {
	t.render(os.Stdout)
}

func (t *table) render(w io.Writer) {
	widths := make([]int, len(t.header))
	for i, h := range t.header {
		widths[i] = utf8.RuneCountInString(h)
	}
	for _, row := range t.rows {
		for i := 0; i < len(row) && i < len(widths); i++ {
			widths[i] = max(widths[i], utf8.RuneCountInString(row[i]))
		}
	}

	var b strings.Builder
	writeLine := func(cell func(i int) string) {
		b.WriteString("|")
		for i := range widths {
			b.WriteString(" " + cell(i) + " |")
		}
		b.WriteString("\n")
	}

	writeLine(func(i int) string { return pad(t.header[i], widths[i], center) })
	b.WriteString("|")
	for _, width := range widths {
		b.WriteString(strings.Repeat("-", width+2) + "|")
	}
	b.WriteString("\n")
	for _, row := range t.rows {
		writeLine(func(i int) string {
			var cell string
			if i < len(row) {
				cell = row[i]
			}
			return pad(cell, widths[i], t.columnAlignment(i))
		})
	}
	if t.caption != "" {
		b.WriteString(t.caption + "\n")
	}
	fmt.Fprint(w, b.String())
}

func (t *table) columnAlignment(i int) alignment {
	if i < len(t.align) {
		return t.align[i]
	}
	return alignLeft
}

// center is only used for headers and is not exposed as a column alignment.
const center alignment = -1

func pad(s string, width int, align alignment) string {
	space := width - utf8.RuneCountInString(s)
	switch align {
	case alignRight:
		return strings.Repeat(" ", space) + s
	case center:
		left := space / 2
		return strings.Repeat(" ", left) + s + strings.Repeat(" ", space-left)
	default:
		return s + strings.Repeat(" ", space)
	}
}
