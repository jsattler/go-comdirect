package cmd

import (
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
)

// newTable creates a table writing to stdout using the markdown-like style of
// the CLI: left and right borders, no top and bottom borders, and "|" as the
// column separator. Optional alignments are applied per column to the rows.
func newTable(header []string, alignment ...tw.Align) *tablewriter.Table {
	table := tablewriter.NewTable(os.Stdout,
		tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{
			Borders: tw.Border{Left: tw.On, Top: tw.Off, Right: tw.On, Bottom: tw.Off},
			Symbols: tw.NewSymbols(tw.StyleMarkdown),
		})),
		tablewriter.WithRowAlignmentConfig(tw.CellAlignment{Global: tw.AlignLeft, PerColumn: alignment}),
	)
	table.Header(header)
	return table
}

// setTableCaption adds a caption below the table, left aligned.
func setTableCaption(table *tablewriter.Table, text string) {
	table.Caption(tw.Caption{Text: text, Spot: tw.SpotBottomLeft, Align: tw.AlignLeft})
}
