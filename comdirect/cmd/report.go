package cmd

import (
	"log"
	"os"

	"github.com/jsattler/go-comdirect/pkg/comdirect"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	reportsHeader = []string{"ID", "TYPE", "BALANCE"}

	reportCmd = &cobra.Command{
		Use:   "report",
		Short: "list aggregated account information",
		Run:   report,
	}
)

func report(cmd *cobra.Command, args []string) {
	client := initClient()
	ctx, cancel := contextWithTimeout()
	defer cancel()
	reports, err := client.Reports(ctx)
	if err != nil {
		return
	}
	switch formatFlag {
	case "json":
		printJSON(reports)
	case "markdown":
		printReportsTable(reports)
	case "csv":
		printReportsCSV(reports)
	default:
		printReportsTable(reports)
	}
}

func printReportsCSV(reports *comdirect.Reports) {
	t, err := getCSVTemplate("report.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	err = t.ExecuteTemplate(os.Stdout, t.Name(), reports)
	if err != nil {
		log.Fatal(err)
	}
}

func printReportsTable(reports *comdirect.Reports) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(reportsHeader)
	table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_RIGHT})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	for _, r := range reports.Values {
		var balance string
		if r.Balance.Balance.Value == "" {
			balance = formatAmountValue(r.Balance.PrevDayValue)
		} else {
			balance = formatAmountValue(r.Balance.Balance)
		}
		table.Append([]string{r.ProductID, r.ProductType, balance})
	}
	table.Append([]string{"", "TOTAL", formatAmountValue(reports.ReportAggregated.BalanceEUR)})
	table.Render()
}
