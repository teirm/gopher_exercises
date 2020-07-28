// Quiz program based on input csv file
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

// Asks a question with an integer  answer
// If the user answers correctly, true is returned.
// Otherwise the function returns false
func askQuestion(question *string, answer int) bool {
	var response int

	fmt.Printf("%s\n", *question)
	_, err := fmt.Scanf("%d", &response)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failure to read input: %v\n", err)
		os.Exit(1)
	}
	return (response == answer)
}

func main() {
	// allow the user to specify the file name for the quiz
	var fileName = flag.String("filename", "", "File name for quiz")
	var timeOut = flag.Int64("timeout", 30, "Timeout for quiz")
	flag.Parse()

	// report an error if no quiz file is specified
	if *fileName == "" {
		fmt.Fprintf(os.Stderr, "usage: ./quiz -filename=<filename>\n")
		os.Exit(1)
	}

	// open the file name specified
	file, err := os.Open(*fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failure to open filename '%s': %v\n", *fileName, err)
		os.Exit(1)
	}

	// read the file in as a default csv file
	csvReader := csv.NewReader(file)
	records, readErr := csvReader.ReadAll()
	if readErr != nil {
		fmt.Fprintf(os.Stderr, "%v", readErr)
		os.Exit(1)
	}

	// print out each record
	response := make(chan bool)

	go func() {
		for _, record := range records {
			if len(record) != 2 {
				fmt.Fprintf(os.Stderr, "Bad question: %v\n", record)
				continue
			}
			question := &record[0]
			answer, err := strconv.Atoi(record[1])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failure to convert record value '%s' to int: %v\n", record[1], err)
				os.Exit(1)
			}
			response <- askQuestion(question, answer)
		}
	}()

	correct := 0
	incorrect := len(records)
	questions := len(records)

loop:
	for {
		select {
		case <-time.After(time.Duration(*timeOut) * time.Second):
			fmt.Println("TIMES UP")
			break loop
		case ans := <-response:
			if ans == true {
				correct++
				incorrect--
			}
			questions--
			if questions == 0 {
				break loop
			}
		}
	}
	fmt.Printf("Quiz results (out of %v): %v correct %v incorrect\n", len(records), correct, incorrect)
}
