package cmd

import (
	"log"
	"os"

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
		printJSON(balances)
	case "markdown":
		printBalanceTable(balances)
	case "csv":
		printBalanceCSV(balances)
	default:
		printBalanceTable(balances)
	}
}

func printBalanceCSV(balances *comdirect.AccountBalances) {
	t, err := getCSVTemplate("balance.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	err = t.ExecuteTemplate(os.Stdout, t.Name(), balances)
	if err != nil {
		log.Fatal(err)
	}
}

func printBalanceTable(account *comdirect.AccountBalances) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "TYPE", "IBAN", "BALANCE"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	for _, a := range account.Values {
		table.Append([]string{a.AccountId, a.Account.AccountType.Text, a.Account.Iban, a.Balance.Value})
	}
	table.Render()
}
