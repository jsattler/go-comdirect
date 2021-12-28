package cmd

import (
	"context"
	"fmt"
	"github.com/jsattler/go-comdirect/pkg/comdirect"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

var (
	folder string

	documentCmd = &cobra.Command{
		Use:   "document",
		Short: "list and download postbox documents",
		Run:   document,
	}
)

func document(cmd *cobra.Command, args []string) {
	client := InitClient()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	options := comdirect.EmptyOptions()
	options.Add(comdirect.PagingFirstQueryKey, pageIndex)
	options.Add(comdirect.PagingCountQueryKey, pageCount)
	documents, err := client.Documents(ctx, options)
	if err != nil {
		log.Fatal(err)
	}
	if len(args) == 0 {
		printDocumentTable(documents)
	}

	context.WithTimeout(context.Background(), 20*time.Second)
	for _, d := range documents.Values {
		for _, a := range args {
			if a == d.DocumentID {
				err := client.DownloadDocument(ctx, &d, folder)
				if err != nil {
					log.Fatal("failed to download document: ", err)
				}
				fmt.Printf("Download complete for document with ID %s", d.DocumentID)
			}
		}
	}
}

func printDocumentTable(documents *comdirect.Documents) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "NAME", "DATE", "OPENED", "TYPE"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetCaption(true, fmt.Sprintf("%d out of %d", len(documents.Values), documents.Paging.Matches))
	for _, d := range documents.Values {
		name := strings.ReplaceAll(d.Name, " ", "-")
		if len(name) > 30 {
			name = name[:30]
		}
		table.Append([]string{d.DocumentID, name + "...", d.DateCreation, fmt.Sprintf("%t", d.DocumentMetaData.AlreadyRead), d.MimeType})
	}
	table.Render()
}

func formatDocumentTable(tw *tabwriter.Writer, doc comdirect.Document) {
	shortName := doc.Name
	if len(doc.Name) > 30 {
		shortName = doc.Name[:30]
	}
	fmt.Fprintf(tw, "%s\t%s\t%s...\t%t\n", doc.DocumentID, doc.DateCreation, shortName, !doc.DocumentMetaData.AlreadyRead)
}
