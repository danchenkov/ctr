package main

import (
	"bufio"
	"encoding/csv"

	// "encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Record struct {
	Query       string
	Clicks      int
	Impressions int
	CTR         float64
	Position    float64
}

func main() {
	// Read the first argument, check file existence, open
	// Read the first argument or take the current folder, check the folder existence, locate most recent csv that matches the pattern, open
	csvDataFile := ""
	if len(os.Args) > 1 {
		csvDataFile = os.Args[1]
	}

	fi, err := os.Stat(csvDataFile)
	if err != nil {
		log.Fatalf("Error %v", err)
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		// do directory stuff
		fmt.Println("directory")
	case mode.IsRegular():
		// do file stuff
		fmt.Println("file")
	}

	csvFile, err := os.Open(csvDataFile)
	if err != nil {
		log.Fatalf("Error %v", err)
	}
	// fmt.Printf("csvFile: %#v", csvFile)

	reader := csv.NewReader(bufio.NewReader(csvFile))
	var records []Record
	var counter int
	var header []string
	var totalCTR float64
	var totalPosition float64
	var qualifiedCounter int

	for {
		counter++

		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		if header == nil {
			header = record
		} else {

			clicks, err := strconv.Atoi(record[1])
			if err != nil {
				log.Fatalf("Error parsing clicks on line %d", counter, err)
			}

			impressions, err := strconv.Atoi(record[2])
			if err != nil {
				log.Fatalf("Error parsing impressions on line %d", counter, err)
			}

			ctr, err := strconv.ParseFloat(record[3][:strings.Index(record[3], "%")], 10)
			if err != nil {
				log.Fatalf("Error parsing CTR on line %d", counter, err)
			}

			position, err := strconv.ParseFloat(record[4], 10)
			if err != nil {
				log.Fatalf("Error parsing position on line %d", counter, err)
			}

			if impressions > 1000 {
				qualifiedCounter++
				totalCTR += ctr / 100
				totalPosition += position
				records = append(records, Record{
					Query:       record[0],
					Clicks:      clicks,
					Impressions: impressions,
					CTR:         ctr / 100,
					Position:    position,
				})
			}
		}
	}

	// fmt.Println(header)
	// for _, r := range records {
	// 	fmt.Printf("%30s\t%d\t%d\t%.2f\t%.1f\n", r.Query, r.Clicks, r.Impressions, r.CTR, r.Position)
	// }

	// fmt.Printf("Average CTR: %.3f%%; Average Position: %.2f\n", totalCTR*100/float64(qualifiedCounter), totalPosition/float64(qualifiedCounter))
	fmt.Printf("CTR: %.2f%%; Position: %.2f\n", totalCTR*100/float64(qualifiedCounter), totalPosition/float64(qualifiedCounter))
}
