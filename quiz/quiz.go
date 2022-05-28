package quiz

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Questions []*Quiz

func (q Questions) Result() string {
	total := len(q)
	answered, right, wrong := 0, 0, 0
	for _, qz := range q {
		if qz.IsAnswerRight() {
			right++
		}
		if qz.userAnswer > -1 {
			answered++
		}
	}
	wrong = total - right

	return fmt.Sprintf("Out of %d questions, %d answered. %d were right, %d were wrong",
		total, answered, right, wrong)
}

func (qs Questions) TotalTimeout() (tot time.Duration) {
	for _, q := range qs {
		tot += q.Timeout()
	}
	tot = (tot * 80) / 100
	return
}

func PrepareQuestions(f *os.File) (questions Questions) {
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		csv := strings.Split(scanner.Text(), ",")
		q, err := newQuiz(csv)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		questions = append(questions, q)
	}
	return
}

func newQuiz(a []string) (q *Quiz, err error) {
	q = new(Quiz)
	err = q.Marshal(a)
	if err != nil {
		return
	}
	q.opts = opts{"A", "B", "C", "D"}
	q.userAnswer = -1
	return
}

type Quiz struct {
	question   string
	options    []string
	answer     int
	userAnswer int
	timeLimit  time.Duration
	opts       opts
}

type opts []string

func (o opts) doesExist(a string) int {
	for i, opt := range o {
		if opt == strings.ToUpper(strings.TrimSpace(a)) {
			return i
		}
	}
	return -1
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

func (q Quiz) Display(i int) {
	fmt.Printf("%d. %s ? (%s to answer)\n", i, q.question, q.Timeout())
	for i, opt := range q.opts {
		fmt.Printf("%s) %s \n", opt, q.options[i])
	}
	fmt.Printf("Answer: ")
}

func (q *Quiz) ReadAnswer(a string) {
	q.userAnswer = q.opts.doesExist(a)
}

func (q Quiz) Timeout() time.Duration {
	return q.timeLimit
}

func (q *Quiz) IsAnswerRight() bool {
	return q.answer == q.userAnswer
}
