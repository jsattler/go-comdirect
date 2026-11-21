package cmd

import (
	"encoding/csv"
	"github.com/jsattler/go-comdirect/pkg/comdirect"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"
	"os"
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
	table := csv.NewWriter(os.Stdout)
	table.Write(reportsHeader)
	for _, r := range reports.Values {
		var balance string
		if r.Balance.Balance.Value == "" {
			balance = formatAmountValue(r.Balance.PrevDayValue)
		} else {
			balance = formatAmountValue(r.Balance.Balance)
		}
		table.Write([]string{r.ProductID, r.ProductType, balance})
	}
	table.Flush()
}

func printReportsTable(reports *comdirect.Reports) {
	table := newTable(reportsHeader, tw.AlignLeft, tw.AlignLeft, tw.AlignRight)
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
