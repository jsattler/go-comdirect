package cmd

import (
	"encoding/json"
	"log"
	"os"
)

// Format JSON input and write to global output file which might be stdout.
func writeJSON(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	getOutputFile().Write(b)
}

func getOutputFile() *os.File {
	if outputFile == nil {
		return os.Stdout
	}
	return outputFile
}

// Read global fileFlag variable and return a file pointer to the target file if any.
// This function return nil if the target file is stdout.
func getOutputFileFromFlag() *os.File {
	if fileFlag != "" && fileFlag != "-" {
		outputFile, err := os.OpenFile(fileFlag, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			println(err)
		} else {
			return outputFile
		}
	}
	return nil
}
