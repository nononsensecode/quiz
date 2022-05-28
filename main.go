package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"nononsensecode.com/quiz/quiz"
)

func main() {
	q, err := os.OpenFile("quiz.csv", os.O_RDONLY, 0600)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer q.Close()

	var questions []*quiz.Quiz
	scanner := bufio.NewScanner(q)
	for scanner.Scan() {
		csv := strings.Split(scanner.Text(), ",")
		qz, err := quiz.New(csv)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		questions = append(questions, qz)
	}

	fmt.Println("Here are your questions:")
	reader := bufio.NewReader(os.Stdin)
	for i, qz := range questions {
		qz.Present(i + 1)
		for alive := true; alive; {
			timer := time.NewTimer(qz.Timeout())
			select {
			case result := <-display(reader, qz):
				fmt.Println(result)
				timer.Stop()
			case <-timer.C:
				alive = false
				fmt.Println("Time over. Over to next question...")
			}
		}
	}
}

func display(reader *bufio.Reader, qz *quiz.Quiz) chan string {
	var s chan string
	answer, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
	qz.ReadAnswer(answer)
	s <- "Answered"
	return s
}
