package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
)

func main() {

	var fileName = flag.String("filename", "", "File name for quiz")
	flag.Parse()

	if *fileName == "" {
		fmt.Fprintf(os.Stderr, "usage: ./quiz -filename=<filename>\n")
		os.Exit(1)
	}

	file, err := os.Open(*fileName)
	if err != nil {
		fmt.Print(err)
	}

	csvReader := csv.NewReader(file)
	records, read_err := csvReader.ReadAll()
	if read_err != nil {
		fmt.Print(read_err)
	}

	for _, record := range records {
		fmt.Println(record)
	}

}
