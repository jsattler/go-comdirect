package cmd

import (
	"encoding/csv"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jsattler/go-comdirect/pkg/comdirect"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	transactionHeader = []string{"REMITTER", "DEPTOR", "BOOKING DATE", "STATUS", "TYPE", "VALUE", "UNIT"}
	transactionCmd    = &cobra.Command{
		Use:   "transaction",
		Short: "list account transactions",
		Args:  cobra.MinimumNArgs(1),
		Run:   transaction,
	}
)

func transaction(cmd *cobra.Command, args []string) {
	var transactions = &comdirect.AccountTransactions{}
	var err error
	client := initClient()
	ctx, cancel := contextWithTimeout()
	defer cancel()

	if sinceFlag == "" {
		options := comdirect.EmptyOptions()
		options.Add(comdirect.PagingCountQueryKey, countFlag)
		options.Add(comdirect.PagingFirstQueryKey, indexFlag)
		transactions, err = client.Transactions(ctx, args[0], options)
		if err != nil {
			log.Fatalf("Failed to retrieve transactions: %s", err)
		}
	} else {
		transactions = getTransactionsSince(sinceFlag, client, args[0])
	}

	switch formatFlag {
	case "json":
		writeJSON(transactions)
	case "markdown":
		printTransactionTable(transactions)
	case "csv":
		printTransactionCSV(transactions)
	default:
		printTransactionTable(transactions)
	}
}

func getTransactionsSince(since string, client *comdirect.Client, accountID string) *comdirect.AccountTransactions {
	const dateLayout = "2006-01-02"
	var transactions = &comdirect.AccountTransactions{}
	ctx, cancel := contextWithTimeout()
	defer cancel()

	s, err := time.Parse(dateLayout, since)
	if err != nil {
		log.Fatalf("Failed to parse date from command line: %s", err)
	}

	page := 1
	pageCount, err := strconv.Atoi(countFlag)
	if err != nil {
		log.Fatalf("Can't convert string to int: %s", err)
	}
	for {
		options := comdirect.EmptyOptions()
		options.Add(comdirect.PagingCountQueryKey, fmt.Sprint(pageCount*page))
		options.Add(comdirect.PagingFirstQueryKey, fmt.Sprint(0))

		transactions, err = client.Transactions(ctx, accountID, options)
		if err != nil {
			log.Fatalf("Error retrieving transactions: %e", err)
		}

		if len(transactions.Values) == 0 {
			break
		}

		lastDate, err := time.Parse(dateLayout, transactions.Values[len(transactions.Values)-1].BookingDate)
		if err != nil {
			lastDate = time.Now()

		}
		if transactions.Paging.Matches == len(transactions.Values) || lastDate.Before(s) {
			break
		}
		page++
	}

	transactions, err = transactions.FilterSince(s)
	if err != nil {
		log.Fatalf("Error filtering transactions by date: %e", err)
	}
	return transactions
}

func printTransactionCSV(transactions *comdirect.AccountTransactions) {
	table := csv.NewWriter(getOutputFile())
	table.Write(transactionHeader)
	for _, t := range transactions.Values {
		holderName := t.Remitter.HolderName
		if len(holderName) > 30 {
			holderName = holderName[:30]
		} else if holderName == "" {
			holderName = "N/A"
		}
		table.Write([]string{holderName, t.Creditor.HolderName, t.BookingDate, t.BookingStatus, t.TransactionType.Text, formatAmountValue(t.Amount), t.Amount.Unit})
	}
	table.Flush()
}
func printTransactionTable(transactions *comdirect.AccountTransactions) {
	table := tablewriter.NewWriter(getOutputFile())
	table.SetHeader(transactionHeader)
	table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_RIGHT})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetCaption(true, fmt.Sprintf("%d out of %d", len(transactions.Values), transactions.Paging.Matches))
	for _, t := range transactions.Values {
		holderName := t.Remitter.HolderName
		if len(holderName) > 30 {
			holderName = holderName[:30]
		} else if holderName == "" {
			holderName = "N/A"
		}
		table.Append([]string{holderName, t.Creditor.HolderName, t.BookingDate, t.BookingStatus, t.TransactionType.Text, formatAmountValue(t.Amount), t.Amount.Unit})
	}
	table.Render()
}
