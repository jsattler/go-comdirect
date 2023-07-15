package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/jsattler/go-comdirect/comdirect/tpl"
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
	const dateLayout = "2006-01-02"
	since, err := time.Parse(dateLayout, sinceFlag)
	if err != nil {
		log.Fatalf("Failed to parse date from command line: %s", err)
	}

	client := initClient()
	var transactions = &comdirect.AccountTransactions{}

	page := 1
	pageCount, err := strconv.Atoi(countFlag)
	if err != nil {
		log.Fatalf("Can't convert string to int: %s", err)
	}
	for {
		options := comdirect.EmptyOptions()
		options.Add(comdirect.PagingCountQueryKey, fmt.Sprint(pageCount*page))
		options.Add(comdirect.PagingFirstQueryKey, fmt.Sprint(0))
		ctx, cancel := contextWithTimeout()
		defer cancel()

		transactions, err = client.Transactions(ctx, args[0], options)
		if err != nil {
			log.Fatalf("Error retrieving transactions: %e", err)
		}

		if len(transactions.Values) == 0 {
			break
		}

		lastDate, err := time.Parse(dateLayout, transactions.Values[len(transactions.Values)-1].BookingDate)
		if err != nil {
			log.Fatalf("Failed to parse date from command line: %s", err)
		}
		if transactions.Paging.Matches == len(transactions.Values) || lastDate.Before(since) {
			break
		}
		page++
	}

	transactions, err = transactions.FilterSince(since)
	if err != nil {
		log.Fatalf("Error filtering transactions by date: %e", err)
	}

	switch formatFlag {
	case "json":
		printJSON(transactions)
	case "markdown":
		printTransactionTable(transactions)
	case "csv":
		printTransactionCSV(transactions)
	default:
		printTransactionTable(transactions)
	}
}

func printJSON(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}

func printTransactionCSV(transactions *comdirect.AccountTransactions) {
	t, err := template.New("transaction.tmpl").Funcs(
		template.FuncMap{
			"formatAmountValue": formatAmountValue,
			"holderName":        holderName,
		},
	).ParseFS(tpl.Default, "transaction.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	err = t.ExecuteTemplate(os.Stdout, "transaction.tmpl", transactions)
	if err != nil {
		log.Fatal(err)
	}
}
func printTransactionTable(transactions *comdirect.AccountTransactions) {
	table := tablewriter.NewWriter(os.Stdout)
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

func holderName(holderName string) string {
	if len(holderName) > 30 {
		holderName = holderName[:30]
	} else if holderName == "" {
		holderName = "N/A"
	}
	return holderName
}
