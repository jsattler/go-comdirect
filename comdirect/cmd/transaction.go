package cmd

import (
	"context"
	"fmt"
	"github.com/jsattler/go-comdirect/pkg/comdirect"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"log"
	"os"
	"time"
)

var (
	transactionCmd = &cobra.Command{
		Use:   "transaction",
		Short: "list account transactions",
		Args:  cobra.MinimumNArgs(1),
		Run:   transaction,
	}
)

func transaction(cmd *cobra.Command, args []string) {
	client := InitClient()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	options := comdirect.EmptyOptions()
	options.Add(comdirect.PagingCountQueryKey, pageCount)
	options.Add(comdirect.PagingFirstQueryKey, pageIndex)
	transactions, err := client.Transactions(ctx, args[0], options)
	if err != nil {
		log.Fatal(err)
	}
	printTransactionTable(transactions)
}

func printTransactionTable(transactions *comdirect.AccountTransactions) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"REMITTER", "DEPTOR", "BOOKING DATE", "STATUS", "TYPE", "AMOUNT"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetCaption(true, fmt.Sprintf("%d out of %d", len(transactions.Values), transactions.Paging.Matches))
	for _, t := range transactions.Values {
		holderName := t.Remitter.HolderName
		if len(holderName) > 30 {
			holderName = holderName[:30]
		} else if holderName == "" {
			holderName = "UNKNOWN"
		}
		value := t.Amount.Value
		if value[0] != '-' {
			value = "+" + value
		}
		table.Append([]string{holderName, t.Creditor.HolderName, t.BookingDate, t.BookingStatus, t.TransactionType.Text, boldGreen(value)})
	}
	table.Render()
}
