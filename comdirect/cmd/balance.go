package cmd

import (
	"context"
	"github.com/jsattler/go-comdirect/pkg/comdirect"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"log"
	"os"
	"time"
)

var (
	balanceCmd = &cobra.Command{
		Use:   "balance",
		Short: "list account balances",
		Run:   balance,
	}
)

func balance(cmd *cobra.Command, args []string) {
	client := InitClient()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	balances, err := client.Balances(ctx)
	if err != nil {
		log.Fatal(err)
	}
	printBalanceTable(balances)
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
