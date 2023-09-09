package cmd

import (
	"log"
	"os"

	"github.com/jsattler/go-comdirect/pkg/comdirect"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	depotCmd = &cobra.Command{
		Use:   "depot",
		Short: "list basic depot information",
		Run:   depot,
	}
)

func depot(cmd *cobra.Command, args []string) {
	client := initClient()
	ctx, cancel := contextWithTimeout()
	defer cancel()

	depots, err := client.Depots(ctx)
	if err != nil {
		return
	}
	switch formatFlag {
	case "json":
		printJSON(depots)
	case "markdown":
		printDepotsTable(depots)
	case "csv":
		printDepotsCSV(depots)
	default:
		printDepotsTable(depots)
	}
}

func printDepotsCSV(depots *comdirect.Depots) {
	t, err := getCSVTemplate("depot.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	err = t.ExecuteTemplate(os.Stdout, t.Name(), depots)
	if err != nil {
		log.Fatal(err)
	}
}
func printDepotsTable(depots *comdirect.Depots) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"DEPOT ID", "DISPLAY ID", "HOLDER NAME", "CLIENT ID"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	for _, d := range depots.Values {
		table.Append([]string{d.DepotId, d.DepotDisplayId, d.HolderName, d.ClientId})
	}
	table.Render()
}
