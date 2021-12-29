package cmd

import (
	"github.com/jsattler/go-comdirect/pkg/comdirect"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
)

var (
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
	printReportsTable(reports)
}

func printReportsTable(reports *comdirect.Reports) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "TYPE", "BALANCE"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetFooter([]string{"", "TOTAL", reports.ReportAggregated.BalanceEUR.Value})
	for _, r := range reports.Values {
		balance := r.Balance.Balance.Value
		if balance == "" {
			balance = r.Balance.PrevDayValue.Value
		}
		table.Append([]string{r.ProductID, r.ProductType, balance})
	}
	table.Render()
}
