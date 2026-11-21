package cmd

import (
	"github.com/jsattler/go-comdirect/pkg/comdirect"
	"github.com/spf13/cobra"
	"log"
)

var (
	accountCmd = &cobra.Command{
		Use:   "account",
		Short: "list all available accounts",
		Args:  cobra.MinimumNArgs(0),
		Run:   Account,
	}
)

func Account(cmd *cobra.Command, args []string) {
	ctx, cancel := contextWithTimeout()
	defer cancel()
	client := initClient()
	balances, err := client.Balances(ctx)
	if err != nil {
		log.Fatal(err)
	}
	switch formatFlag {
	case "json":
		printJSON(balances)
	case "markdown":
		printAccountTable(balances)
	default:
		printAccountTable(balances)
	}
}

func printAccountTable(account *comdirect.AccountBalances) {
	table := newTable([]string{"ID", "TYPE", "IBAN", "CREDIT LIMIT"})
	for _, a := range account.Values {
		table.Append([]string{a.AccountId, a.Account.AccountType.Text, a.Account.Iban, a.Account.CreditLimit.Value})
	}
	table.Render()
}
