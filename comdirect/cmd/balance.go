package cmd

import (
	"encoding/csv"
	"log"

	"github.com/jsattler/go-comdirect/pkg/comdirect"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	balanceCmd = &cobra.Command{
		Use:   "balance",
		Short: "list account balances",
		Run:   balance,
	}
)

func balance(cmd *cobra.Command, args []string) {
	client := initClient()
	ctx, cancel := contextWithTimeout()
	defer cancel()
	balances, err := client.Balances(ctx)
	if err != nil {
		log.Fatal(err)
	}

	switch formatFlag {
	case "json":
		writeJSON(balances)
	case "markdown":
		printBalanceTable(balances)
	case "csv":
		printBalanceCSV(balances)
	default:
		printBalanceTable(balances)
	}
}

func printBalanceCSV(balances *comdirect.AccountBalances) {
	table := csv.NewWriter(getOutputBuffer())
	table.Write([]string{"ID", "TYPE", "IBAN", "BALANCE"})
	for _, a := range balances.Values {
		table.Write([]string{a.AccountId, a.Account.AccountType.Text, a.Account.Iban, a.Balance.Value})
	}
	table.Flush()
}

func printBalanceTable(account *comdirect.AccountBalances) {
	table := tablewriter.NewWriter(getOutputBuffer())
	table.SetHeader([]string{"ID", "TYPE", "IBAN", "BALANCE"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	for _, a := range account.Values {
		table.Append([]string{a.AccountId, a.Account.AccountType.Text, a.Account.Iban, a.Balance.Value})
	}
	table.Render()
}
