package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type problem struct {
	question string
	answer   string
}

func main() {
	csvfilename := flag.String("csv", "problems.csv", "a csv file in the fromat of 'question, answers'")
	timeLimit := flag.Int("limit", 30, "the time limit of the quiz in seconds")
	flag.Parse()
	file, err := os.Open(*csvfilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open filename %s\n", *csvfilename))
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("failed to parse the provided csv file")
	}
	problems := parseLines(lines)
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	correctAnswer := 0
	for _, p := range problems {
		fmt.Printf("what is the solution of %s?\n", p.question)
		answerChannel := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerChannel <- answer
		}()
		select {
		case <-timer.C:
			fmt.Printf("\nYou scored %d out of %d\n", correctAnswer, len(problems))
			return
		case answer := <-answerChannel:
			if answer == p.answer {
				correctAnswer++
			}
		}
	}
}

func parseLines(lines [][]string) []problem {
	r := make([]problem, len(lines))
	for i, v := range lines {
		r[i] = problem{
			question: v[0],
			answer:   strings.TrimSpace(v[1]),
		}
	}
	return r
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
