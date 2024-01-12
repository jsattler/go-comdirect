package cmd

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// Format JSON input and write to global output file which might be stdout.
func writeJSON(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	getOutputBuffer().Write(b)
}

func getOutputBuffer() io.Writer {
	return &outputBuffer
}

func writeToOutputFile(cmd *cobra.Command, args []string) {
	if fileFlag == "-" {
		// We have either no file flag at all or the user explicitly specified "-", so we simply write to stdout
		outputBuffer.WriteTo(os.Stdout)
	} else {
		f, err := os.OpenFile(fileFlag, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			println(err)
		} else {
			outputBuffer.WriteTo(f)
			f.Close()
		}
	}
}
