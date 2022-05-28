package quiz

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Quiz struct {
	question   string
	options    []string
	answer     int
	userAnswer int
	timeLimit  time.Duration
	opts       opts
}

type opts []string

func (o opts) IsExist(a string) int {
	for i, opt := range o {
		if opt == strings.ToUpper(a) {
			return i
		}
	}
	return -1
}

func New(a []string) (q *Quiz, err error) {
	q = new(Quiz)
	err = q.Marshal(a)
	if err != nil {
		return
	}
	q.opts = opts{"A", "B", "C", "D"}
	return
}

func (q *Quiz) Marshal(a []string) (err error) {
	if len(a) != 7 {
		err = fmt.Errorf("invalid quiz. Length should be 8, but got %d", len(a))
		return
	}

	q.question = a[0]

	q.options = append(q.options, a[1:5]...)

	q.answer, err = strconv.Atoi(a[5])
	if err != nil {
		err = fmt.Errorf("invalid answer: %w", err)
		return
	}

	q.timeLimit, err = time.ParseDuration(a[6])
	if err != nil {
		err = fmt.Errorf("question time limit is invalid: %w", err)
		return
	}
	return
}

func (q Quiz) Present(i int) {
	fmt.Printf("%d. %s ?\n", i, q.question)
	for i, opt := range q.opts {
		fmt.Printf("%s) %s \n", opt, q.options[i])
	}
	fmt.Printf("Answer: ")
}

func (q *Quiz) ReadAnswer(a string) {
	q.userAnswer = q.opts.IsExist(a)
}

func (q Quiz) Timeout() time.Duration {
	return q.timeLimit
}
