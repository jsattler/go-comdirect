package cmd

import (
	"encoding/json"
	"log"
	"os"
)

func printJSON(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	outputFile.Write(b)
}

func getWriteTarget() *os.File {
	t := os.Stdout
	if fileFlag != "" && fileFlag != "-" {
		outputFile, err := os.OpenFile(fileFlag, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			println(err)
		} else {
			t = outputFile
		}
	}
	return t
}
