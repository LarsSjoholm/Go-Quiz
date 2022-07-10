package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var quizCmd = &cobra.Command{
	Use:   "quiz",
	Short: "A simple Math-Quiz",
	Long:  `A simple Math-Quiz written in Go.`,
	Run: func(cmd *cobra.Command, args []string) {
		r := result{}
		q := quiz{}
		getQuiz(&q)
		RunQuiz(&q)
		postQuiz(&q, &r)
		ShowResult(&r)
	},
}

func init() {
	go rootCmd.AddCommand(quizCmd)
}

type quiz struct {
	Title     string     `json:"title"`
	Questions []question `json:"questions"`
	Scores    []int      `json:"-"`
}

type question struct {
	Problem         string   `json:"problem"`
	Choices         []choice `json:"choices"`
	ChoosenChoiceID string   `json:"choosenChoiceID"`
}

type choice struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type result struct {
	Score      int    `json:"score"`
	Percentile string `json:"percentile"`
}

func getQuiz(q *quiz) {
	url := "http://localhost:9090/quiz"
	responseBytes := getQuizData(url)
	if err := json.Unmarshal(responseBytes, &q); err != nil {
		log.Printf("Could not unmarshal response - %v", err)
	}
}

func postQuiz(q *quiz, r *result) {
	url := "http://localhost:9090/quiz"
	responseBytes := postQuizData(url, q)
	if err := json.Unmarshal(responseBytes, r); err != nil {
		log.Printf("Could not unmarshal response - %v", err)
	}

	return
}

func RunQuiz(q *quiz) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\n%v\n", q.Title)
	for i, b := range q.Questions {
		fmt.Printf("\n%v\n\n", b.Problem)
		for _, c := range b.Choices {
			fmt.Printf("%v, %v\n", c.ID, c.Text)
		}
		fmt.Printf("\nChoose an option(A,B,C,D): ")
		a, _ := reader.ReadString('\n')
		q.Questions[i].ChoosenChoiceID = strings.Trim(strings.Replace(a, "\r\n", "", -1), "\n")
	}
	return
}

func ShowResult(r *result) {
	fmt.Printf("\nScore: %v\n", r.Score)
	fmt.Println(r.Percentile)
}

func getQuizData(baseAPI string) []byte {
	request, err := http.NewRequest(http.MethodGet, baseAPI, nil)
	if err != nil {
		log.Printf("Could not request a quiz - %v", err)
	}
	request.Header.Add("Accept", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Printf("Could not make a request - %v", err)
	}

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Could not read a response body - %v", err)
	}

	return responseBytes
}

func postQuizData(baseAPI string, q *quiz) []byte {
	reqBody, err := json.Marshal(q)
	if err != nil {
		print(err)
	}
	resp, err := http.Post(baseAPI,
		"application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		print(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print(err)
	}

	return body
}
