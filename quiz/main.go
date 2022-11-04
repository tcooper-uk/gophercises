package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	inputFile := flag.String("csv", "problems.csv", "The input CSV (question,answer) of problems. The CSV should not contain a header")
	secondLimit := flag.Int64("time", 30, "The time limit (in seconds) to complete the quiz.")
	shuffle := flag.Bool("shuffle", false, "Indicates that the order of the questions should be shuffled.")

	flag.Parse()

	file, err := os.Open(*inputFile)
	defer file.Close()

	if err != nil {
		fmt.Println("unable to open input file")
		os.Exit(1)
	}

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("unable to read input csv")
		os.Exit(1)
	}

	if *shuffle {
		shuffleRecords(records)
	}

	score, total := 0, len(records)

	fmt.Println("Are you ready to start the quiz? Press any key to continue.")
	fmt.Scanln()

	done := make(chan bool)
	timeout := make(chan bool)

	go startTimer(*secondLimit, timeout)
	go runQuiz(records, &score, done)

	select {
	case <-done:
		break
	case <-timeout:
		fmt.Printf("\nYou have run out of time.\n")
		break
	}

	fmt.Printf("You scored %d out of %d\n", score, total)
}

func runQuiz(records [][]string, score *int, done chan<- bool) {

	for i, record := range records {

		fmt.Printf("Problem #%d: %s = ", i, record[0])

		var answer string
		_, err := fmt.Scan(&answer)
		if err != nil {
			return
		}

		answer = strings.TrimSpace(answer)
		answer = strings.ToLower(answer)

		if strings.Compare(answer, strings.TrimSpace(record[1])) == 0 {
			*score++
		}
	}

	done <- true
}

func startTimer(seconds int64, timeout chan<- bool) {
	s := time.Duration(seconds)
	time.Sleep(time.Second * s)
	timeout <- true
}
func shuffleRecords(records [][]string) {

	// an offset shuffle

	//min := 0
	//max := len(records)
	//rand.Seed(time.Now().UnixNano())
	//offset := rand.Intn(max-min+1) + min
	//
	//newRecords := make([][]string, max, max)
	//
	//// shuffle array
	//for i := 0; i < max; i++ {
	//	push := (i + offset) % max
	//	newRecords[push] = records[i]
	//}
	//
	//return newRecords

	// or use the builtin shuffle
	// uses Fisher-Yates shuffle: https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(records), func(i, j int) {
		records[i], records[j] = records[j], records[i]
	})
}
