package cmd

import (
	"strings"
	"testing"
)

func TestTableRender(t *testing.T) {
	table := newTable([]string{"ID", "TYPE", "BALANCE"}, alignLeft, alignLeft, alignRight)
	table.Append([]string{"A1", "ACCOUNT", "1.234,00"})
	table.Append([]string{"", "TOTAL", "10,00"})
	table.Append([]string{"Ü2"})
	table.SetCaption("3 out of 5")

	var b strings.Builder
	table.render(&b)

	// Same output as github.com/olekukonko/tablewriter v0.0.5, which was used before.
	want := `| ID |  TYPE   | BALANCE  |
|----|---------|----------|
| A1 | ACCOUNT | 1.234,00 |
|    | TOTAL   |    10,00 |
| Ü2 |         |          |
3 out of 5
`
	if got := b.String(); got != want {
		t.Errorf("unexpected output:\n%s\nwant:\n%s", got, want)
	}
}
