package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	fileName string
	limit    int
	shuffle  bool
)

type Problem struct {
	ques string
	ans  string
}

func init() {
	flag.StringVar(&fileName, "csv", "problems.csv", "Name of CSV file")
	flag.IntVar(&limit, "limit", 30, "Time limit in seconds")
	flag.BoolVar(&shuffle, "shuffle", false, "Whether to shuffle the quiz or not")
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func parseCSV() [][]string {
	f, err := os.Open(fileName)
	checkError(err)
	defer f.Close()
	csvReader := csv.NewReader(f)
	csvData, err := csvReader.ReadAll()
	checkError(err)
	return csvData
}

func preprocessData(csvData [][]string) []Problem {
	quizData := make([]Problem, len(csvData))
	for i, data := range csvData {
		quizData[i] = Problem{
			ques: data[0],
			ans:  strings.TrimSpace(data[1]),
		}
	}
	if shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(quizData), func(i, j int) {
			quizData[i], quizData[j] = quizData[j], quizData[i]
		})
	}
	return quizData
}

func takeInput(ansCh chan int) {
	var ansUser int
	fmt.Scanln(&ansUser)
	ansCh <- ansUser
}

func main() {
	flag.Parse()
	csvData := parseCSV()
	quizData := preprocessData(csvData)
	score, length := 0, len(quizData)
	ansCh := make(chan int)
	fmt.Println("Press <ENTER> to start")
	fmt.Scanln()
	timer := time.NewTimer(time.Duration(limit) * time.Second)
	for i, data := range quizData {
		fmt.Printf("Problem #%d: %s = ", i+1, data.ques)
		go takeInput(ansCh)
		select {
		case <-timer.C:
			fmt.Printf("\nYou scored %d out of %d.\n", score, length)
			return
		case ansUser := <-ansCh:
			ans, _ := strconv.Atoi(data.ans)
			if ansUser == ans {
				score++
			}
		}
	}
	fmt.Printf("\nYou scored %d out of %d.\n", score, length)
}
