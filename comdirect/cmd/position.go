package cmd

import (
	"encoding/csv"

	"github.com/jsattler/go-comdirect/pkg/comdirect"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	positionCmd = &cobra.Command{
		Use:   "position",
		Short: "list depot position information",
		Args:  cobra.MinimumNArgs(1),
		Run:   position,
	}
)

func position(cmd *cobra.Command, args []string) {
	client := initClient()
	ctx, cancel := contextWithTimeout()
	defer cancel()

	positions, err := client.DepotPositions(ctx, args[0])
	if err != nil {
		return
	}
	switch formatFlag {
	case "json":
		writeJSON(positions)
	case "markdown":
		printPositionsTable(positions)
	case "csv":
		printPositionsCSV(positions)
	default:
		printPositionsTable(positions)
	}
}

func printPositionsCSV(positions *comdirect.DepotPositions) {
	table := csv.NewWriter(getOutputBuffer())
	table.Write([]string{"POSITION ID", "WKN", "QUANTITY", "CURRENT PRICE", "PREVDAY %", "PURCHASE %", "PURCHASE", "CURRENT"})
	for _, d := range positions.Values {
		table.Write([]string{d.PositionId, d.Wkn, d.Quantity.Value, d.CurrentPrice.Price.Value, d.ProfitLossPrevDayRel, d.ProfitLossPurchaseRel, d.PurchaseValue.Value, d.CurrentValue.Value})
	}
	table.Flush()
}

func printPositionsTable(depots *comdirect.DepotPositions) {
	table := tablewriter.NewWriter(getOutputBuffer())
	table.SetHeader([]string{"POSITION ID", "WKN", "QUANTITY", "CURRENT PRICE", "PREVDAY %", "PURCHASE %", "PURCHASE", "CURRENT"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	for _, d := range depots.Values {
		table.Append([]string{d.PositionId, d.Wkn, d.Quantity.Value, d.CurrentPrice.Price.Value, d.ProfitLossPrevDayRel, d.ProfitLossPurchaseRel, d.PurchaseValue.Value, d.CurrentValue.Value})
	}
	table.Render()
}
