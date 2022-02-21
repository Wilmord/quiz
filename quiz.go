package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func cleanup() {
	fmt.Println("Time is up!!")
	fmt.Printf("Asked: %d Corrects : %d\n", questions, score)
}

var score int = 0
var questions int = 0

func countQuestions() {
	questions++
}

func countScore() {
	score++
}

func quiz() {
	// open file
	f, err := os.Open("problem.csv")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to quiz!")
	fmt.Println("You have 30 secs to finish it!")
	fmt.Println("Please enter a button to start timer:")
	start, _ := reader.ReadString('\n')
	_ = start

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGINT)
	go func() {
		<-c
		cleanup()
		os.Exit(0)
	}()

	go func() {
		timer := time.NewTimer(5 * time.Second)
		<-timer.C
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Printf("What is %v sir?\n", rec[0])
		countQuestions()
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)
		if strings.Compare(rec[1], text) == 0 {
			countScore()
		}

	}
}

func main() {
	quiz()
}
