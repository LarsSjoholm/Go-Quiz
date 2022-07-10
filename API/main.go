package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type quiz struct {
	Title     string     `json:"title"`
	Questions []question `json:"questions"`
	Scores    []int      `json:"-"`
}

type question struct {
	Problem         string   `json:"problem"`
	Choices         []choice `json:"choices"`
	ChoosenChoiceID string   `json:"choosenChoiceID"`
	CorrectChoiceID string   `json:"-"`
}

type choice struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type result struct {
	Score      int    `json:"score"`
	Percentile string `json:"percentile"`
}

var mathQuiz = quiz{
	Title: "Math Quiz",
	Questions: []question{
		{
			Problem: "What is 1 + 1?",
			Choices: []choice{
				{
					ID:   "A",
					Text: "1",
				},
				{
					ID:   "B",
					Text: "0",
				},
				{
					ID:   "C",
					Text: "2",
				},
				{
					ID:   "D",
					Text: "4",
				},
			},
			CorrectChoiceID: "C",
		},
		{
			Problem: "What is 4 / 2?",
			Choices: []choice{
				{
					ID:   "A",
					Text: "1",
				},
				{
					ID:   "B",
					Text: "8",
				},
				{
					ID:   "C",
					Text: "6",
				},
				{
					ID:   "D",
					Text: "2",
				},
			},
			CorrectChoiceID: "D",
		},
		{
			Problem: "What is 3 * 6?",
			Choices: []choice{
				{
					ID:   "A",
					Text: "18",
				},
				{
					ID:   "B",
					Text: "12",
				},
				{
					ID:   "C",
					Text: "9",
				},
				{
					ID:   "D",
					Text: "24",
				},
			},
			CorrectChoiceID: "A",
		},
	},
	Scores: []int{1, 2, 1, 2, 0, 2, 0, 2, 1},
}

func getQuiz(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, mathQuiz)
}

func postQuizScore(c *gin.Context) {
	var q quiz
	if err := c.BindJSON(&q); err != nil {
		return
	}
	r, err := reviewQuiz(&q)
	if err != nil {
		return
	}
	c.IndentedJSON(http.StatusOK, r)
}

func reviewQuiz(q *quiz) (*result, error) {
	var newScore int
	var lowerScores int
	var r result
	if len(q.Questions) != len(mathQuiz.Questions) {
		return nil, errors.New("Question in Quiz was altered")
	}
	for i, b := range q.Questions {
		if strings.ToUpper(b.ChoosenChoiceID) == mathQuiz.Questions[i].CorrectChoiceID {
			newScore++
		}
	}
	r.Score = newScore
	for _, v := range mathQuiz.Scores {
		if v < newScore {
			lowerScores++
		}
	}

	var percentile = int(float32(lowerScores) / float32(len(mathQuiz.Scores)) * float32(100))
	r.Percentile = fmt.Sprintf("You score is higher than %d%% of the quizzers", percentile)

	mathQuiz.Scores = append(mathQuiz.Scores, newScore)

	return &r, nil
}

func main() {
	router := gin.Default()
	router.GET("/quiz", getQuiz)
	router.POST("/quiz", postQuizScore)
	router.Run("localhost:9090")
}
