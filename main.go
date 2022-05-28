package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"nononsensecode.com/quiz/quiz"
)

func main() {
	csvFile, err := os.OpenFile("quiz.csv", os.O_RDONLY, 0600)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer csvFile.Close()

	questions := quiz.PrepareQuestions(csvFile)

	totalTimeout := questions.TotalTimeout()
	fmt.Printf("Here are your questions (All should be answered in %s): \n", totalTimeout)
	totTimer := time.NewTimer(totalTimeout)

	report := make(chan string)
	go askQuestions(questions, report)

	select {
	case <-totTimer.C:
		fmt.Println("Quiz time is over")
		fmt.Println(questions.Result())
		return
	case r := <-report:
		fmt.Println(r)
		return
	}
}

func askQuestions(questions quiz.Questions, c chan string) {
	reader := bufio.NewReader(os.Stdin)
	answer := make(chan string)

	for i, q := range questions {
		q.Display(i + 1)

		timer := time.NewTimer(q.Timeout())
		go getAnswer(reader, answer)

		select {
		case <-timer.C:
			fmt.Println("Answering timed out")
			continue
		case a := <-answer:
			q.ReadAnswer(a)
			timer.Stop()
		}
	}

	c <- questions.Result()
}

func getAnswer(reader *bufio.Reader, a chan string) {
	answer, err := reader.ReadString('\n')
	if err != nil {
		a <- "INVALID"
		return
	}
	a <- answer
}
