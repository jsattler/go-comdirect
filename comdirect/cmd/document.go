package cmd

import (
	"fmt"
	"github.com/jsattler/go-comdirect/pkg/comdirect"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
	"time"
)

var (
	documentCmd = &cobra.Command{
		Use:   "document",
		Short: "list and download postbox documents",
		Run:   document,
	}
)

func document(cmd *cobra.Command, args []string) {
	client := initClient()
	ctx, cancel := contextWithTimeout()
	defer cancel()
	options := comdirect.EmptyOptions()
	options.Add(comdirect.PagingFirstQueryKey, indexFlag)
	options.Add(comdirect.PagingCountQueryKey, countFlag)
	documents, err := client.Documents(ctx, options)
	if err != nil {
		log.Fatal(err)
	}

	filtered := &comdirect.Documents{Values: []comdirect.Document{}, Paging: documents.Paging}
	if len(args) != 0 {
		for _, d := range documents.Values {
			for _, a := range args {
				if a == d.DocumentID {
					filtered.Values = append(filtered.Values, d)
				}
			}
		}
	} else {
		filtered.Values = documents.Values
	}

	if downloadFlag {
		download(client, filtered)
	} else {
		switch formatFlag {
		case "json":
			printJSON(filtered)
		case "markdown":
			printDocumentTable(filtered)
		default:
			printDocumentTable(filtered)
		}
	}

}

func download(client *comdirect.Client, documents *comdirect.Documents) {
	ctx, cancel := contextWithTimeout()
	defer cancel()
	for i, d := range documents.Values {
		err := client.DownloadDocument(ctx, &d, folderFlag)
		if err != nil {
			log.Fatal("failed to download document: ", err)
		}
		// TODO: think about a better solution to limit download requests to 10/sec
		if i % 10 == 0  && i != 0{
			time.Sleep(900 * time.Millisecond)
		}
		fmt.Printf("Download complete for document with ID %s\n", d.DocumentID)
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
