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
	accountCmd = &cobra.Command{
		Use:   "account",
		Short: "list all available accounts",
		Run:   Account,
	}
)

func Account(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client := InitClient()
	balances, err := client.Balances(ctx)
	if err != nil {
		log.Fatal(err)
	}
	printAccountTable(balances)
}

func init() {
	accountCmd.AddCommand(balanceCmd)
	accountCmd.AddCommand(transactionCmd)
}

func printAccountTable(account *comdirect.AccountBalances) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "TYPE", "IBAN", "CREDIT LIMIT"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	for _, a := range account.Values {
		table.Append([]string{a.AccountId, a.Account.AccountType.Text, a.Account.Iban, a.Account.CreditLimit.Value})
	}
	table.Render()
}
